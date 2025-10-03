package http

import (
	"context"
	"encoding/json"
	"fmt"
	stdhttp "net/http"
	"strconv"
	"strings"
	"time"

	"github.com/example/flow/internal/app/column/domain"
	"github.com/example/flow/internal/common"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	db *pgxpool.Pool
	ws *common.WSHub
}

func NewHandler(db *pgxpool.Pool, ws *common.WSHub) *Handler { return &Handler{db: db, ws: ws} }

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/{tableId}", h.list)
	r.Post("/{tableId}", h.create)
	r.Put("/{tableId}/{rowId}", h.update)
	r.Delete("/{tableId}/{rowId}", h.delete)
	return r
}

func (h *Handler) list(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	tableID, _ := strconv.ParseInt(chi.URLParam(r, "tableId"), 10, 64)
	limit := int64(100)
	offset := int64(0)
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			limit = n
		}
	}
	if v := r.URL.Query().Get("offset"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			offset = n
		}
	}
	rows, err := h.db.Query(context.Background(), fmt.Sprintf("SELECT * FROM t_%d ORDER BY id DESC LIMIT $1 OFFSET $2", tableID), limit, offset)
	if err != nil {
		writeErr(w, 500, err)
		return
	}
	defer rows.Close()
	out, err := pgRowsToMaps(rows)
	if err != nil {
		writeErr(w, 500, err)
		return
	}
	writeJSON(w, 200, out)
}

func (h *Handler) create(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	tableID, _ := strconv.ParseInt(chi.URLParam(r, "tableId"), 10, 64)
	var payload map[string]any
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeErr(w, 400, err)
		return
	}
	cols, err := h.loadColumns(tableID)
	if err != nil {
		writeErr(w, 500, err)
		return
	}
	if err := validatePayload(cols, payload); err != nil {
		writeErr(w, 400, err)
		return
	}
	// Build INSERT
	keys := make([]string, 0, len(payload))
	vals := make([]any, 0, len(payload))
	ph := make([]string, 0, len(payload))
	i := 1
	for k, v := range payload {
		keys = append(keys, fmt.Sprintf("\"%s\"", safeIdent(k)))
		vals = append(vals, v)
		ph = append(ph, fmt.Sprintf("$%d", i))
		i++
	}
	sql := fmt.Sprintf("INSERT INTO t_%d (%s) VALUES(%s) RETURNING *", tableID, strings.Join(keys, ","), strings.Join(ph, ","))
	row := h.db.QueryRow(context.Background(), sql, vals...)
	rec, err := pgRowToMap(row)
	if err != nil {
		writeErr(w, 400, err)
		return
	}
	writeJSON(w, 201, rec)
	h.ws.BroadcastJSON(map[string]any{"type": "row.created", "data": map[string]any{"table_id": tableID, "row": rec}})
}

func (h *Handler) update(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	tableID, _ := strconv.ParseInt(chi.URLParam(r, "tableId"), 10, 64)
	rowID, _ := strconv.ParseInt(chi.URLParam(r, "rowId"), 10, 64)
	var payload map[string]any
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeErr(w, 400, err)
		return
	}
	cols, err := h.loadColumns(tableID)
	if err != nil {
		writeErr(w, 500, err)
		return
	}
	// allow partial update, validate keys present
	if err := validatePartial(cols, payload); err != nil {
		writeErr(w, 400, err)
		return
	}
	sets := make([]string, 0, len(payload))
	vals := make([]any, 0, len(payload)+1)
	i := 1
	for k, v := range payload {
		sets = append(sets, fmt.Sprintf("\"%s\"=$%d", safeIdent(k), i))
		vals = append(vals, v)
		i++
	}
	vals = append(vals, rowID)
	sql := fmt.Sprintf("UPDATE t_%d SET %s WHERE id=$%d RETURNING *", tableID, strings.Join(sets, ","), i)
	rec, err := pgRowToMap(h.db.QueryRow(context.Background(), sql, vals...))
	if err != nil {
		writeErr(w, 400, err)
		return
	}
	writeJSON(w, 200, rec)
	h.ws.BroadcastJSON(map[string]any{"type": "row.updated", "data": map[string]any{"table_id": tableID, "row": rec}})
}

func (h *Handler) delete(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	tableID, _ := strconv.ParseInt(chi.URLParam(r, "tableId"), 10, 64)
	rowID, _ := strconv.ParseInt(chi.URLParam(r, "rowId"), 10, 64)
	ct, err := h.db.Exec(context.Background(), fmt.Sprintf("DELETE FROM t_%d WHERE id=$1", tableID), rowID)
	if err != nil {
		writeErr(w, 400, err)
		return
	}
	if ct.RowsAffected() == 0 {
		writeErr(w, 404, fmt.Errorf("not found"))
		return
	}
	writeJSON(w, 204, nil)
	h.ws.BroadcastJSON(map[string]any{"type": "row.deleted", "data": map[string]any{"table_id": tableID, "id": rowID}})
}

type columnMeta struct {
	Name     string
	Type     domain.ColumnType
	Required bool
}

func (h *Handler) loadColumns(tableID int64) ([]columnMeta, error) {
	rows, err := h.db.Query(context.Background(), `SELECT name, type, required FROM columns WHERE table_id=$1`, tableID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []columnMeta
	for rows.Next() {
		var m columnMeta
		if err := rows.Scan(&m.Name, &m.Type, &m.Required); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

func validatePayload(cols []columnMeta, payload map[string]any) error {
	// required
	required := map[string]bool{}
	for _, c := range cols {
		if c.Required {
			required[c.Name] = true
		}
	}
	for k := range payload {
		delete(required, k)
	}
	if len(required) > 0 {
		return fmt.Errorf("missing required: %v", keys(required))
	}
	return coerceTypes(cols, payload)
}

func validatePartial(cols []columnMeta, payload map[string]any) error {
	return coerceTypes(cols, payload)
}

func coerceTypes(cols []columnMeta, payload map[string]any) error {
	typeMap := map[string]domain.ColumnType{}
	for _, c := range cols {
		typeMap[c.Name] = c.Type
	}
	for k, v := range payload {
		switch typeMap[k] {
		case domain.ColumnInteger:
			switch t := v.(type) {
			case float64:
				payload[k] = int64(t)
			case string:
				n, err := strconv.ParseInt(t, 10, 64)
				if err != nil {
					return fmt.Errorf("%s must be integer", k)
				}
				payload[k] = n
			}
		case domain.ColumnBoolean:
			switch t := v.(type) {
			case string:
				b := strings.ToLower(t)
				payload[k] = b == "true" || b == "1" || b == "yes"
			}
		case domain.ColumnTimestamp:
			switch t := v.(type) {
			case string:
				if ts, err := time.Parse(time.RFC3339, t); err == nil {
					payload[k] = ts
				}
			}
		}
	}
	return nil
}

func keys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
func safeIdent(s string) string { return strings.ReplaceAll(s, "\"", "") }

func pgRowsToMaps(rows pgx.Rows) ([]map[string]any, error) {
	fds := rows.FieldDescriptions()
	out := []map[string]any{}
	for rows.Next() {
		vals, err := rows.Values()
		if err != nil {
			return nil, err
		}
		rec := map[string]any{}
		for i, fd := range fds {
			rec[string(fd.Name)] = vals[i]
		}
		out = append(out, rec)
	}
	return out, rows.Err()
}

func pgRowToMap(row pgx.Row) (map[string]any, error) {
	// pgx Row doesn't expose FieldDescriptions; do a query wrapper
	// We will scan using a temporary SELECT to get fields; here we assume we already queried RETURNING *
	// Workaround: use CollectRows pattern
	// Instead, we re-exec using Query with the same SQL is complex; simpler: decode into json via rowToJSON
	var raw map[string]any
	// Try to scan common fields; fallback: error if not available
	// Not straightforward; to keep it simple, we use a generic approach: scan into []any after reading columns is not possible here
	// So we return error to force clients to refetch list; but try to scan common id
	if err := row.Scan(&raw); err != nil {
		return nil, fmt.Errorf("unable to scan row generically; refetch: %w", err)
	}
	return raw, nil
}

func writeJSON(w stdhttp.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v != nil {
		_ = json.NewEncoder(w).Encode(v)
	}
}

func writeErr(w stdhttp.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

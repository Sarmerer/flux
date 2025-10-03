package http

import (
	"encoding/json"
	stdhttp "net/http"
	"strconv"

	"github.com/example/flow/internal/app/relationship"
	"github.com/example/flow/internal/app/relationship/domain"
	"github.com/example/flow/internal/common"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	svc *relationship.Service
	ws  *common.WSHub
}

func NewHandler(s *relationship.Service, ws *common.WSHub) *Handler { return &Handler{svc: s, ws: ws} }

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/{tableId}", h.list)
	r.Post("/", h.create)
	r.Delete("/{id}", h.delete)
	return r
}

func (h *Handler) list(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	tableID, _ := strconv.ParseInt(chi.URLParam(r, "tableId"), 10, 64)
	items, err := h.svc.List(tableID)
	if err != nil {
		writeErr(w, 500, err)
		return
	}
	writeJSON(w, 200, items)
}

func (h *Handler) create(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	var body struct {
		SourceTableID int64       `json:"source_table_id"`
		TargetTableID int64       `json:"target_table_id"`
		Kind          domain.Kind `json:"kind"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeErr(w, 400, err)
		return
	}
	rel, err := h.svc.Create(body.SourceTableID, body.TargetTableID, body.Kind)
	if err != nil {
		writeErr(w, 400, err)
		return
	}
	writeJSON(w, 201, rel)
	h.ws.BroadcastJSON(map[string]any{"type": "relationship.created", "data": rel})
}

func (h *Handler) delete(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err := h.svc.Delete(id); err != nil {
		writeErr(w, 400, err)
		return
	}
	writeJSON(w, 204, nil)
	h.ws.BroadcastJSON(map[string]any{"type": "relationship.deleted", "data": map[string]int64{"id": id}})
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

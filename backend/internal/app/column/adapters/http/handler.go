package http

import (
	"encoding/json"
	stdhttp "net/http"
	"strconv"

	"github.com/example/flow/internal/app/column"
	"github.com/example/flow/internal/app/column/domain"
	"github.com/example/flow/internal/common"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	svc *column.Service
	ws  *common.WSHub
}

func NewHandler(svc *column.Service, ws *common.WSHub) *Handler { return &Handler{svc: svc, ws: ws} }

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/{tableId}", h.list)
	r.Post("/{tableId}", h.create)
	r.Put("/{id}", h.update)
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
	tableID, _ := strconv.ParseInt(chi.URLParam(r, "tableId"), 10, 64)
	var body struct {
		Name         string            `json:"name"`
		Type         domain.ColumnType `json:"type"`
		Required     bool              `json:"required"`
		Unique       bool              `json:"unique"`
		DefaultValue *string           `json:"default_value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeErr(w, 400, err)
		return
	}
	c, err := h.svc.Create(tableID, body.Name, body.Type, body.Required, body.Unique, body.DefaultValue)
	if err != nil {
		writeErr(w, 400, err)
		return
	}
	writeJSON(w, 201, c)
	h.ws.BroadcastJSON(map[string]any{"type": "column.created", "data": c})
}

func (h *Handler) update(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	var body struct {
		Name         string  `json:"name"`
		Required     bool    `json:"required"`
		Unique       bool    `json:"unique"`
		DefaultValue *string `json:"default_value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeErr(w, 400, err)
		return
	}
	c, err := h.svc.Update(id, body.Name, body.Required, body.Unique, body.DefaultValue)
	if err != nil {
		writeErr(w, 400, err)
		return
	}
	writeJSON(w, 200, c)
	h.ws.BroadcastJSON(map[string]any{"type": "column.updated", "data": c})
}

func (h *Handler) delete(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err := h.svc.Delete(id); err != nil {
		writeErr(w, 404, err)
		return
	}
	writeJSON(w, 204, nil)
	h.ws.BroadcastJSON(map[string]any{"type": "column.deleted", "data": map[string]int64{"id": id}})
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

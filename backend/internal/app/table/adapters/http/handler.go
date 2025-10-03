package http

import (
	"encoding/json"
	stdhttp "net/http"
	"strconv"

	"github.com/example/flow/internal/app/table"
	"github.com/example/flow/internal/common"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	svc *table.Service
	ws  *common.WSHub
}

func NewHandler(svc *table.Service, ws *common.WSHub) *Handler { return &Handler{svc: svc, ws: ws} }

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.list)
	r.Post("/", h.create)
	r.Put("/{id}", h.update)
	r.Delete("/{id}", h.delete)
	return r
}

func (h *Handler) list(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	items, err := h.svc.List()
	if err != nil {
		writeErr(w, 500, err)
		return
	}
	writeJSON(w, 200, items)
}

func (h *Handler) create(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	var body struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeErr(w, 400, err)
		return
	}
	t, err := h.svc.Create(body.Name)
	if err != nil {
		writeErr(w, 400, err)
		return
	}
	writeJSON(w, 201, t)
	h.ws.BroadcastJSON(map[string]any{"type": "table.created", "data": t})
}

func (h *Handler) update(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	var body struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeErr(w, 400, err)
		return
	}
	t, err := h.svc.Update(id, body.Name)
	if err != nil {
		writeErr(w, 400, err)
		return
	}
	writeJSON(w, 200, t)
	h.ws.BroadcastJSON(map[string]any{"type": "table.updated", "data": t})
}

func (h *Handler) delete(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err := h.svc.Delete(id); err != nil {
		writeErr(w, 404, err)
		return
	}
	writeJSON(w, 204, nil)
	h.ws.BroadcastJSON(map[string]any{"type": "table.deleted", "data": map[string]int64{"id": id}})
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

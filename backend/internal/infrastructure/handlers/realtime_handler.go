package handlers

import (
	"net/http"

	"github.com/flow/internal/infrastructure/realtime"
)

type RealtimeHandler struct {
	hub       *realtime.Hub
	jwtSecret string
}

func NewRealtimeHandler(hub *realtime.Hub, jwtSecret string) *RealtimeHandler {
	return &RealtimeHandler{
		hub:       hub,
		jwtSecret: jwtSecret,
	}
}

func (h *RealtimeHandler) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	realtime.HandleWebSocket(h.hub, h.jwtSecret)(w, r)
}

package handler

import (
	"net/http"
)

type GameplayHandler struct{}

func NewGameplayHandler() *GameplayHandler {
	return &GameplayHandler{}
}

func (h *GameplayHandler) Health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

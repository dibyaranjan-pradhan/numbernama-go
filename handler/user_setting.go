package handler

import (
	"encoding/json"
	"net/http"

	"numbernama-go/service"
)

type UserSettingHandler struct {
	svc *service.UserSettingService
}

func NewUserSettingHandler(svc *service.UserSettingService) *UserSettingHandler {
	return &UserSettingHandler{svc: svc}
}

func (h *UserSettingHandler) Get(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(h.svc.Get())
}

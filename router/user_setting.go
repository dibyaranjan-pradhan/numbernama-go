package router

import (
	"net/http"

	"github.com/gorilla/mux"

	"numbernama-go/handler"
	"numbernama-go/service"
)

func RegisterUserSetting(r *mux.Router, svc *service.UserSettingService) {
	h := handler.NewUserSettingHandler(svc)
	r.Methods(http.MethodGet).Path("/api/settings").HandlerFunc(h.Get)
}

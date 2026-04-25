package router

import (
	"net/http"

	"github.com/gorilla/mux"

	"numbernama-go/handler"
	"numbernama-go/service"
)

func RegisterUser(r *mux.Router, svc *service.UserService) {
	h := handler.NewUserHandler(svc)
	r.Methods(http.MethodGet).Path("/api/user/me").HandlerFunc(h.Me)
}

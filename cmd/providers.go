package main

import (
	"net/http"

	gosocket "github.com/dibyaranjan-pradhan/go-socket"
	"github.com/gorilla/mux"

	"numbernama-go/handler"
	"numbernama-go/middleware"
	"numbernama-go/router"
	"numbernama-go/service"
	"numbernama-go/utils"
)

func provideLogger() *utils.Logger {
	utils.InitLogger("info", false)
	return utils.GetLogger()
}

func provideSocketServer(logger *utils.Logger) *gosocket.Server {
	return gosocket.New(gosocket.Config{
		Logger: utils.GoSocketDiag{L: logger},
	})
}

func assembleRouter(
	logger *utils.Logger,
	socketServer *gosocket.Server,
	gameplaySocketHandler *handler.GameplaySocketHandler,
	gameplayHandler *handler.GameplayHandler,
	userService *service.UserService,
	userSettingService *service.UserSettingService,
) *mux.Router {
	r := mux.NewRouter()
	middleware.Register(r, logger)
	router.RegisterGameplay(r, socketServer, gameplaySocketHandler, gameplayHandler)
	router.RegisterUser(r, userService)
	router.RegisterUserSetting(r, userSettingService)
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		http.NotFound(w, req)
	})
	return r
}

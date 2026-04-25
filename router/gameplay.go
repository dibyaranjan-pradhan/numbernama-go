package router

import (
	"io/fs"
	"net/http"

	gosocket "github.com/dibyaranjan-pradhan/go-socket"
	"github.com/gorilla/mux"

	"numbernama-go/handler"
	"numbernama-go/web"
)

func RegisterGameplay(r *mux.Router, socketServer *gosocket.Server, socketHandler *handler.GameplaySocketHandler, gameplayHandler *handler.GameplayHandler) {
	socketHandler.Register(socketServer)
	r.Path("/ws/gameplay").Handler(socketServer.Handler())

	sub, err := fs.Sub(web.Assets, ".")
	if err != nil {
		panic("gameplay static: " + err.Error())
	}
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.FS(sub))))
	r.Methods(http.MethodGet).Path("/numbers-game").HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		data, readErr := fs.ReadFile(web.Assets, "index.html")
		if readErr != nil {
			http.Error(w, readErr.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(data)
	})
	r.Methods(http.MethodGet).Path("/health").HandlerFunc(gameplayHandler.Health)
}

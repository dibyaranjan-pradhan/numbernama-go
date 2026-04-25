package middleware

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"numbernama-go/utils"
)

func Register(r *mux.Router, logger *utils.Logger) {
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			start := time.Now()
			w.Header().Set("Access-Control-Allow-Origin", "*")
			next.ServeHTTP(w, req)
			logger.Infof("http_request method=%s path=%s duration=%s", req.Method, req.URL.Path, time.Since(start))
		})
	})
}

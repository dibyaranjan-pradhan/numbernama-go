package main

import (
	"net/http"
	"os"
	"strconv"
)

func main() {
	router, err := InitializeMux()
	if err != nil {
		panic(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "7002"
	}
	if _, err := strconv.Atoi(port); err != nil {
		panic("invalid PORT: " + port)
	}

	srv := &http.Server{Addr: ":" + port, Handler: router}
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}

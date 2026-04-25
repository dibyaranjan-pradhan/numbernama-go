//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/gorilla/mux"

	"numbernama-go/handler"
	"numbernama-go/repo"
	"numbernama-go/service"
)

// InitializeMux builds the HTTP router via Google Wire.
func InitializeMux() (*mux.Router, error) {
	wire.Build(
		repo.NewMemoryGameplay,
		service.NewGameplayService,
		service.NewUserService,
		service.NewUserSettingService,
		provideLogger,
		handler.NewGameplayHandler,
		handler.NewGameplaySocketHandler,
		provideSocketServer,
		assembleRouter,
	)
	return nil, nil
}

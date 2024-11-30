package main

import (
	"os"

	"github.com/PRYVT/posting/pkg/command/httphandler"
	"github.com/PRYVT/posting/pkg/command/httphandler/controller"
	"github.com/PRYVT/utils/pkg/auth"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	tokenManager, err := auth.NewTokenManager()
	if err != nil {
		log.Error().Err(err).Msg("Unsuccessful initialization of token manager")
		return
	}
	aut := auth.NewAuthMiddleware(tokenManager)
	uc := controller.NewPostController(tokenManager)
	h := httphandler.NewHttpHandler(uc, aut)

	h.Start()
}

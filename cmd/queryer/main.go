package main

import (
	"os"

	"github.com/L4B0MB4/EVTSRC/pkg/client"
	"github.com/PRYVT/posting/pkg/query/eventhandling"
	"github.com/PRYVT/posting/pkg/query/httphandler"
	"github.com/PRYVT/posting/pkg/query/httphandler/controller"
	"github.com/PRYVT/posting/pkg/query/store"
	"github.com/PRYVT/posting/pkg/query/store/repository"
	"github.com/PRYVT/utils/pkg/auth"
	"github.com/PRYVT/utils/pkg/eventpolling"
	utilsRepo "github.com/PRYVT/utils/pkg/store/repository"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	db := store.DatabaseConnection{}
	db.SetUp()
	conn, err := db.GetDbConnection()
	if err != nil {
		log.Error().Err(err).Msg("Unsuccessfull initalization of db")
		return
	}
	log.Debug().Msg("Db Connection was successful")

	c, err := client.NewEventSourcingHttpClient(client.RetrieveEventSourcingClientUrl())
	if err != nil {
		log.Error().Err(err).Msg("Unsuccessful initialization of client")
		return
	}
	tokenManager, err := auth.NewTokenManager()
	if err != nil {
		log.Error().Err(err).Msg("Unsuccessful initialization of token manager")
		return
	}
	eventRepo := utilsRepo.NewEventRepository(conn)
	userRepo := repository.NewUserRepository(conn)
	uc := controller.NewPostController(userRepo, tokenManager)
	aut := auth.NewAuthMiddleware(tokenManager)
	h := httphandler.NewHttpHandler(uc, aut)

	userEventHandler := eventhandling.NewPostEventHandler(userRepo)

	eventPolling := eventpolling.NewEventPolling(c, eventRepo, userEventHandler)
	go eventPolling.PollEvents()

	h.Start()
}

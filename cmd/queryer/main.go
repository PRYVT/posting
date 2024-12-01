package main

import (
	"os"

	"github.com/L4B0MB4/EVTSRC/pkg/client"
	tcpClient "github.com/L4B0MB4/EVTSRC/pkg/tcp/client"
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
	eventRepo := utilsRepo.NewEventRepository(conn)
	userRepo := repository.NewUserRepository(conn)
	userEventHandler := eventhandling.NewPostEventHandler(userRepo)
	uc := controller.NewPostController(userRepo, userEventHandler)
	aut := auth.NewAuthMiddleware()
	wsH := controller.NewWsController(userEventHandler)
	h := httphandler.NewHttpHandler(uc, aut, wsH)

	eventPolling := eventpolling.NewEventPolling(c, eventRepo, userEventHandler)

	tcpC, err := tcpClient.NewTcpEventClient()
	if err != nil {
		log.Error().Err(err).Msg("Unsuccessful initialization of tcp client")
		return
	}
	channel := make(chan string, 1)
	go tcpC.ListenForEvents(channel)

	eventPolling.PollEventsUntilEmpty()
	go func() {
		for {
			select {
			case event := <-channel:
				log.Info().Msgf("Received event: %s", event)
				eventPolling.PollEventsUntilEmpty()
			}
			log.Debug().Msg("New event received finished")
		}
	}()
	h.Start()
}

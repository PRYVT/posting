package eventhandling

import (
	"sync"

	"github.com/L4B0MB4/EVTSRC/pkg/models"
	"github.com/PRYVT/posting/pkg/aggregates"
	"github.com/PRYVT/posting/pkg/query/store/repository"
	ws "github.com/PRYVT/posting/pkg/query/websocket"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type PostEventHandler struct {
	postRepo      *repository.PostRepository
	wsConnections []*ws.WebsocketConnection
	mu            sync.Mutex
}

func NewPostEventHandler(postRepo *repository.PostRepository) *PostEventHandler {
	return &PostEventHandler{
		postRepo:      postRepo,
		wsConnections: []*ws.WebsocketConnection{},
	}
}

func (eh *PostEventHandler) AddWebsocketConnection(conn *ws.WebsocketConnection) {
	eh.mu.Lock()
	defer eh.mu.Unlock()
	eh.wsConnections = append(eh.wsConnections, conn)
}

func removeDisconnectedSockets(slice []*ws.WebsocketConnection) []*ws.WebsocketConnection {
	output := []*ws.WebsocketConnection{}
	for _, element := range slice {
		if element.IsConnected {
			output = append(output, element)
		}
	}
	return output
}

func (eh *PostEventHandler) HandleEvent(event models.Event) error {
	log.Debug().Msg("Handling event")
	if event.AggregateType == "post" {
		log.Debug().Msg("Handling post event")
		ua, err := aggregates.NewPostAggregate(uuid.MustParse(event.AggregateId))
		if err != nil {
			return err
		}
		p := aggregates.GetPostModelFromAggregate(ua)
		err = eh.postRepo.AddOrReplacePost(p)
		if err != nil {
			log.Err(err).Msg("Error while processing user event")
			return err
		}
		for _, conn := range eh.wsConnections {
			if !conn.IsAuthenticated {
				continue
			}
			err := conn.WriteJSON(p)
			if err != nil {
				log.Warn().Err(err).Msg("Error while writing to websocket connection")
			}
		}
		eh.mu.Lock()
		defer eh.mu.Unlock()
		eh.wsConnections = removeDisconnectedSockets(eh.wsConnections)
		log.Trace().Msgf("Number of active connections: %d", len(eh.wsConnections))
	}
	return nil
}

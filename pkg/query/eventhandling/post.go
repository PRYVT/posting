package eventhandling

import (
	"github.com/L4B0MB4/EVTSRC/pkg/models"
	"github.com/PRYVT/posting/pkg/aggregates"
	"github.com/PRYVT/posting/pkg/query/store/repository"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type UserEventHandler struct {
	postRepo *repository.PostRepository
}

func NewPostEventHandler(postRepo *repository.PostRepository) *UserEventHandler {
	return &UserEventHandler{
		postRepo: postRepo,
	}
}

func (eh *UserEventHandler) HandleEvent(event models.Event) error {
	if event.AggregateType == "post" {
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
	}
	return nil
}

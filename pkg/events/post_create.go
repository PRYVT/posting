package events

import (
	"time"

	"github.com/L4B0MB4/EVTSRC/pkg/models"
	m "github.com/PRYVT/posting/pkg/models/command"
	"github.com/google/uuid"
)

type PostCreatedEvent struct {
	Text         string
	ImageBase64  string
	CreationDate time.Time
	UserId       uuid.UUID
}

func NewPostCreateEvent(cp m.CreatePost) *models.ChangeTrackedEvent {

	b := UnsafeSerializeAny(PostCreatedEvent{
		Text:         cp.Text,
		ImageBase64:  cp.ImageBase64,
		CreationDate: time.Now(),
		UserId:       cp.UserId,
	})
	return &models.ChangeTrackedEvent{
		Event: models.Event{
			Name: "PostCreatedEvent",
			Data: b,
		},
	}
}

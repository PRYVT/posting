package events

import (
	"time"

	"github.com/L4B0MB4/EVTSRC/pkg/models"
	m "github.com/PRYVT/posting/pkg/models/command"
)

type PostCreatedEvent struct {
	Text         string
	ImageBase64  string
	CreationDate time.Time
}

func NewPostCreateEvent(cp m.CreatePost) *models.ChangeTrackedEvent {

	b := UnsafeSerializeAny(PostCreatedEvent{
		Text:         cp.Text,
		ImageBase64:  cp.ImageBase64,
		CreationDate: time.Now(),
	})
	return &models.ChangeTrackedEvent{
		Event: models.Event{
			Name: "PostCreatedEvent",
			Data: b,
		},
	}
}

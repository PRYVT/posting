package aggregates

import (
	"fmt"
	"time"

	"github.com/L4B0MB4/EVTSRC/pkg/client"
	"github.com/L4B0MB4/EVTSRC/pkg/models"
	"github.com/PRYVT/posting/pkg/events"
	"github.com/PRYVT/posting/pkg/models/command"
	"github.com/google/uuid"
)

type PostAggregate struct {
	Text          string
	ImageBase64   string
	UserId        uuid.UUID
	ChangeDate    time.Time
	Events        []models.ChangeTrackedEvent
	aggregateType string
	AggregateId   uuid.UUID
	client        *client.EventSourcingHttpClient
}

func NewPostAggregate(id uuid.UUID) (*PostAggregate, error) {

	c, err := client.NewEventSourcingHttpClient(client.RetrieveEventSourcingClientUrl())
	if err != nil {
		panic(err)
	}
	iter, err := c.GetEventsOrdered(id.String())
	if err != nil {
		return nil, fmt.Errorf("COULDN'T RETRIEVE EVENTS ")
	}
	ua := &PostAggregate{
		client:        c,
		Events:        []models.ChangeTrackedEvent{},
		aggregateType: "post",
		AggregateId:   id,
		ChangeDate:    time.Date(2000, 0, 0, 0, 0, 0, 0, time.UTC),
	}

	for {
		ev, ok := iter.Next()
		if !ok {
			break
		}
		changeTrackedEv := models.ChangeTrackedEvent{
			Event: *ev,
			IsNew: false,
		}
		ua.addEvent(&changeTrackedEv)
	}
	return ua, nil
}

func (pa *PostAggregate) apply_PostCreatedEvent(e *events.PostCreatedEvent) {
	pa.Text = e.Text
	pa.ImageBase64 = e.ImageBase64
	pa.ChangeDate = e.CreationDate
	pa.UserId = e.UserId
}

func (ua *PostAggregate) addEvent(ev *models.ChangeTrackedEvent) {
	switch ev.Name {
	case "PostCreatedEvent":
		e := events.UnsafeDeserializeAny[events.PostCreatedEvent](ev.Data)
		ua.apply_PostCreatedEvent(e)
	default:
		panic(fmt.Errorf("NO KNOWN EVENT %v", ev))
	}
	if ev.Version == 0 {
		ev.IsNew = true
	}
	v := len(ua.Events) + 1 //for validation we need to start at 1
	ev.Version = int64(v)
	ev.AggregateType = ua.aggregateType
	ev.AggregateId = ua.AggregateId.String()
	ua.Events = append(ua.Events, *ev)
}

func (ua *PostAggregate) saveChanges() error {
	return ua.client.AddEvents(ua.AggregateId.String(), ua.Events)
}

func (ua *PostAggregate) CreatePost(postCreate command.CreatePost) error {

	if len(ua.Events) != 0 {
		return fmt.Errorf("post already exists")
	}
	if postCreate.Text == "" && postCreate.ImageBase64 == "" {
		return fmt.Errorf("post must have text or image")
	}

	ua.addEvent(events.NewPostCreateEvent(postCreate))
	err := ua.saveChanges()
	if err != nil {
		return fmt.Errorf("ERROR ")
	}
	return nil
}

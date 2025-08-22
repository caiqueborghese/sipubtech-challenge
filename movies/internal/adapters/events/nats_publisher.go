package events

import (
	"context"
	"encoding/json"
	"time"

	"github.com/caiqueborghese/sipubtech-challenge/movies/internal/domain"
	"github.com/caiqueborghese/sipubtech-challenge/movies/internal/ports"
	"github.com/nats-io/nats.go"
)

type NatsPublisher struct {
	nc             *nats.Conn
	subjectCreated string
	subjectDeleted string
}

func NewNatsPublisher(nc *nats.Conn, subjCreated, subjDeleted string) ports.EventPublisher {
	return &NatsPublisher{nc: nc, subjectCreated: subjCreated, subjectDeleted: subjDeleted}
}

type eventEnvelope struct {
	Type       string      `json:"type"`
	OccurredAt time.Time   `json:"occurred_at"`
	Payload    interface{} `json:"payload"`
}

func (p *NatsPublisher) MovieCreated(ctx context.Context, m domain.Movie) error {
	ev := eventEnvelope{
		Type:       "movies.created",
		OccurredAt: time.Now().UTC(),
		Payload: struct {
			ID    string `json:"id"`
			Title string `json:"title"`
			Year  int    `json:"year"`
		}{ID: m.ID, Title: m.Title, Year: m.Year},
	}
	b, _ := json.Marshal(ev)
	return p.nc.Publish(p.subjectCreated, b)
}

func (p *NatsPublisher) MovieDeleted(ctx context.Context, id string) error {
	ev := eventEnvelope{
		Type:       "movies.deleted",
		OccurredAt: time.Now().UTC(),
		Payload: struct {
			ID string `json:"id"`
		}{ID: id},
	}
	b, _ := json.Marshal(ev)
	return p.nc.Publish(p.subjectDeleted, b)
}

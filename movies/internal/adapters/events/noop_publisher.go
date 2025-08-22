package events

import (
	"context"

	"github.com/caiqueborghese/sipubtech-challenge/movies/internal/domain"
	"github.com/caiqueborghese/sipubtech-challenge/movies/internal/ports"
)

type NoopPublisher struct{}

func NewNoopPublisher() ports.EventPublisher { return &NoopPublisher{} }

func (NoopPublisher) MovieCreated(ctx context.Context, m domain.Movie) error { return nil }
func (NoopPublisher) MovieDeleted(ctx context.Context, id string) error      { return nil }

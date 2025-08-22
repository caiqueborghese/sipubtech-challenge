package ports

import (
	"context"

	"github.com/caiqueborghese/sipubtech-challenge/movies/internal/domain"
)

type EventPublisher interface {
	MovieCreated(ctx context.Context, m domain.Movie) error
	MovieDeleted(ctx context.Context, id string) error
}

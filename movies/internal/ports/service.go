package ports

import (
	"context"

	"github.com/caiqueborghese/sipubtech-challenge/movies/internal/domain"
)

type MovieService interface {
	List(ctx context.Context) ([]domain.Movie, error)
	Get(ctx context.Context, id string) (*domain.Movie, error)
	Create(ctx context.Context, m domain.Movie) (*domain.Movie, error)
	Delete(ctx context.Context, id string) error

	// Usado no bootstrap do servidor para popular base, se necessário
	EnsureSeed(ctx context.Context, seed []domain.Movie) (inserted int, err error)
}

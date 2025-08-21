package ports

import (
	"context"

	"github.com/caiqueborghese/sipubtech-challenge/api-gateway/internal/domain"
)

// Porta usada pelos handlers HTTP.
type MovieService interface {
	List(ctx context.Context) ([]domain.Movie, error)
	Get(ctx context.Context, id string) (*domain.Movie, error)
	Create(ctx context.Context, in domain.Movie) (*domain.Movie, error)
	Delete(ctx context.Context, id string) error
}

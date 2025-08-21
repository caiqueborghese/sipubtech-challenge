package ports

import (
	"context"

	"github.com/caiqueborghese/sipubtech-challenge/api-gateway/internal/domain"
)

// MoviesClient porta de saída ( gRPC implementa)
type MoviesClient interface {
	List(ctx context.Context) ([]domain.Movie, error)
	Get(ctx context.Context, id string) (*domain.Movie, error)
	Create(ctx context.Context, m domain.Movie) (*domain.Movie, error)
	Delete(ctx context.Context, id string) error
}

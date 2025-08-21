package ports

import (
	"context"

	"github.com/caiqueborghese/sipubtech-challenge/movies/internal/domain"
)

type MovieRepository interface {
	List(ctx context.Context) ([]domain.Movie, error)
	Get(ctx context.Context, id string) (*domain.Movie, error)
	Create(ctx context.Context, m *domain.Movie) (*domain.Movie, error)
	Delete(ctx context.Context, id string) error

	// Suporte a seed idempotente
	Count(ctx context.Context) (int64, error)
	BulkInsertIgnoreDuplicates(ctx context.Context, ms []domain.Movie) (int, error)
}

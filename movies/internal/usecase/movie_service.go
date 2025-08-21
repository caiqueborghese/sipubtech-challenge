package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/caiqueborghese/sipubtech-challenge/movies/internal/domain"
	"github.com/caiqueborghese/sipubtech-challenge/movies/internal/ports"
)

var _ ports.MovieService = (*movieService)(nil)

type movieService struct {
	repo ports.MovieRepository
}

func NewMovieService(repo ports.MovieRepository) ports.MovieService {
	return &movieService{repo: repo}
}

func (s *movieService) List(ctx context.Context) ([]domain.Movie, error) {
	return s.repo.List(ctx)
}

func (s *movieService) Get(ctx context.Context, id string) (*domain.Movie, error) {
	if id == "" {
		return nil, domain.ErrInvalidID
	}
	return s.repo.Get(ctx, id)
}

func (s *movieService) Create(ctx context.Context, m domain.Movie) (*domain.Movie, error) {
	m.Normalize()
	if err := m.Validate(); err != nil {
		return nil, err
	}
	return s.repo.Create(ctx, &m)
}

func (s *movieService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return domain.ErrInvalidID
	}
	return s.repo.Delete(ctx, id)
}

func (s *movieService) EnsureSeed(ctx context.Context, seed []domain.Movie) (int, error) {
	if len(seed) == 0 {
		return 0, nil
	}

	clean := make([]domain.Movie, 0, len(seed))
	for i := range seed {
		seed[i].Normalize()

		if seed[i].Title == "" || seed[i].Year == 0 {
			continue
		}
		if err := seed[i].Validate(); err != nil {

			continue
		}
		clean = append(clean, seed[i])
	}

	if len(clean) == 0 {
		return 0, errors.New("no valid items to seed")
	}

	inserted, err := s.repo.BulkInsertIgnoreDuplicates(ctx, clean)
	if err != nil {
		return 0, fmt.Errorf("seed insert failed: %w", err)
	}
	return inserted, nil
}

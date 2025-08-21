package usecase

import (
	"context"
	"errors"

	"github.com/caiqueborghese/sipubtech-challenge/api-gateway/internal/domain"
	moviespb "github.com/caiqueborghese/sipubtech-challenge/proto/moviespb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type MovieService interface {
	List() ([]domain.Movie, error)
	Get(id string) (*domain.Movie, error)
	Create(m *domain.Movie) (*domain.Movie, error)
	Delete(id string) error
}

type movieService struct {
	client moviespb.MovieServiceClient
}

func NewMovieService(client moviespb.MovieServiceClient) MovieService {
	return &movieService{client: client}
}

func (s *movieService) List() ([]domain.Movie, error) {
	res, err := s.client.ListMovies(context.Background(), &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	out := make([]domain.Movie, 0, len(res.Movies))
	for _, m := range res.Movies {
		out = append(out, domain.Movie{
			ID:    m.GetId(),
			Title: m.GetTitle(),
			Year:  int(m.GetYear()),
		})
	}
	return out, nil
}

func (s *movieService) Get(id string) (*domain.Movie, error) {
	if id == "" {
		return nil, errors.New("id required")
	}
	res, err := s.client.GetMovie(context.Background(), &moviespb.GetMovieRequest{Id: id})
	if err != nil {
		return nil, err
	}
	m := res.GetMovie()
	if m == nil {
		return nil, domain.ErrNotFound
	}
	return &domain.Movie{
		ID:    m.GetId(),
		Title: m.GetTitle(),
		Year:  int(m.GetYear()),
	}, nil
}

func (s *movieService) Create(in *domain.Movie) (*domain.Movie, error) {
	if in == nil {
		return nil, errors.New("movie required")
	}
	in.Normalize()
	if err := in.Validate(); err != nil {
		return nil, err
	}
	req := &moviespb.CreateMovieRequest{
		Title: in.Title,
		Year:  int32(in.Year),
	}
	res, err := s.client.CreateMovie(context.Background(), req)
	if err != nil {
		return nil, err
	}
	m := res.GetMovie()
	return &domain.Movie{
		ID:    m.GetId(),
		Title: m.GetTitle(),
		Year:  int(m.GetYear()),
	}, nil
}

func (s *movieService) Delete(id string) error {
	if id == "" {
		return errors.New("id required")
	}
	_, err := s.client.DeleteMovie(context.Background(), &moviespb.DeleteMovieRequest{Id: id})
	return err
}

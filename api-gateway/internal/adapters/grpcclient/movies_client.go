package grpcclient

import (
	"context"
	"fmt"

	"github.com/caiqueborghese/sipubtech-challenge/api-gateway/internal/domain"
	"github.com/caiqueborghese/sipubtech-challenge/api-gateway/internal/ports"
	moviespb "github.com/caiqueborghese/sipubtech-challenge/proto/moviespb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ ports.MoviesClient = (*Client)(nil)

type Client struct {
	cli moviespb.MovieServiceClient
}

// Conecta via endereço (ex.: "movies:50051"). Retorna o client, um close() e erro.
func NewFromAddr(ctx context.Context, addr string) (ports.MoviesClient, func() error, error) {
	conn, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("grpc dial %s: %w", addr, err)
	}
	return &Client{cli: moviespb.NewMovieServiceClient(conn)}, conn.Close, nil
}

// Para injetar um *grpc.ClientConn já existente (útil em testes/injeção).
func New(conn *grpc.ClientConn) ports.MoviesClient {
	return &Client{cli: moviespb.NewMovieServiceClient(conn)}
}

func (c *Client) List(ctx context.Context) ([]domain.Movie, error) {
	res, err := c.cli.ListMovies(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	out := make([]domain.Movie, 0, len(res.Movies))
	for _, m := range res.Movies {
		out = append(out, fromPB(m))
	}
	return out, nil
}

func (c *Client) Get(ctx context.Context, id string) (*domain.Movie, error) {
	res, err := c.cli.GetMovie(ctx, &moviespb.GetMovieRequest{Id: id})
	if err != nil {
		return nil, err
	}
	dm := fromPB(res.Movie)
	return &dm, nil
}

func (c *Client) Create(ctx context.Context, m domain.Movie) (*domain.Movie, error) {
	res, err := c.cli.CreateMovie(ctx, &moviespb.CreateMovieRequest{
		Title: m.Title,
		Year:  int32(m.Year),
	})
	if err != nil {
		return nil, err
	}
	dm := fromPB(res.Movie)
	return &dm, nil
}

func (c *Client) Delete(ctx context.Context, id string) error {
	_, err := c.cli.DeleteMovie(ctx, &moviespb.DeleteMovieRequest{Id: id})
	return err
}

func fromPB(m *moviespb.Movie) domain.Movie {
	if m == nil {
		return domain.Movie{}
	}
	return domain.Movie{
		ID:    m.Id,
		Title: m.Title,
		Year:  int(m.Year),
	}
}

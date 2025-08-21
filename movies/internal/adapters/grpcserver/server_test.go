package grpcserver

import (
	"context"
	"net"
	"testing"

	"github.com/caiqueborghese/sipubtech-challenge/movies/internal/domain"
	"github.com/caiqueborghese/sipubtech-challenge/movies/internal/ports"
	moviespb "github.com/caiqueborghese/sipubtech-challenge/proto/moviespb"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/emptypb"
)

type fakeSvc struct{}

func (f fakeSvc) List(ctx context.Context) ([]domain.Movie, error) {
	return []domain.Movie{{ID: "8", Title: "X", Year: 2000}}, nil
}
func (f fakeSvc) Get(ctx context.Context, id string) (*domain.Movie, error) {
	return &domain.Movie{ID: id, Title: "One", Year: 1999}, nil
}
func (f fakeSvc) Create(ctx context.Context, m domain.Movie) (*domain.Movie, error) {
	m.ID = "new"
	return &m, nil
}
func (f fakeSvc) Delete(ctx context.Context, id string) error { return nil }
func (f fakeSvc) EnsureSeed(ctx context.Context, seed []domain.Movie) (int, error) {
	return 0, nil
}

var _ ports.MovieService = (*fakeSvc)(nil)

func dialBuf(s *grpc.Server) (*grpc.ClientConn, func(), error) {
	const bufSize = 1024 * 1024
	lis := bufconn.Listen(bufSize)
	go func() { _ = s.Serve(lis) }()
	conn, err := grpc.DialContext(
		context.Background(),
		"bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) { return lis.Dial() }),
		grpc.WithInsecure(),
	)
	cleanup := func() { conn.Close(); s.Stop(); lis.Close() }
	return conn, cleanup, err
}

func TestListMovies_Bufconn(t *testing.T) {
	s := grpc.NewServer()
	moviespb.RegisterMovieServiceServer(s, New(fakeSvc{}))

	conn, cleanup, err := dialBuf(s)
	require.NoError(t, err)
	defer cleanup()

	cli := moviespb.NewMovieServiceClient(conn)
	resp, err := cli.ListMovies(context.Background(), &emptypb.Empty{})
	require.NoError(t, err)
	require.Len(t, resp.GetMovies(), 1)
	require.Equal(t, "8", resp.GetMovies()[0].GetId())
	require.Equal(t, int32(2000), resp.GetMovies()[0].GetYear())
}

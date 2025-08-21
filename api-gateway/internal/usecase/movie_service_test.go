package usecase

import (
	"context"
	"testing"

	gdomain "github.com/caiqueborghese/sipubtech-challenge/api-gateway/internal/domain"
	moviespb "github.com/caiqueborghese/sipubtech-challenge/proto/moviespb"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type fakeClient struct {
	list   []*moviespb.Movie
	get    *moviespb.GetMovieResponse
	create *moviespb.CreateMovieResponse
	delErr error
}

func (f *fakeClient) ListMovies(ctx context.Context, _ *emptypb.Empty, _ ...grpc.CallOption) (*moviespb.ListMoviesResponse, error) {
	return &moviespb.ListMoviesResponse{Movies: f.list}, nil
}
func (f *fakeClient) GetMovie(ctx context.Context, in *moviespb.GetMovieRequest, _ ...grpc.CallOption) (*moviespb.GetMovieResponse, error) {
	return f.get, nil
}
func (f *fakeClient) CreateMovie(ctx context.Context, in *moviespb.CreateMovieRequest, _ ...grpc.CallOption) (*moviespb.CreateMovieResponse, error) {
	return f.create, nil
}
func (f *fakeClient) DeleteMovie(ctx context.Context, in *moviespb.DeleteMovieRequest, _ ...grpc.CallOption) (*moviespb.DeleteMovieResponse, error) {
	if f.delErr != nil {
		return nil, f.delErr
	}
	return &moviespb.DeleteMovieResponse{Success: true}, nil
}

func TestGatewayUsecase_List_MapsFields(t *testing.T) {
	cli := &fakeClient{
		list: []*moviespb.Movie{
			{Id: "8", Title: "X", Year: 1999},
		},
	}
	svc := NewMovieService(cli)

	got, err := svc.List()
	require.NoError(t, err)
	require.Equal(t, []gdomain.Movie{{ID: "8", Title: "X", Year: 1999}}, got)
}

func TestGatewayUsecase_Create_Validate(t *testing.T) {
	cli := &fakeClient{
		create: &moviespb.CreateMovieResponse{
			Movie: &moviespb.Movie{Id: "new", Title: "Ok", Year: 2020},
		},
	}
	svc := NewMovieService(cli)

	// inválido
	_, err := svc.Create(&gdomain.Movie{Title: "", Year: 2020})
	require.Error(t, err)

	// válido
	out, err := svc.Create(&gdomain.Movie{Title: "Ok", Year: 2020})
	require.NoError(t, err)
	require.Equal(t, "new", out.ID)
}

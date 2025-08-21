package grpcserver

import (
	"context"
	"net"

	"github.com/caiqueborghese/sipubtech-challenge/movies/internal/domain"
	"github.com/caiqueborghese/sipubtech-challenge/movies/internal/ports"
	moviespb "github.com/caiqueborghese/sipubtech-challenge/proto/moviespb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	moviespb.UnimplementedMovieServiceServer
	svc ports.MovieService
}

func New(svc ports.MovieService) *Server { return &Server{svc: svc} }

func toPB(m domain.Movie) *moviespb.Movie {
	return &moviespb.Movie{
		Id:    m.ID,
		Title: m.Title,
		Year:  int32(m.Year),
	}
}

func toStatusErr(err error) error {
	switch err {
	case nil:
		return nil
	case domain.ErrNotFound:
		return status.Error(codes.NotFound, err.Error())
	case domain.ErrInvalidID:
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		// validações de domínio diversas
		return status.Error(codes.Unknown, err.Error())
	}
}

func (s *Server) ListMovies(ctx context.Context, _ *emptypb.Empty) (*moviespb.ListMoviesResponse, error) {
	ms, err := s.svc.List(ctx)
	if err != nil {
		return nil, toStatusErr(err)
	}
	out := make([]*moviespb.Movie, 0, len(ms))
	for _, m := range ms {
		out = append(out, toPB(m))
	}
	return &moviespb.ListMoviesResponse{Movies: out}, nil
}

func (s *Server) GetMovie(ctx context.Context, in *moviespb.GetMovieRequest) (*moviespb.GetMovieResponse, error) {
	m, err := s.svc.Get(ctx, in.GetId())
	if err != nil {
		return nil, toStatusErr(err)
	}
	return &moviespb.GetMovieResponse{Movie: toPB(*m)}, nil
}

func (s *Server) CreateMovie(ctx context.Context, in *moviespb.CreateMovieRequest) (*moviespb.CreateMovieResponse, error) {
	n := domain.Movie{
		Title: in.GetTitle(),
		Year:  int(in.GetYear()),
	}
	created, err := s.svc.Create(ctx, n)
	if err != nil {
		return nil, toStatusErr(err)
	}
	return &moviespb.CreateMovieResponse{Movie: toPB(*created)}, nil
}

func (s *Server) DeleteMovie(ctx context.Context, in *moviespb.DeleteMovieRequest) (*moviespb.DeleteMovieResponse, error) {
	if err := s.svc.Delete(ctx, in.GetId()); err != nil {
		return nil, toStatusErr(err)
	}
	return &moviespb.DeleteMovieResponse{Success: true}, nil
}

func RunGRPCServer(svc ports.MovieService, grpcAddr string) error {
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	moviespb.RegisterMovieServiceServer(s, New(svc))
	reflection.Register(s)
	return s.Serve(lis)
}

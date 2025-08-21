package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/caiqueborghese/sipubtech-challenge/movies/internal/domain"
	"github.com/caiqueborghese/sipubtech-challenge/movies/internal/ports/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestList_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockMovieRepository(ctrl)
	svc := NewMovieService(mockRepo)

	want := []domain.Movie{{ID: "8", Title: "Movie 8", Year: 1999}}
	mockRepo.EXPECT().List(gomock.Any()).Return(want, nil)

	got, err := svc.List(context.Background())
	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestGet_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockMovieRepository(ctrl)
	svc := NewMovieService(mockRepo)

	_, err := svc.Get(context.Background(), "")
	require.ErrorIs(t, err, domain.ErrInvalidID)
}

func TestCreate_Validation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockMovieRepository(ctrl)
	svc := NewMovieService(mockRepo)

	_, err := svc.Create(context.Background(), domain.Movie{Title: "", Year: 1999})
	require.Error(t, err)

	in := domain.Movie{Title: "Ok", Year: 2000}
	out := in
	out.ID = "new-id"
	mockRepo.
		EXPECT().
		Create(gomock.Any(), &in).
		Return(&out, nil)

	got, err := svc.Create(context.Background(), in)
	require.NoError(t, err)
	require.Equal(t, &out, got)
}

func TestDelete_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockMovieRepository(ctrl)
	svc := NewMovieService(mockRepo)

	mockRepo.EXPECT().Delete(gomock.Any(), "8").Return(nil)
	require.NoError(t, svc.Delete(context.Background(), "8"))
}

func TestDelete_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockMovieRepository(ctrl)
	svc := NewMovieService(mockRepo)

	require.ErrorIs(t, svc.Delete(context.Background(), ""), domain.ErrInvalidID)
}

func TestEnsureSeed_FiltersAndInserts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockMovieRepository(ctrl)
	svc := NewMovieService(mockRepo)

	seed := []domain.Movie{
		{Title: "Valid", Year: 2001},
		{Title: "", Year: 2002},
		{Title: "OutOfRange", Year: 0},
	}

	mockRepo.
		EXPECT().
		BulkInsertIgnoreDuplicates(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, ms []domain.Movie) (int, error) {
			require.Len(t, ms, 1)
			require.Equal(t, "Valid", ms[0].Title)
			return 1, nil
		})

	ins, err := svc.EnsureSeed(context.Background(), seed)
	require.NoError(t, err)
	require.Equal(t, 1, ins)
}

func TestEnsureSeed_NoValidItems(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockMovieRepository(ctrl)
	svc := NewMovieService(mockRepo)

	seed := []domain.Movie{
		{Title: "", Year: 2000},
		{Title: "X", Year: 0},
	}

	ins, err := svc.EnsureSeed(context.Background(), seed)
	require.Error(t, err)
	require.True(t, errors.Is(err, errors.New("no valid items to seed")) || err != nil)
	require.Equal(t, 0, ins)
}

package usecase

import (
	"context"
	"errors"
	"strconv"
	"testing"

	"github.com/caiqueborghese/sipubtech-challenge/movies/internal/domain"
	"github.com/caiqueborghese/sipubtech-challenge/movies/internal/ports"
	"github.com/stretchr/testify/require"
)

type memRepo struct {
	byID  map[string]domain.Movie
	byKey map[string]string
	next  int
}

func newMemRepo() *memRepo {
	return &memRepo{
		byID:  make(map[string]domain.Movie),
		byKey: make(map[string]string),
		next:  1,
	}
}

func keyOf(m domain.Movie) string { return m.Title + "|" + strconv.Itoa(m.Year) }

func (r *memRepo) List(ctx context.Context) ([]domain.Movie, error) {
	out := make([]domain.Movie, 0, len(r.byID))
	for _, m := range r.byID {
		out = append(out, m)
	}
	return out, nil
}
func (r *memRepo) Get(ctx context.Context, id string) (*domain.Movie, error) {
	m, ok := r.byID[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return &m, nil
}
func (r *memRepo) Create(ctx context.Context, m *domain.Movie) (*domain.Movie, error) {
	k := keyOf(*m)
	if _, dup := r.byKey[k]; dup {
		return nil, errors.New("duplicate (title,year)")
	}
	id := "gen-" + strconv.Itoa(r.next)
	r.next++
	cp := *m
	cp.ID = id
	r.byID[id] = cp
	r.byKey[k] = id
	return &cp, nil
}
func (r *memRepo) Delete(ctx context.Context, id string) error {
	m, ok := r.byID[id]
	if !ok {
		return errors.New("not found")
	}
	delete(r.byID, id)
	delete(r.byKey, keyOf(m))
	return nil
}
func (r *memRepo) Count(ctx context.Context) (int64, error) { return int64(len(r.byID)), nil }
func (r *memRepo) BulkInsertIgnoreDuplicates(ctx context.Context, ms []domain.Movie) (int, error) {
	ins := 0
	for i := range ms {
		k := keyOf(ms[i])
		if _, dup := r.byKey[k]; dup {
			continue
		}
		id := "gen-" + strconv.Itoa(r.next)
		r.next++
		cp := ms[i]
		cp.ID = id
		r.byID[id] = cp
		r.byKey[k] = id
		ins++
	}
	return ins, nil
}

var _ ports.MovieRepository = (*memRepo)(nil)

func TestUsecase_NoMock(t *testing.T) {
	repo := newMemRepo()
	svc := NewMovieService(repo)

	ins, err := svc.EnsureSeed(context.Background(), []domain.Movie{
		{Title: "A", Year: 1999},
		{Title: "", Year: 2000},  // inválido
		{Title: "A", Year: 1999}, // duplicado
	})
	require.NoError(t, err)
	require.Equal(t, 1, ins)

	m, err := svc.Create(context.Background(), domain.Movie{Title: "B", Year: 2001})
	require.NoError(t, err)
	require.NotEmpty(t, m.ID)

	all, err := svc.List(context.Background())
	require.NoError(t, err)
	require.Len(t, all, 2)

	got, err := svc.Get(context.Background(), m.ID)
	require.NoError(t, err)
	require.Equal(t, "B", got.Title)

	require.NoError(t, svc.Delete(context.Background(), m.ID))
	_, err = svc.Get(context.Background(), m.ID)
	require.Error(t, err)
}

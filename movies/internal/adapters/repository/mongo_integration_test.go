//go:build integration
// +build integration

package repository

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/caiqueborghese/sipubtech-challenge/movies/internal/domain"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func newTestDB(t *testing.T) *mongo.Database {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}
	dbname := os.Getenv("MONGODB_DB")
	if dbname == "" {
		dbname = "moviesdb_test"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	t.Cleanup(cancel)

	cl, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	require.NoError(t, err)

	db := cl.Database(dbname)
	_ = db.Collection("movies").Drop(ctx)

	return db
}

func TestMongoRepository_Integration(t *testing.T) {
	db := newTestDB(t)
	repo, err := NewMongoRepository(db.Collection("movies"))
	require.NoError(t, err)

	ctx := context.Background()

	m1 := &domain.Movie{Title: "Legacy 8", Year: 1894, LegacyID: "8"}
	created, err := repo.Create(ctx, m1)
	require.NoError(t, err)
	require.NotEmpty(t, created.ID)

	got, err := repo.Get(ctx, created.ID)
	require.NoError(t, err)
	require.Equal(t, "Legacy 8", got.Title)

	got2, err := repo.Get(ctx, "8")
	require.NoError(t, err)
	require.Equal(t, "Legacy 8", got2.Title)

	cnt, err := repo.Count(ctx)
	require.NoError(t, err)
	require.Equal(t, int64(1), cnt)

	err = repo.Delete(ctx, created.ID)
	require.NoError(t, err)
	_, err = repo.Get(ctx, created.ID)
	require.Error(t, err)
}

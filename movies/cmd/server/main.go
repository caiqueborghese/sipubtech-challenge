package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/caiqueborghese/sipubtech-challenge/movies/internal/adapters/grpcserver"
	"github.com/caiqueborghese/sipubtech-challenge/movies/internal/adapters/repository"
	"github.com/caiqueborghese/sipubtech-challenge/movies/internal/seed"
	"github.com/caiqueborghese/sipubtech-challenge/movies/internal/usecase"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func main() {
	mongoURI := env("MONGODB_URI", "mongodb://mongo:27017/moviesdb")
	dbName := env("MONGODB_DB", "moviesdb")
	grpcAddr := ":" + env("GRPC_PORT", "50051")
	seedFile := os.Getenv("SEED_FILE")

	// conecta no Mongo com timeout e ping
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("mongo connect: %v", err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("mongo ping: %v", err)
	}

	col := client.Database(dbName).Collection("movies")

	// repositório mais caso de uso
	repo, err := repository.NewMongoRepository(col)
	if err != nil {
		log.Fatalf("new repo: %v", err)
	}
	svc := usecase.NewMovieService(repo)

	if seedFile != "" {
		movies, err := seed.LoadSeed(seedFile)
		if err != nil {
			log.Fatalf("load seed: %v", err)
		}
		inserted, err := svc.EnsureSeed(context.Background(), movies)
		if err != nil {
			log.Fatalf("ensure seed: %v", err)
		}
		if inserted > 0 {
			log.Printf("🌱 seeded %d movies", inserted)
		}
	}

	log.Printf(" gRPC listening on %s | mongo=%s db=%s", grpcAddr, mongoURI, dbName)
	if err := grpcserver.RunGRPCServer(svc, grpcAddr); err != nil {
		log.Fatal(err)
	}
}

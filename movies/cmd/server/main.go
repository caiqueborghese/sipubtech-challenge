package main

import (
	"context"
	"log"
	"os"
	"time"

	ae "github.com/caiqueborghese/sipubtech-challenge/movies/internal/adapters/events"
	"github.com/caiqueborghese/sipubtech-challenge/movies/internal/adapters/grpcserver"
	"github.com/caiqueborghese/sipubtech-challenge/movies/internal/adapters/repository"
	"github.com/caiqueborghese/sipubtech-challenge/movies/internal/ports"
	"github.com/caiqueborghese/sipubtech-challenge/movies/internal/seed"
	"github.com/caiqueborghese/sipubtech-challenge/movies/internal/usecase"
	"github.com/nats-io/nats.go"
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

	// ---- NATS (opcional) ----
	natsEnabled := env("NATS_ENABLED", "false")
	natsURL := env("NATS_URL", "nats://nats:4222")
	subjCreated := env("NATS_SUBJECT_CREATED", "movies.created")
	subjDeleted := env("NATS_SUBJECT_DELETED", "movies.deleted")

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

	// repositório
	repo, err := repository.NewMongoRepository(col)
	if err != nil {
		log.Fatalf("new repo: %v", err)
	}

	// publisher de eventos (pode ser nil)
	var pub ports.EventPublisher
	var nc *nats.Conn
	if natsEnabled == "true" {
		nc, err = nats.Connect(natsURL, nats.Name("movies-publisher"))
		if err != nil {
			log.Printf("NATS disabled (connect error): %v", err)
		} else {
			pub = ae.NewNatsPublisher(nc, subjCreated, subjDeleted)
			log.Printf("NATS connected at %s (created=%s, deleted=%s)", natsURL, subjCreated, subjDeleted)
			defer nc.Close()
		}
	}

	// serviço (com ou sem publisher)
	var svc ports.MovieService
	if pub != nil {
		svc = usecase.NewMovieServiceWithPublisher(repo, pub)
	} else {
		svc = usecase.NewMovieService(repo)
	}

	// seed opcional
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

	log.Printf("🎬 gRPC listening on %s | mongo=%s db=%s", grpcAddr, mongoURI, dbName)
	if err := grpcserver.RunGRPCServer(svc, grpcAddr); err != nil {
		log.Fatal(err)
	}
}

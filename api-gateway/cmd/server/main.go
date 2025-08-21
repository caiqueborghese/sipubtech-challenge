package main

import (
	"log"
	"os"

	docs "github.com/caiqueborghese/sipubtech-challenge/api-gateway/docs"
	"github.com/caiqueborghese/sipubtech-challenge/api-gateway/internal/handlers"
	"github.com/caiqueborghese/sipubtech-challenge/api-gateway/internal/usecase"
	moviespb "github.com/caiqueborghese/sipubtech-challenge/proto/moviespb"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Movies API Gateway
// @version 1.0
// @description Gateway REST que expõe os serviços de filmes via gRPC
// @host localhost:8080
// @BasePath /
func main() {
	addr := getenv("MOVIES_ADDR", "movies:50051")
	listen := getenv("HTTP_ADDR", ":8080")
	mode := getenv("GIN_MODE", "release")
	gin.SetMode(mode)

	// ---- Swagger runtime overrides (garante schemes/host/basePath) ----
	docs.SwaggerInfo.Title = "Movies API Gateway"
	docs.SwaggerInfo.Description = "Gateway REST que expõe os serviços de filmes via gRPC"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Host = getenv("SWAGGER_HOST", "localhost:8080")
	docs.SwaggerInfo.Schemes = []string{"http"}

	// gRPC client p/ o serviço Movies
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("dial movies gRPC (%s): %v", addr, err)
	}
	defer conn.Close()

	client := moviespb.NewMovieServiceClient(conn)
	movieSvc := usecase.NewMovieService(client)

	// HTTP (Gin)
	r := gin.Default()
	handlers.RegisterMovieRoutes(r, movieSvc)

	// Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.URL("/swagger/doc.json"),
		ginSwagger.DefaultModelsExpandDepth(-1), // desativa a aba "Models" (evita bug de render)
		ginSwagger.DocExpansion("none"),
	))

	log.Printf("HTTP listening on %s, talking to gRPC at %s", listen, addr)
	if err := r.Run(listen); err != nil {
		log.Fatal(err)
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

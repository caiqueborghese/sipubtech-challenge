package ports

// Usa go run para evitar depender do PATH local.
// Gera mocks de MovieRepository em internal/ports/mocks/ .
//
//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -destination=./mocks/mock_repository.go -package=mocks github.com/caiqueborghese/sipubtech-challenge/movies/internal/ports MovieRepository

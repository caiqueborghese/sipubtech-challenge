# Makefile — Sipub Tech Challenge
# Uso:
#   make up       -> gera gRPC, arruma deps e sobe API, Movies e Mongo
#   make down     -> derruba containers e volumes
#   make proto    -> (re)gera stubs gRPC
#   make test     -> roda testes dos serviços
#   make logs     -> segue logs
#   make tools    -> instala plugins gRPC e swag
#   make mocks    -> gera mocks (gomock) do módulo movies
#   make tidy     -> go mod tidy nos serviços
#   make clean    -> limpa artefatos gerados

APP_NAME := sipubtech-challenge
COMPOSE  := docker compose
GO       := go

# Caminhos do proto
PROTO_DIR  := proto
PROTO_FILE := $(PROTO_DIR)/moviespb/movies.proto

.PHONY: up down proto test logs tools tidy clean mocks

up: proto tidy
	$(COMPOSE) up -d --build
	@echo " Stack no ar. API: http://localhost:8080  | Swagger: http://localhost:8080/swagger/index.html"

down:
	$(COMPOSE) down -v
	@echo " Containers e volumes removidos."

proto: tools
	@echo "  Gerando stubs gRPC..."
	@# Injeta o bin dos plugins no PATH para evitar depender do shell do usuário
	@PLUGIN_BIN=$$(go env GOPATH)/bin; \
	if ! command -v protoc >/dev/null 2>&1; then \
	  echo " 'protoc' não encontrado. Instale o Protobuf Compiler (e.g. 'brew install protobuf') e rode novamente."; \
	  exit 1; \
	fi; \
	PATH="$$PLUGIN_BIN:$$PATH" protoc \
	  -I $(PROTO_DIR) \
	  --go_out=$(PROTO_DIR) --go_opt=paths=source_relative \
	  --go-grpc_out=$(PROTO_DIR) --go-grpc_opt=paths=source_relative \
	  $(PROTO_FILE)
	@echo "Stubs gerados em $(PROTO_DIR)"

test:
	@echo "Rodando testes (movies + api-gateway)..."
	cd movies && $(GO) test ./... -count=1 -v
	@if [ -d "api-gateway" ]; then cd api-gateway && $(GO) test ./... -count=1 -v; fi
	@echo "Testes finalizados."

logs:
	$(COMPOSE) logs -f

tidy:
	@echo "go mod tidy (movies, api-gateway se existir)..."
	cd movies && $(GO) mod tidy
	@if [ -d "api-gateway" ]; then cd api-gateway && $(GO) mod tidy; fi

tools:
	@echo "Verificando/instalando ferramentas (protoc-gen-go, protoc-gen-go-grpc, swag)..."
	@if ! command -v protoc >/dev/null 2>&1 ; then \
	  echo "'protoc' não encontrado. Instale com: brew install protobuf (macOS)"; \
	fi
	@if ! command -v protoc-gen-go >/dev/null 2>&1 ; then \
	  echo "Instalando protoc-gen-go..."; \
	  $(GO) install google.golang.org/protobuf/cmd/protoc-gen-go@latest; \
	fi
	@if ! command -v protoc-gen-go-grpc >/dev/null 2>&1 ; then \
	  echo "Instalando protoc-gen-go-grpc..."; \
	  $(GO) install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest; \
	fi
	@if ! command -v swag >/dev/null 2>&1 ; then \
	  echo "Instalando swag (Swagger CLI)..."; \
	  $(GO) install github.com/swaggo/swag/cmd/swag@latest; \
	fi
	@echo "Ferramentas OK."

mocks:
	@echo "Gerando mocks (gomock) do módulo movies..."
	cd movies && mkdir -p internal/ports/mocks && \
		go generate ./internal/ports
	@echo "Mocks gerados em movies/internal/ports/mocks"

clean: down
	@echo "Limpando artefatos gerados..."
	@find $(PROTO_DIR) -maxdepth 1 -type f -name "*.pb.go" -delete || true

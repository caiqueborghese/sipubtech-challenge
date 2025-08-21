# Sipub Tech Challenge — Movies

Microserviço de catálogo de filmes com **arquitetura hexagonal**, split em **API Gateway (HTTP/REST)** e **Movies (gRPC)** com **MongoDB** como storage.  
A comunicação entre os serviços é feita via **gRPC/Protobuf**. A API HTTP é documentada com **Swagger**.


---

## 🔧 Tecnologias

- Go (módulos: `api-gateway`, `movies`, `proto`)
- gRPC + Protobuf
- MongoDB
- Docker & Docker Compose
- Swagger (swaggo / gin-swagger)
- Testes com `go test` (mocks com `gomock/testify`)

---

## 📂 Arquitetura de Pastas

```
.
├─ api-gateway/
│  ├─ cmd/server/
│  │  └─ main.go                 # inicialização HTTP + Swagger
│  ├─ internal/
│  │  ├─ adapters/               # ADAPTADORES (saída) do gateway
│  │  │  └─ grpcclient/
│  │  │     └─ movies_client.go  # client gRPC p/ serviço Movies (porta de saída)
│  │  ├─ domain/                 # entidades expostas no gateway (shape HTTP)
│  │  │  └─ movie.go
│  │  ├─ handlers/               # ADAPTADORES (entrada) HTTP
│  │  │  └─ movie_handler.go     # HTTP <-> usecase, mapeia erros p/ HTTP
│  │  ├─ ports/                  # PORTAS do gateway
│  │  │  ├─ service.go           # porta de ENTRADA (MovieService)
│  │  │  └─ movies_client.go     # porta de SAÍDA (client gRPC)
│  │  └─ usecase/                # CASOS DE USO do gateway
│  │     ├─ movie_service.go
│  └─ docs/                      # artefatos gerados do Swagger
│     ├─ docs.go
│     ├─ swagger.json
│     └─ swagger.yaml
│
├─ movies/
│  ├─ cmd/server/
│  │  └─ main.go                 # bootstrap gRPC + seed
│  ├─ internal/
│  │  ├─ adapters/
│  │  │  ├─ grpcserver/          # ADAPTADOR (entrada) — server gRPC
│  │  │  │  └─ server.go         # implementa protobuf, expõe legacy_id
│  │  │  └─ repository/          # ADAPTADOR (saída) — MongoDB
│  │  │     └─ mongo.go          # índices, mapeamento domain <-> BSON, erros 11000
│  │  ├─ domain/                 # REGRAS de domínio (entidades/validação/erros)
│  │  │  └─ movie.go
│  │  ├─ ports/                  # PORTAS do domínio Movies
│  │  │  ├─ repository.go        # porta de SAÍDA (persistência)
│  │  │  ├─ service.go           # porta de ENTRADA (serviço de aplicação)
│  │  │  └─ mocks/               # mocks gerados (gomock)
│  │  │     └─ mock_repository.go
│  │  ├─ seed/                   # carregamento do movies.json -> LegacyID
│  │  │  └─ seed.go
│  │  └─ usecase/                # CASOS DE USO do serviço Movies
│  │     ├─ movie_service.go
│  └─ seed/
│     └─ movies.json             # arquivo de seed fornecido
│
├─ proto/
│  └─ moviespb/                  # contrato gRPC
│     ├─ movies.proto
│     ├─ movies.pb.go
│     └─ movies_grpc.pb.go
│
├─ deploy/
│  └─ docker/
│     ├─ api.Dockerfile
│     └─ movies.Dockerfile
│
├─ go.work                        # workspace Go (3 módulos)
├─ docker-compose.yml
├─ Makefile
└─ README.md
```

### Comentários rápidos (Hexagonal)

- **API Gateway**
  - `handlers`: adaptador de **entrada** HTTP (Gin). Converte HTTP ↔ domínio do *gateway*; chama `usecase.MovieService`.
  - `usecase`: orquestra o client gRPC; validações do lado HTTP.
  - `ports`: interfaces do gateway. `service.go` (entrada) e `movies_client.go` (saída).
  - `adapters/grpcclient`: implementa porta de saída chamando o gRPC de `movies`.

- **Movies**
  - `adapters/grpcserver`: adaptador de **entrada** gRPC. Implementa protobuf, traduz pb ↔ domínio e chama `usecase`.
  - `adapters/repository`: adaptador de **saída** (MongoDB). Índices únicos: **`uniq_title_year`** e **`uniq_legacy_id`**.
  - `usecase`: regras de negócio e `EnsureSeed` (idempotente).
  - `ports`: interfaces centrais (`MovieService`, `MovieRepository`) — facilitam mocks e troca de storage.
  - `domain`: entidade `Movie`, normalização e erros.
  - `seed`: lê `seed/movies.json`, preenche `LegacyID` e dispara `EnsureSeed`.

---

## ▶️ Subir tudo com **um comando**

**Pré-requisitos**
- Docker + Docker Compose
- Go 1.25+ (apenas se for rodar local sem Docker)

**Com Make (recomendado)**
```bash
make up
# API: http://localhost:8080
# Swagger: http://localhost:8080/swagger/index.html
```

**Com Docker Compose**
```bash
docker compose up -d --build
```

**Serviços**
- **API Gateway**: http://localhost:8080  
- **Swagger**: http://localhost:8080/swagger/index.html  
- **Movies (gRPC)**: interno em `movies:50051` (mapeado em `localhost:50051` para debug)  
- **MongoDB**: `localhost:27017`

---

## 🔌 Variáveis de Ambiente

| Serviço       | Variável        | Padrão                                 | Descrição                               |
|---------------|-----------------|----------------------------------------|-----------------------------------------|
| api-gateway   | `MOVIES_ADDR`   | `movies:50051`                         | Endereço do gRPC do serviço `movies`    |
| api-gateway   | `HTTP_ADDR`     | `:8080`                                | Porta HTTP                              |
| api-gateway   | `SWAGGER_HOST`  | `localhost:8080`                       | Host do Swagger (override runtime)      |
| movies        | `MONGODB_URI`   | `mongodb://mongo:27017/moviesdb`       | URI do Mongo                            |
| movies        | `MONGODB_DB`    | `moviesdb`                             | Nome do banco                           |
| movies        | `GRPC_PORT`     | `50051`                                | Porta gRPC                              |
| movies        | `SEED_FILE`     | `/app/seed/movies.json`                | Caminho do seed (habilita seed)         |

---

## 📘 Swagger

- UI: **http://localhost:8080/swagger/index.html**  
- JSON: **http://localhost:8080/swagger/doc.json**

> A listagem `GET /movies` aceita `?limit=` (default 50, máx 200) para não travar a UI com payload gigante.

---

## 🧭 Rotas HTTP

### `GET /movies?limit=50`
Lista os filmes; ordenação padrão por `title` no Mongo.  
**Exemplos**
```bash
curl -s "http://localhost:8080/movies" | jq .
curl -s "http://localhost:8080/movies?limit=5" | jq '.|length'
```

**Modelo de resposta**
```json
[
  {
    "id": "8",
    "title": "Edison Kinetoscopic Record of a Sneeze (1894)",
    "year": 1894
  }
]
```

---

### `GET /movies/{id}`
Busca por **ID externo** (`legacy_id` do JSON). Também aceita `_id` (ObjectID) dos itens criados via `POST`.

```bash
# pelo legacy_id do seed
curl -s http://localhost:8080/movies/8 | jq .

# pelo ObjectID de um POST:
# curl -s http://localhost:8080/movies/68a60b2b457c7c8d2c09d81f | jq .
```

---

### `POST /movies`
Cria um filme. **Não envie `id` no corpo**.

```bash
curl -s -X POST http://localhost:8080/movies   -H "Content-Type: application/json"   -d '{"title":"Meu Filme de Teste","year":2025}' | jq .
```

**Respostas**
- `201 Created`
- `400 invalid body`
- `409 movie already exists (title+year)` (se mapeado) ou `502` com detalhe do Mongo

---

### `DELETE /movies/{id}`
Remove por **ID externo** (`legacy_id`) ou ObjectID.

```bash
ID=$(curl -s -X POST http://localhost:8080/movies   -H "Content-Type: application/json"   -d '{"title":"Apagar Depois","year":2026}' | jq -r '.id')

curl -i -X DELETE "http://localhost:8080/movies/$ID"
```

**Respostas**
- `204 No Content`
- `404 movie not found`

---

## 🌱 Seed — popular / resetar banco

**Reset rápido (drop + reseed)**
```bash
docker compose exec mongo   mongosh "mongodb://localhost:27017/moviesdb" --quiet   --eval 'db.movies.drop()'

docker compose restart movies
```

**Reset total**
```bash
docker compose down -v
docker compose up -d --build
```

---

## 🧪 Testes

Rode tudo:
```bash
make test
```

Ou por módulo:

```bash
# movies (usecase + grpcserver)
cd movies
go test ./... -count=1 -v

# api-gateway (handlers HTTP + usecase)
cd ../api-gateway
go test ./... -count=1 -v
```

**Mocks** (no módulo `movies`):
```bash
make mocks
```

> Há testes **com mocks** (ports/mock) e **sem mocks** (ex.: grpc com `bufconn`).  
> Teste de repositório Mongo pode ser habilitado com build tag `integration` (opcional), apontando `MONGODB_URI` para um Mongo local de testes.

---


# ---------- builder ----------
FROM golang:1.25-alpine AS builder
WORKDIR /src
RUN apk add --no-cache ca-certificates


COPY movies/go.mod movies/go.sum ./movies/
COPY proto/go.mod proto/go.sum ./proto/
RUN cd movies && go mod download


COPY proto ./proto
COPY movies ./movies


ENV CGO_ENABLED=0
ARG TARGETOS=linux
ARG TARGETARCH=arm64
RUN cd movies/cmd/server && \
    GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /out/movies-server .

# ---------- runtime (ALPINE) ----------
FROM alpine:3.20
WORKDIR /app

RUN apk add --no-cache ca-certificates bash curl busybox-extras
COPY --from=builder /out/movies-server /app/movies-server

COPY movies/seed/movies.json /app/seed/movies.json


RUN adduser -D -H app && chown app:app /app/movies-server
USER app

ENV GRPC_PORT=50051
EXPOSE 50051
ENTRYPOINT ["/app/movies-server"]

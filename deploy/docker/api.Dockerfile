# ---------- builder ----------
FROM golang:1.25-alpine AS builder
WORKDIR /src
RUN apk add --no-cache ca-certificates


COPY api-gateway/go.mod api-gateway/go.sum ./api-gateway/
COPY proto/go.mod proto/go.sum ./proto/
RUN cd api-gateway && go mod download


COPY . .


ARG TARGETOS
ARG TARGETARCH
RUN cd api-gateway/cmd/server && \
    CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-arm64} \
    go build -trimpath -ldflags="-s -w" -o /out/api-gateway

# ---------- runtime ----------
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /out/api-gateway /app/api-gateway
RUN adduser -D -H app && chown app:app /app/api-gateway
USER app
EXPOSE 8080
ENTRYPOINT ["/app/api-gateway"]




# syntax=docker/dockerfile:1

FROM golang:1.24 AS builder

WORKDIR /src

COPY go.mod go.sum ./
COPY auth-service/go.mod auth-service/go.sum ./auth-service/

# Modules cachs and build-cache with BuildKit (faster second build)
RUN --mount=type=cache,target=/go/pkg/mod \
    cd auth-service && go mod download

COPY pkg ./pkg
COPY auth-service/. ./auth-service

WORKDIR /src/auth-service

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    mkdir -p /out && \
    CGO_ENABLED=1 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} \
    go build -trimpath -ldflags="-s -w" -o /out/auth-service ./cmd/server

# ---------- runtime stage ----------
FROM debian:stable-slim

RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates librdkafka1 && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /out/auth-service ./auth-service
RUN chmod +x /app/auth-service

EXPOSE 8000

ENTRYPOINT ["/app/auth-service"]
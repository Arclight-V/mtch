# syntax=docker/dockerfile:1

FROM golang:1.24 AS builder

WORKDIR /src

COPY go.mod go.sum ./
COPY user-service/go.mod user-service/go.sum ./user-service/

# Modules cachs and build-cache with BuildKit (faster second build)
RUN --mount=type=cache,target=/go/pkg/mod \
    cd user-service && go mod download

COPY pkg ./pkg
COPY user-service/. ./user-service

WORKDIR /src/user-service

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    mkdir -p /out && \
    CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} \
    go build -trimpath -ldflags="-s -w" -o /out/user-service ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app
COPY --from=builder /out/user-service ./user-service

EXPOSE 50051

USER nonroot:nonroot

ENTRYPOINT ["/app/user-service"]
# syntax=docker/dockerfile:1

FROM golang:1.24 AS builder

WORKDIR /src

COPY go.mod go.sum ./
COPY notification/go.mod notification/go.sum ./notification/

# Modules cachs and build-cache with BuildKit (faster second build)
RUN --mount=type=cache,target=/go/pkg/mod \
    cd notification && go mod download

COPY pkg ./pkg
COPY notification/. ./notification

WORKDIR /src/notification

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    mkdir -p /out && \
    CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} \
    go build -trimpath -ldflags="-s -w" -o /out/notification ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app
COPY --from=builder /src/notification/config/flags.flagd.json ./flags.flagd.json
COPY --from=builder /out/notification ./notification

EXPOSE 50052

USER nonroot:nonroot

ENTRYPOINT ["/app/notification"]
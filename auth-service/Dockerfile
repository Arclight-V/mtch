# syntax=docker/dockerfile:1

FROM golang:1.24 AS builder

WORKDIR /src

COPY pkg ./pkg
COPY go.mod go.sum ./
COPY auth-service/go.mod auth-service/go.sum ./auth-service/

WORKDIR /src/auth-service
RUN go mod download

COPY auth-service/. .
RUN go mod tidy

RUN mkdir -p /out && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/auth-service ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app
COPY --from=builder /out/auth-service ./auth-service

EXPOSE 8000

USER nonroot:nonroot

ENTRYPOINT ["/app/auth-service"]
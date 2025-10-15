# syntax=docker/dockerfile:1

FROM golang:1.24 AS builder

WORKDIR /src

COPY proto ./proto
COPY user-service/go.mod user-service/go.sum ./user-service/

WORKDIR /src/user-service
RUN go mod download

COPY user-service/. ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/userservice-service ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app
COPY --from=builder /out/user-service ./user-service

EXPOSE 50051

USER nonroot:nonroot

ENTRYPOINT ["/app/user-service"]
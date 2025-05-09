.PHONY: build run test

build:
	go build -o auth-service cmd/server/main.go

run:
	./auth-service

test:
	./... -race -cover
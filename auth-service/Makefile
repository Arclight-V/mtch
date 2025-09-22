.PHONY: build run test clean mocks swag docker

BIN_DIR := bin
SERVICE := auth-service
VERSION ?= latest
AUTH_IMAGE ?= ${SERVICE}:${VERSION}

build: swag
	go build -o $(BIN_DIR)/$(SERVICE) cmd/server/main.go

run:
	./$(BIN_DIR)/$(SERVICE)

test:
	go test -race -covermode=atomic -coverprofile=coverage.out ./...

clean:
	rm -rf $(BIN_DIR)

mocks:
	go generate ./...

mocks-check: mocks
	git diff --exit-code

swag:
	swag init -g cmd/server/main.go -o ./docs

docker-build:
	docker build --tag mtch-${AUTH_IMAGE} -f Dockerfile .
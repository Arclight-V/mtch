.PHONY: build run test test-integration clean docker

BIN_DIR := bin
SERVICE := user-service
VERSION ?= latest
USER_IMAGE ?= ${SERVICE}:${VERSION}

build:
	go build -o $(BIN_DIR)/$(SERVICE) cmd/server/main.go

run:
	./$(BIN_DIR)/$(SERVICE)

test:
	go test ./... -race -cover

test-integration:
	go test -tags=integration ./... -race -cover

clean:
	rm -rf $(BIN_DIR)

mocks:
	go generate ./...

docker-build:
	docker build --tag mtch-${USER_IMAGE} -f Dockerfile .
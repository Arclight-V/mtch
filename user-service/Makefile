.PHONY: build run test clean docker

BIN_DIR := bin
SERVICE := user-service

build:
	go build -o $(BIN_DIR)/$(SERVICE) cmd/server/main.go

run:
	./$(BIN_DIR)/$(SERVICE)

test:
	./... -race -cover

clean:
	rm -rf $(BIN_DIR)
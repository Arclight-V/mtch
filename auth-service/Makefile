.PHONY: build run test clean

BIN_DIR := bin
SERVICE := auth-service

build:
	go build -o $(BIN_DIR)/$(SERVICE) cmd/server/main.go

run:
	.$(BIN_DIR)/$(SERVICE)

test:
	./... -race -cover

clean:
	rm -rf $(BIN_DIR)
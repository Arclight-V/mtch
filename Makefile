.PHONY: proto build  run test docker

TAG ?= latest
PREFIX ?= mtch-
SERVICES = auth-service user-service

proto:
	protoc	-I proto \
          	--go_out=proto	--go_opt=paths=source_relative \
          	--go-grpc_out=proto     --go-grpc_opt=paths=source_relative\
            proto/user.proto

build:
	#$(MAKE) proto
	$(MAKE) -C auth-service build
	$(MAKE) -C user-service build

run:
	$(MAKE) -C auth-service run
	$(MAKE) -C user-service run

test:
	$(MAKE) -C auth-service test
	$(MAKE) -C user-service test

docker-build: $(SERVICES:%=docker-build-%)

docker-build-%:
	docker build -t $(PREFIX)$*:$(TAG) -f $*/Dockerfile .
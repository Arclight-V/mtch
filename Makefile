.PHONY: proto build  run test docker

TAG ?= latest
PREFIX ?= mtch-
SERVICES ?= auth-service user-service

# -------- Help --------
.PHONY: help
help:
	@echo "make proto        					- Generate protobuf stubs"
	@echo "make build        					- Build all services"
	@echo "make build-auth-service    			- Build one service (SERVICES=auth-service)"
	@echo "make test         					- Run all unit tests"
	@echo "make test-auth-service         		- Run one service unit tests (SERVICES=auth-service)"
	@echo "make run          					- Run all services (It's doesn't work)"
	@echo "make run-auth-service          		- Run one service locally (SERVICE=auth-service)"
	@echo "make docker-build					- Build all docker images"
	@echo "make docker-build-auth-service       - Build one docker image (SERVICE=auth-service)"



# -------- Proto --------
proto:
	protoc	-I proto \
          	--go_out=proto			--go_opt=paths=source_relative \
          	--go-grpc_out=proto     --go-grpc_opt=paths=source_relative\
            proto/user.proto

# -------- Build --------
build: $(SERVICES:%=build-%)

build-%:
	echo "==> building $*"
	$(MAKE) -C $* build

# -------- Run --------
run: $(SERVICES:%=run-%)

run-%:
	$(MAKE) -C $* run

# -------- Test --------
test: $(SERVICES:%=test-%)

test-%:
	$(MAKE) -C $* test

# -------- Docker build --------
docker-build: $(SERVICES:%=docker-build-%)

docker-build-%:
	docker build -t $(PREFIX)$*:$(TAG) -f $*/Dockerfile .
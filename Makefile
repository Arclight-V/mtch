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
	@echo "make compose-up						- Build & start all containers"
	@echo "make compose-down					- Stop and remove containers"
	@echo "make compose-build					- Rebuild images without cache"
	@echo "make compose-restart					- Restart the whole stack"
	@echo "make compose-logs					- Tail logs of all services"
	@echo "make compose-ps						- Show container statuses"
	@echo "make compose-rebuild-auth / make compose-rebuild-user	Rebuild and restart one service"
	@echo "make compose-up-observability	    - Start Grafana + Prometheus + Loki + Promtail"
	@echo "make compose-open-grafana / make open-prom	-Open Grafana or Prometheus in browser"
	@echo "make compose-prune					- Remove all unused Docker data"
	@echo "make compose-clean					- Tear down everything including volumes"



# -------- Proto --------
proto:
	protoc	-I . --go_out=. --go_opt=paths=source_relative \
          	--go-grpc_out=. --go-grpc_opt=paths=source_relative\
            pkg/userservice/userservicepb/v1/userservice.proto

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


#-------- Docker Compose helpers --------

##### Docker Compose helpers ----------------------------------------------------

# Project name in Docker Compose (containers will have mtch-* prefix)
PROJECT := mtch

# Short aliases
COMPOSE  := docker compose -p $(PROJECT)

# Services defined in docker-compose.yml
SERVICES := auth user prometheus grafana loki promtail

.PHONY: compose-up compose-down compose-build compose-restart compose-logs compose-ps \
        compose-up-auth compose-up-user compose-up-observability \
        compose-rebuild-auth compose-rebuild-user \
        compose-logs-auth compose-logs-user compose-logs-prom compose-logs-graf \
        compose-shell-auth compose-shell-user compose-shell-net \
        open-grafana open-prom targets \
        compose-reload-prom compose-prune compose-clean

## Build and start all services in detached mode
compose-up:
	$(COMPOSE) up -d --build

## Stop and remove all containers (keep volumes)
compose-down:
	$(COMPOSE) down

## Rebuild all images without using cache (no start)
compose-build:
	$(COMPOSE) build --no-cache

## Restart all services
compose-restart:
	$(COMPOSE) down && $(COMPOSE) up -d

## Show container status
compose-ps:
	$(COMPOSE) ps

## View logs from all services (last 100 lines, follow mode)
compose-logs:
	$(COMPOSE) logs -f --tail=100

# --- Service-specific operations ----------------------------------------------

## Start only the auth service
compose-up-auth:
	$(COMPOSE) up -d --build auth

## Start only the user service
compose-up-user:
	$(COMPOSE) up -d --build user

## Start observability stack (Prometheus + Grafana + Loki + Promtail)
compose-up-observability:
	$(COMPOSE) up -d --build prometheus grafana loki promtail

## Rebuild and restart a single service
compose-rebuild-auth:
	$(COMPOSE) build auth && $(COMPOSE) up -d auth

compose-rebuild-user:
	$(COMPOSE) build user && $(COMPOSE) up -d user

## View logs per service
compose-logs-auth:
	$(COMPOSE) logs -f --tail=200 auth

compose-logs-user:
	$(COMPOSE) logs -f --tail=200 user

compose-logs-prom:
	$(COMPOSE) logs -f --tail=200 prometheus

compose-logs-graf:
	$(COMPOSE) logs -f --tail=200 grafana

# --- Interactive / debugging ---------------------------------------------------

# Note: distroless images have no /bin/sh inside.
# For interactive debugging, use a temporary helper container.
compose-shell-auth:
	@echo "auth is built on distroless — no shell inside. Use 'make shell-net' and connect to auth via network."

compose-shell-user:
	@echo "user is built on distroless — no shell inside. Use 'make shell-net' and connect to user via network."

## Run a temporary Alpine shell container inside the same network (for debugging)
compose-shell-net:
	docker run --rm -it --network $$(basename $$(pwd))_default alpine:3.20 sh

# --- Quick links / monitoring --------------------------------------------------

compose-open-grafana:
	@python3 -c "import webbrowser; webbrowser.open('http://localhost:9090')"

open-prom:
	@python3 -c "import webbrowser; webbrowser.open('http://localhost:9090')"

## Display Prometheus targets (should show UP)
targets:
	@echo "Open http://localhost:9090/targets"

## Soft reload Prometheus configuration
compose-reload-prom:
	$(COMPOSE) exec prometheus wget -qO- http://localhost:9090/-/reload || true

# --- Cleanup -------------------------------------------------------------------

## Remove all unused Docker data (images, cache, networks, volumes)
compose-prune:
	docker system prune -a --volumes -f

## Completely remove the compose project (including volumes)
compose-clean:
	$(COMPOSE) down -v


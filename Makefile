include .mk/migrate.mk

.PHONY: up down seed mock test-unit test-integration test-e2e test-web test-web-e2e test

COMPOSE ?= podman compose

up:
	$(COMPOSE) up --build -d

down:
	$(COMPOSE) down

seed:
	cd apps/api && SEED=1 go run ./cmd/seed

mock:
	cd apps/api && mockery

test-unit:
	cd apps/api && go test -tags=unit ./...

test-integration:
	cd apps/api && go test -tags=integration ./internal/adapter/postgres/... -count=1

test-e2e:
	cd apps/api && go test -tags=e2e ./tests/e2e/... -count=1

test-web:
	cd apps/web && npm test -- --run

test-web-e2e:
	cd apps/web && npm run test:e2e

test: test-unit test-web

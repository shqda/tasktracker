.PHONY: gen test utest itest migration-up migration-down migration-status migration-create

## Migrations
POSTGRES_HOST ?= localhost
POSTGRES_PORT ?= 5432
POSTGRES_USER ?= postgres
POSTGRES_PASSWORD ?= postgres
POSTGRES_DB ?= tasktracker

export GOOSE_DRIVER := postgres
export GOOSE_DBSTRING := host=$(POSTGRES_HOST) port=$(POSTGRES_PORT) user=$(POSTGRES_USER) password=$(POSTGRES_PASSWORD) dbname=$(POSTGRES_DB) sslmode=disable
export GOOSE_MIGRATION_DIR := ./migrations

## Generates all mocks using mockery
gen:
	go run github.com/vektra/mockery/v2@latest

## Run all test
test:
	go test -count=1 ./...

## Run unit-test
utest:
	go test ./... --short

## Run integration-tests
itest:
	go test ./internal/storage/postgres/... -v

## Apply all migrations
migration-up:
	goose up

## Rollback all the migrations
migration-down:
	goose down

## Show migration status
migration-status:
	goose status

## Create new migration file. Usage: make migration-create NAME=...
migration-create:
ifndef NAME
	$(error NAME is required, e.g. make migration-create NAME=add_users)
endif
	goose create $(NAME) sql



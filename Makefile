-include .env
export

CURRENT_DIR=$(shell pwd)

# run service
.PHONY: run
run:
	go run cmd/app/main.go

# go generate

proto-gen:
	./scripts/gen_proto.sh

# migrate
.PHONY: migrate
migrate:
	migrate -source file://db/migrations -database postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable up

DB_URL := "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable"

migrate-up:
	migrate -path db/migrations -database $(DB_URL) -verbose up

migrate-down:
	migrate -path db/migrations -database $(DB_URL) -verbose down

migrate-force:
	migrate -path migrations -database $(DB_URL) -verbose force 1

migrate-file:
	migrate create -ext sql -dir db/migrations/ -seq url_shortern


pull-proto-module:
	git submodule update --init --recursive

update-proto-module:
	git submodule update --remote --merge
swag-init:
	swag init -g internal/app/app.go -o api/openapi  --outputTypes "go,json,yaml" --overridesFile .swaggo
.PHONY: sqlc
sqlc:
	sqlc generate
	go run tools/sqlc-gen-custom/main.go
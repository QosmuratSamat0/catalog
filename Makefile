
run:
	go run cmd/catalog/main.go --config=./config/local.yaml

generate:
	protoc -I proto \
		proto/catalog/*.proto \
		--go_out=./proto/gen/go/ --go_opt=paths=source_relative \
		--go-grpc_out=./proto/gen/go/ --go_opt=paths=source_relative

include .env
export $(shell sed 's/=.*//' .env)

migrate:
	@if echo $(DB_URL) | grep -q '?'; then \
		STORAGE="$(DB_URL)&sslmode=disable"; \
	else \
		STORAGE="$(DB_URL)?sslmode=disable"; \
	fi; \
	go run ./cmd/migrator --storage-path="$$STORAGE" --migrations-path="./migrations"
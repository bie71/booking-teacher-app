.PHONY: run build tidy migrate

migrate:
	go run migrations/migrate.go

run:
	go run cmd/app/main.go

build:
	go build -o bin/app cmd/app/main.go

tidy:
	go mod tidy

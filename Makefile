.PHONY: run build test lint format tidy compose-up compose-down

run:
	go run .

build:
	go build -o bin/server .

test:
	go test ./... -cover

format:
	golangci-lint fmt

lint:
	golangci-lint run

tidy:
	go mod tidy

compose-up:
	docker compose up -d --wait --remove-orphans

compose-down:
	docker compose down

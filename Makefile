.PHONY: run build image container test lint format tidy compose-up compose-down

IMAGE ?= sgt-noti-service
PORT ?= 8088

run:
	go run .

build:
	go build -o bin/server .

image:
	docker build -t $(IMAGE) .

container: image
	docker run --rm --env-file .env -p $(PORT):8080 $(IMAGE)

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

.PHONY: build run test clean docker-build docker-run

APP_NAME=polkadot-rewards-api

build:
	go build -o bin/$(APP_NAME) .

run:
	go run .

test:
	go test ./... -v

docker-build:
	docker build -t $(APP_NAME) .

docker-run:
	docker run -p 8080:8080 $(APP_NAME)
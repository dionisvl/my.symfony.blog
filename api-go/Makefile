.PHONY: build test lint run

build:
	go build -o bin/api-go ./cmd/api

run:
	go run ./cmd/api

test:
	go test ./... -v -count=1

lint:
	golangci-lint run ./...
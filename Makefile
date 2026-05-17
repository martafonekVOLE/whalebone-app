.PHONY: run

run:
	docker compose up -d --build

stop:
	docker compose down

build-local:
	go build ./cmd/simple-microservice/main.go

run-local:
	go run ./cmd/simple-microservice/main.go

test:
	go test -count=1 -race -cover ./...

lint:
	golangci-lint run --fix
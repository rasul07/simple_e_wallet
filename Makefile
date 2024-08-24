.PHONY: swag test docker-build docker-up docker-down docker-logs

swag:
	swag init -g internal/handlers/api.go -o api/docs

test:
	go test ./internal/service

build:
	docker-compose build

run:
	docker-compose up -d

stop:
	docker-compose down

logs:
	docker-compose logs -f app
swag:
	swag init -g internal/handlers/api.go -o api/docs

test:
	go test ./internal/service
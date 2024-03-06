run:
	go run cmd/core-service/main.go -dev-logger

migrate:
	go run cmd/form-service/main.go -migrate
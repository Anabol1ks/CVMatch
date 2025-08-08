run:
	go run cmd/main.go

swag:
	swag init -g cmd/main.go

doc:
	docker-compose up -d --build

test:
	go test ./...
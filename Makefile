run:
	cd srcs/ && go run main.go

up:
	docker-compose up --build

down:
	docker-compose down

tidy:
	cd srcs/ && go mod tidy

.PHONY: run up down tidy

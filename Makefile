run:
	cd srcs/ && go run main.go

up:
	docker-compose up --build

start:
	docker-compose start
stop:
	docker-compose stop

down:
	docker-compose down

tidy:
	cd srcs/ && go mod tidy

.PHONY: run up down tidy

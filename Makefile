run:
	cd srcs/ && go run main.go

watch:
	./watch.sh
#	while true; do find srcs/ -name "*.go" | entr -r make run; sleep(1); done

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

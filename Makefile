all: docker-stop docker-build docker-run

format: 
	go fmt ./...

lint:
	golangci-lint run --config ./.golangci.yml


docker-run:
	docker compose up

docker-stop:
	docker compose down

docker-build:
	docker compose build --no-cache

docker-list:
	docker image ls

docker-ps:
	docker ps

format: 
	go fmt ./...

lint:
	golangci-lint run --config ./.golangci.yml

docker-run:
	docker compose up

docker-stop:
	docker compose stop

docker-list:
	docker image ls

docker-ps:
	docker ps

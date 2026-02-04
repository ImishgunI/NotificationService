format: 
	go fmt ./...

lint:
	golangci-lint run --config ./.golangci.yml

docker-run:
	docker compose up

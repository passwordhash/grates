
prod: docker-build-prod docker-up-prod timeout migrate-up

dev: docker-build-dev docker-up-dev timeout migrate-up go-run

migrate-up:
	docker compose up migrate

timeout:
	sleep 4

docker-build-dev:
	ENV_FILE_NAME=.env docker compose build db rdb migrate

docker-build-prod:
	ENV_FILE_NAME=.prod.env docker compose build

docker-up-prod:
	ENV_FILE_NAME=.prod.env docker compose up -d

docker-up-dev:
	ENV_FILE_NAME=.env docker compose up -d db rdb migrate

go-run:
	go run ./cmd/http/main.go
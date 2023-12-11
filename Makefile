
prod: docker-build-prod docker-up-prod timeout migrate-prod

dev: docker-build-dev docker-up timeout migrate-prod

migrate-prod:
	docker compose up migrate

timeout:
	sleep 4

docker-build-dev:
	ENV_FILE_NAME=.env docker compose build

docker-build-prod:
	ENV_FILE_NAME=.prod.env docker compose build

docker-up-prod:
	ENV_FILE_NAME=.prod.env docker compose up -d

docker-up:
	docker compose up -d


prod: docker-build-prod docker-up

dev: docker-build-dev docker-up

docker-build-dev:
	ENV_FILE_NAME=.env docker compose build

docker-build-prod:
	ENV_FILE_NAME=.prod.env docker compose build

docker-up:
	docker compose up -d

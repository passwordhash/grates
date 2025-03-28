
prod: docker-build-prod docker-up-prod timeout migrate-up


prod-reload:

dev: docker-build-dev docker-up-dev timeout migrate-up go-run

db-reload: db-down db-up

# ==================================================================================================

migrate-up:
	docker compose up migrate

timeout:
	sleep 4


docker-build-dev:
	docker compose build db rdb migrate

docker-build-prod:
	CONFIG_FILE_NAME=config.prod docker compose build

docker-up-dev:
	docker compose up -d db rdb migrate

docker-up-prod:
	CONFIG_FILE_NAME=config.prod docker compose up -d



go-run:
	go run ./cmd/http/main.go


db-up:
	migrate -path ./schema -database 'postgres://postgres:root@localhost:54320/postgres?sslmode=disable' up

db-down:
	migrate -path ./schema -database 'postgres://postgres:root@localhost:54320/postgres?sslmode=disable' down

generate-docs:
	swag init --output ./docs --generalInfo ./cmd/http/main.go

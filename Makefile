
setup: docker-build up-dbs timeout up-prod

# Костыль
timeout:
	sleep 5
	echo "hello"

docker-build:
	docker compose build db rdb migrate app

docker-build-prod:
	docker compose build db rdb migrate app-prod

up-dbs:
	docker compose up db rdb -d

# TODO: решить как передать флаг в golang
up-prod:
	docker compose up app-prod migrate

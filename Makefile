
prod: export-prod docker-build docker-up

#timeout:
#	sleep 5
#	echo "hello"

export-prod:
	export ENV_FILE_NAME=.prod.env

docker-build:
	docker compose build

docker-up:
	docker compose up -d
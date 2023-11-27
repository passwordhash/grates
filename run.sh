#!/bin/bash

for t in $@;do
    if [ $t == "-n" ];then
        is_run=0
    fi
done

check() {
    local code=$?
    local msg=$1
    if ! [ $? -eq $code ];then
        echo $msg >&2
        exit 2
    fi
}


# Генерация документации
swag init -g ./cmd/http/main.go

# Build проекта
go build -o main ./cmd/http/main.go
check "go build error"

docker-compose build
check "docker compose build error"

docker-compose up -d
check "docker compose run error"

migrate -path ./schema -database 'postgres://postgres:root@localhost:54320/postgres?sslmode=disable' up
check "migration up error"

if [ $is_run ];then
    # Запуск go
    chmod +x main
    ./main
fi

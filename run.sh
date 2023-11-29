#!/bin/bash

check() {
    local code=$?
    local msg=$1
    if ! [ $? -eq $code ];then
        echo $msg >&2
        exit 2
    fi
}


progname=$0

usage () {
    cat <<HELP_USAGE

    $progname  [--no-build]

    --no-build   Если проект уже скомпилирован, можно запустить без
                 повторной компиляции.
HELP_USAGE
    exit 0
}
#               -f <file>
#     -a           All the instances.
#     -f           File to write all the log lines

#
# Что-то наподобие конфигурации
is_no_build=false
is_help=false
for t in $@;do
    if [ $t == "--no-build" ];then
        is_no_build=true
    elif [ $t == "-r" ];then
        is_run=true
    elif [ $t == "-h" ] || [ $t == "--help" ];then
        is_help=true
    fi
done

# Если передан соответствующий флаг, вывести usage info
$is_help && usage


# Генерация документации
swag init -g ./cmd/http/main.go

# Build проекта
if ! $is_no_build;then
    go build -0 ./cmd/main.go
else
    echo "run without go build"
fi

docker compose build
check "docker compose build error"

docker compose up db -d
sleep 3
docker compose up rdb migrate -d

# check "docker compose run error"

chmod +x main
./main
# if [ $is_run ];then
#     Запуск go
#     chmod +x main
#     ./main
# fi

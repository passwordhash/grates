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
    --docs       Генерация документации API.
                 !!! Необходимо установить swag (go get -u github.com/swaggo/swag/cmd/swag)
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
is_docs=false
for t in $@;do
    if [ $t == "--no-build" ];then
        is_no_build=true
    elif [ $t == "-r" ];then
        is_run=true
    elif [ $t == "-h" ] || [ $t == "--help" ];then
        is_help=true
    elif [ $t == "--docs" ];then
        is_docs=true
    else
        echo "unknown flag $t"
        usage
        exit 1
    fi
done

# Если передан соответствующий флаг, вывести usage info
$is_help && usage

# Генерация документации
if $is_docs;then
    echo "run with swag init"
    swag init -g ./cmd/http/main.go
fi

# Build проекта
if ! $is_no_build;then
    echo "run with go build"
    go build -o main cmd/http/main.go
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

#!/bin/zsh

swag init -g ./cmd/http/main.go

# go run ./cmd/http/main.go
go build -o main ./cmd/http/main.go

chmod +x main

./main


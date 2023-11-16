FROM golang:1.21.4-alpine3.17

RUN mkdir /app

ADD . /app

WORKDIR /app

RUN chmod +x wait-for-postgres.sh

RUN go build -o main ./cmd/main.go

EXPOSE 8080

CMD ["/app/main"]
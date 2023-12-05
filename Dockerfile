FROM golang:1.21.4-alpine3.17

RUN mkdir /app

ADD . /app

WORKDIR /app

# RUN #go build -o main ./cmd/http/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/http/main.go

EXPOSE 8000

CMD ["/app/main"]
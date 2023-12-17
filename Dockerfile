FROM golang:1.21.4-alpine3.17

RUN apk add --no-cache bash

RUN mkdir /app

ADD . /app

WORKDIR /app

# RUN #go build -o main ./cmd/http/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/http/main.go

EXPOSE 8000

COPY wait-for-it.sh wait-for-it.sh

RUN chmod +x wait-for-it.sh

ENTRYPOINT [ "/bin/bash", "-c" ]

#CMD ["./wait-for-it.sh" , "db:5432", "--timeout=300" , "--" , "/app/main"]
#CMD ["sleep", "4"]
CMD ["/app/main"]
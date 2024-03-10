FROM golang:1.21.3

WORKDIR /

COPY . /app


WORKDIR /app

RUN go mod download
RUN go mod tidy
RUN go build cmd/time_tracker/main.go

EXPOSE 8080

CMD ["./main", "--config-path=./config.toml"]
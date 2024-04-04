FROM golang:1.21.3

WORKDIR /

COPY . /app


WORKDIR /app

RUN go mod download
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init --parseDependency --parseInternal -g cmd/time_tracker/main.go
RUN go mod tidy
RUN go build cmd/time_tracker/main.go

EXPOSE 8080

CMD ["./main", "--config-path=./config.toml"]
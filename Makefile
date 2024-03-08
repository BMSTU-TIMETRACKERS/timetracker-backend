# Локальная разработка.
build:
	go build ./cmd/time_tracker/main.go
start: build
	./main --config-path=./config.toml
clean: 
	rm -rf ./main

# Докер


# Генерирует документацию
docs:
	swag init --parseDependency --parseInternal -g cmd/time_tracker/main.go
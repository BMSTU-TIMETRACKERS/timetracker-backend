.PHONY:docs build start clean up
# Локальная разработка.
build:
	go build ./cmd/time_tracker/main.go
start: build
	./main --config-path=./config.toml
clean: 
	rm -rf ./main

# Докер
up:
	docker compose up -d

# Генерирует документацию
docs:
	swag init --parseDependency --parseInternal -g cmd/time_tracker/main.go
.PHONY: setup lint test cover format build clean help

setup: ## Установить зависимости
	go get -u go.uber.org/zap
	go get -u golang.org/x/lint/golint
	go mod tidy

lint: ## Проверить исходный кода на соответствие стандартам
	golint ./...

test: ## Запустить выполнение тестов
	go test -v

cover: ## Запустить отчёт о покрытии кода тестами
	go test -cover

format: ## Отформатировать исходный код
	go fmt ./...

build: ## Собрать версию программы
	GOOS=linux GOARCH=amd64 go build -v ./...

clean: ## Удалить временные файлы
	go clean

help: ## Информация по командам
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := build
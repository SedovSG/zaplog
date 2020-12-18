PROJECT_NAME=$(shell basename "$(PWD)")

.PHONY: setup lint test cover format build clean help

setup: ## Установить зависимости
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.33.0
	go get -u go.uber.org/zap
	go mod tidy

lint: ## Проверить исходный кода на соответствие стандартам
	golangci-lint run ./...

test: ## Запустить выполнение тестов
	go test -v

cover: ## Запустить отчёт о покрытии кода тестами
	go test -cover

format: ## Отформатировать исходный код
	go fmt -n ./...

build: ## Собрать версию программы
	GOOS=linux GOARCH=amd64 go build -o bin/${PROJECT_NAME} -v

clean: ## Удалить временные файлы
	go clean
	rm -rf bin/

help: ## Информация по командам
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := build
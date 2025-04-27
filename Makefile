IMAGE_NAME = sbermortgagecalculator
TAG = $(shell date +%Y%m%d)
TESTING_IMAGE = tparse
LINT_IMAGE = golangci-lint

## Сборка и запуск приложения
.PHONY: build
build: ## Собирает приложение в Docker-контейнер
	docker build -f deployments/Dockerfile.release -t $(IMAGE_NAME):$(TAG) .
	docker tag $(IMAGE_NAME):$(TAG) $(IMAGE_NAME):latest

.PHONY: run
run: ## Запускает приложение в Docker-контейнере
	docker run --rm -p 8080:8080 $(IMAGE_NAME):latest

## Тестирование
.PHONY: test
test: ## Запускает все тесты с покрытием
	docker run --rm -v "$(shell pwd):/app" $(TESTING_IMAGE) sh -c "cd /app; go test -v -cover ./... -json | tparse -all"

.PHONY: test-race
test-race: ## Запускает тесты с проверкой на гонки данных
	docker run --rm -v "$(shell pwd):/app" -e CGO_ENABLED=1 $(TESTING_IMAGE) sh -c "cd /app; go test -race -v ./... -json | tparse -all"

## Линтинг и форматирование
.PHONY: lint
lint: ## Запускает линтер golangci-lint
	@if [ -z "$$(docker images -q $(LINT_IMAGE))" ]; then \
		docker build -f deployments/Dockerfile.lint -t $(LINT_IMAGE) .; \
	fi
	docker run --rm -v "$(shell pwd):/app" $(LINT_IMAGE) sh -c "cd /app; golangci-lint run ./..."

.PHONY: fmt
fmt: ## Форматирует код с помощью gofmt
	docker run --rm -v "$(shell pwd):/app" golang:1.21 sh -c "cd /app; gofmt -w ."

## Очистка
.PHONY: clean
clean: ## Удаляет скомпилированные файлы и Docker-образы
	docker rmi $(IMAGE_NAME):$(TAG) $(IMAGE_NAME):latest $(docker images -f "dangling=true" -q) || true
	find . -name "*.test" -delete

## Справка
.PHONY: help
help: ## Выводит справку по доступным целям
	@echo "Доступные цели:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

dev: clean build
	docker-compose up --build -d

logs:
	docker logs -f calculator

stop_dev:
	docker-compose down

deps:
	docker run --rm -v "$(shell pwd):/app" golang:1.21-alpine sh -c "cd /app; go mod vendor"

image_testing:
	@if [ -z "$$(docker images -q $(TESTING_IMAGE))" ]; then \
		docker build -f deployments/Dockerfile.testing -t $(TESTING_IMAGE) .; \
	fi

image: clean
	docker build -f deployments/Dockerfile.release -t $(IMAGE_NAME):$(TAG) .
	docker tag $(IMAGE_NAME):$(TAG) $(IMAGE_NAME):latest
IMAGE_NAME = sbermortgagecalculator
TAG = $(shell date +%Y%m%d)
TESTING_IMAGE = tparse
LINT_IMAGE = golangci-lint

image_lint:
	@if [ -z "$$(docker images -q $(LINT_IMAGE))" ]; then \
		docker build -f deployments/Dockerfile.lint -t $(LINT_IMAGE) .; \
	fi

lint: image_lint
	docker run --rm -v "$(shell pwd):/app" $(LINT_IMAGE) sh -c "cd /app; golangci-lint run ./..."

image: clean
	docker build -f deployments/Dockerfile.release -t $(IMAGE_NAME):$(TAG) .
	docker tag $(IMAGE_NAME):$(TAG) $(IMAGE_NAME):latest

dev: clean image
	docker-compose up --build -d

logs:
	docker logs -f calculator

stop_dev:
	docker-compose down

run: image
	docker run --rm $(IMAGE_NAME):latest

deps:
	docker run --rm -v "$(shell pwd):/app" golang:1.21-alpine sh -c "cd /app; go mod vendor"

image_testing:
	@if [ -z "$$(docker images -q $(TESTING_IMAGE))" ]; then \
		docker build -f deployments/Dockerfile.testing -t $(TESTING_IMAGE) .; \
	fi

test: image_testing
	docker run --rm -v "$(shell pwd):/app" $(TESTING_IMAGE) sh -c "cd /app; go test -v -cover ./... -json | tparse -all"

clean:
	docker rmi $(IMAGE_NAME):$(TAG) $(IMAGE_NAME):latest $(docker images -f "dangling=true" -q) || true
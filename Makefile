IMAGE_NAME = sbermortgagecalculator
TAG = $(shell date +%Y%m%d)

lint:
	golangci-lint run ./...

image: clean
	docker build -f deployments/Dockerfile.release -t $(IMAGE_NAME):$(TAG) .
	docker tag $(IMAGE_NAME):$(TAG) $(IMAGE_NAME):latest

dev: clean image
	docker-compose up --build -d

stop_dev:
	docker-compose down

run: image
	docker run --rm $(IMAGE_NAME):latest

vendor:
	docker run --rm -v "$(shell pwd):/app" golang:alpine sh -c "cd /app; go mod tidy; go mod vendor"

test: vendor
	docker run --rm -v "$(shell pwd):/app" golang:alpine sh -c "go install github.com/mfridman/tparse@latest; cd /app; go test -v -cover ./... -json | tparse -all"

clean:
	docker rmi $(IMAGE_NAME):$(TAG) $(IMAGE_NAME):latest $(docker images -f "dangling=true" -q) || true
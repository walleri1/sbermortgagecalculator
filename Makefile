all: releases

release:
	docker build -f deployments/Dockerfile.release -t test .

run: clean
	docker-compose up -d

stop:
	docker-compose down

clean:
	docker system prune & docker builder prune & docker image prune -a
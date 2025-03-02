all: releases

release:
	docker build -f deployments/Dockerfile.release -t test .

run:
	docker run -p 8080:8080 test
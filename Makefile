all: releases

release:
	docker build -f deployments/Dockerfile.release -t test .
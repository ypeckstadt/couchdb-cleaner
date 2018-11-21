SHELL := /bin/bash

.PHONY: env-up env-down build-cleaner run-cleaner docker-build-cleaner

##### ENV
env-up: env-down
	@echo "Start environment ..."
	@cd fixtures && docker-compose up -d --force-recreate --build
	@sleep 5
	@cd scripts && ./channel.sh
	@echo "Environment up"

env-down:
	@echo "Stop environment ..."
	@cd fixtures && docker-compose down
	@echo "Environment down"

##### BUILD
build-cleaner:
	@echo "Build cleaner ..."
	@cd cmd/cleaner && go build
	@echo "Build done"


##### RUN
run-cleaner: build-cleaner
	@echo "Start cleaner ..."
	@cd cmd/cleaner/environment && source export-all.sh && cd .. && ./cleaner

##### DOCKER
docker-build-cleaner: build-cleaner
	@echo "Building docker image cleaner ..."
	@cd cmd/cleaner && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ../../fixtures/docker/cleaner
	@cd fixtures/docker && docker build -t couchdb-cleaner -f Dockerfile .
	@cd fixtures/docker && rm cleaner

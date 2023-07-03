
.DEFAULT_GOAL := all

PROJECT_NAME := chatik

.PHONY: all
all:
	@echo "Choose target ..."


.PHONY: lint
lint:
	golangci-lint run ./...


.PHONY: test
test:
	go test -v -race -timeout 5m ./...


.PHONY: ci
ci: lint test


.PHONY: build-local
build-local: ci
	go build -v -o ./build/${PROJECT_NAME} ./cmd/${PROJECT_NAME}


.PHONY: start-local
start-local: build-local
	./build/${PROJECT_NAME} -conf ./configs/local.env


.PHONY: start-infra-local
start-infra-local:
	docker compose -p ${PROJECT_NAME} -f ./deploy/local.docker-compose.yaml up -d


.PHONY: stop-infra-local
stop-infra-local:
	docker compose -p ${PROJECT_NAME} -f ./deploy/local.docker-compose.yaml down


.PHONY: start-web-local
start-web-local:
	cd ./web && yarn dev


.PHONY: clean
clean:
	rm -rf ./build
	rm -rf ./logs



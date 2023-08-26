
.DEFAULT_GOAL := all

PROJECT_NAME = chatik

GO_VERSION := 1.21.0
TOOLS_DIR := tools

LINTER := golangci-lint

.PHONY: all
all: run-stage


#***********************************************************************
# QUALITY CONTROL
#***********************************************************************

.PHONY: lint
lint:
	${LINTER} run ./...


.PHONY: test
test:
	go test -v -race -timeout 5m ./...


.PHONY: ci
ci: lint test


#***********************************************************************
# LOCAL                                   
#***********************************************************************

# TODO: add ci before run-local
.PHONY: run-local
run-local:
	go run ./cmd/chatik -conf ./configs/local/app.yaml


.PHONY: run-infra-local
run-infra-local:
	docker compose -p ${PROJECT_NAME}-local -f ./deploy/local/docker-compose.yaml up -d


.PHONY: stop-infra-local
stop-infra-local:
	docker compose -p ${PROJECT_NAME}-local -f ./deploy/local/docker-compose.yaml down


.PHONY: run-web-local
run-web-local: API_URL="localhost:8080"
run-web-local:
	cd ./web && VITE_API_URL=${API_URL} yarn dev


#***********************************************************************
# STAGE
#***********************************************************************

.PHONY: run-stage
run-stage: API_URL="localhost:8080"
run-stage:
	APP_URL=${API_URL} \
		docker compose -p ${PROJECT_NAME}-stage -f ./deploy/stage/docker-compose.yaml up -d --build


.PHONY: stop-stage
stop-stage:
	docker compose -p ${PROJECT_NAME}-stage -f ./deploy/stage/docker-compose.yaml down


#***********************************************************************
# Tools
#***********************************************************************

.PHONY: mkdir-tools
mkdir-tools:
	mkdir -p ./${TOOLS_DIR}


.PHONY: install-lint
install-lint: LINT_DIR="./${TOOLS_DIR}"
install-lint: mkdir-tools
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ${LINT_DIR} v1.53.3


.PHONY: install-go
install-go: GO_DIR="/usr/local"
install-go:
	curl -O https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz
	rm -rf ${GO_DIR}/go
	tar -C ${GO_DIR} -xzf go${GO_VERSION}.linux-amd64.tar.gz
	rm go1.20.5.linux-amd64.tar.gz
	@echo "Add ${GO_DIR}/go/bin to PATH"

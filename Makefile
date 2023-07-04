
.DEFAULT_GOAL := all

PROJECT_NAME = chatik

.PHONY: all
all:
	@echo "Choose target ..."


#***********************************************************************
# QUALITY CONTROL
#***********************************************************************

.PHONY: lint
lint:
	golangci-lint run ./...


.PHONY: test
test:
	go test -v -race -timeout 5m ./...


.PHONY: ci
ci: lint test



#***********************************************************************
# LOCAL                                   
#***********************************************************************

.PHONY: build-local
build-local: ci
	go build -v -o ./build/${PROJECT_NAME} ./cmd/${PROJECT_NAME}


.PHONY: start-local
start-local: HTTP_ADDR="localhost:8080"
start-local: build-local
	HTTP_ADDR=${HTTP_ADDR} ./build/${PROJECT_NAME} -conf ./configs/local.env


.PHONY: start-infra-local
start-infra-local:
	docker compose -p ${PROJECT_NAME}-local -f ./deploy/local.docker-compose.yaml up -d


.PHONY: stop-infra-local
stop-infra-local:
	docker compose -p ${PROJECT_NAME}-local -f ./deploy/local.docker-compose.yaml down


.PHONY: start-web-local
start-web-local: WEB_API_URL="localhost:8080"
start-web-local: 
	cd ./web && VITE_API_URL=${WEB_API_URL} yarn dev


.PHONY: clean
clean:
	rm -rf ./build
	rm -rf ./logs


#***********************************************************************
# STAGE
#***********************************************************************

# TODO: add to deps: ci
.PHONY: start-stage
start-stage: JWT_SECRET="" MONGO_PASSWORD="" MONGO_USERNAME="" WEB_API_URL=""
start-stage:
	JWT_SECRET=${JWT_SECRET} \
	WEB_API_URL=${WEB_API_URL} \
	MONGO_INITDB_ROOT_PASSWORD=${MONGO_PASSWORD} MONGO_INITDB_ROOT_USERNAME=${MONGO_USERNAME} \
		docker compose -p ${PROJECT_NAME}-stage -f ./deploy/stage.docker-compose.yaml up -d --build


.PHONY: stop-stage
stop-stage:
	docker compose -p ${PROJECT_NAME}-stage -f ./deploy/stage.docker-compose.yaml down


#***********************************************************************
# Tools
#***********************************************************************

.PHONY: mkdir-tools
mkdir-tools:
	mkdir -p ./tools


.PHONY: install-lint
install-lint: LINT_DIR="/usr/local/bin"
install-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ${LINT_DIR} v1.53.3


.PHONY: install-go
install-go: GO_DIR="/usr/local"
install-go:
	curl -O https://dl.google.com/go/go1.20.5.linux-amd64.tar.gz
	rm -rf ${GO_DIR}/go
	tar -C ${GO_DIR} -xzf go1.20.5.linux-amd64.tar.gz
	rm go1.20.5.linux-amd64.tar.gz
	@ehco "Add ${GO_DIR}/go/bin to PATH"

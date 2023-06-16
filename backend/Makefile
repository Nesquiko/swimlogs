.DEFAULT_GOAL := help
DOCKER_IMAGE_VERSION?=latest

.PHONY: audit
audit: ## tidy dependencies and vet
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...

.PHONY: docker-build
docker-build: ## [version=$1] build docker image with the given version
	@echo 'Building docker image with version ${version}'
	cp ../swimlogsApi.yaml ./swimlogsApi.yaml
	docker build --tag docker.io/nesquiko/swimlogs:${DOCKER_IMAGE_VERSION} .
	rm ./swimlogsApi.yaml

.PHONY: help
help: ## show this help message
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m\033[0m\n"} /^[$$()% a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: generate
generate: ## generate Go server code from openapi spec
	@echo "Generating Go server code from openapi spec"
	mkdir -p ./pkg/openapi
	oapi-codegen -config ./openapi/server.cfg.yaml swimlogsApi.yaml
	oapi-codegen -config ./openapi/models.cfg.yaml swimlogsApi.yaml
	sed -i '5s|^|//lint:file-ignore ST1005 generated code\n|' ./pkg/openapi/server.gen.go

.PHONY: local
local: ## run local server
	@echo "Running local server"
	go run . -host localhost -port 42069 \
		-db-host localhost -db-port 5432 \
		-db-user swimlogs -db-pass swimlogs \
		-db-name swimlogs \
		-fe-origin http://localhost:3000

.PHONY: migrations-new
migrations-new: ## [name=$1] create a new database migration with the given name
	@echo 'Creating migration files for ${name}'
	migrate create -seq -ext=.sql -dir=./migrations/ ${name}

.PHONY: test-it
test-it: ## run integration tests against a dockerized postgres database
	@echo "Running integration tests"
	docker run --rm --name ${TEST_DB_DOCKER_NAME} \
		-p ${TEST_DB_PORT}:5432 \
		-e POSTGRES_USER=${TEST_DB_USER} \
		-e POSTGRES_PASSWORD=${TEST_DB_PASSWORD} \
		-e POSTGRES_DB=${TEST_DB_NAME} \
		-d postgres:15.2
	sleep 1 # wait for db to start
	-go test -p 1 -count=1 -v ./...
	docker stop ${TEST_DB_DOCKER_NAME}
test-it: export TEST_DB_HOST = localhost 
test-it: export TEST_DB_PORT = 2345
test-it: export TEST_DB_USER = swimlogs_test
test-it: export TEST_DB_PASSWORD = swimlogs_test
test-it: export TEST_DB_NAME = swimlogs_test
test-it: export TEST_DB_DOCKER_NAME = swimlogs_testdb
test-it: export INTEGRATION = true

.PHONY: test-ut
test-ut: ## run unit tests
	@echo "Running unit tests"
	go test -v ./...


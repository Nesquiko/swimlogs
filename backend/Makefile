TEST_BUILD_ARGS=-tags=integration
.DEFAULT_GOAL := help

.PHONY: help
help: ## show this help message
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m\033[0m\n"} /^[$$()% a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: tidy
tidy: ## format code and tidy modfile
	go fmt ./...
	go mod tidy -v

.PHONY: audit
audit: ## vet and staticcheck code
	go mod verify
	go vet ${TEST_BUILD_ARGS} ./...
	staticcheck ${TEST_BUILD_ARGS} ./...

.PHONY: local
local: ## run local server, needs a local postgres database
	@echo "Running local server"
	air -- -host localhost -port 42069 \
		-debug-level 0 \
		-db-host localhost -db-port 5432 \
		-db-user swimlogs -db-pass swimlogs \
		-db-name swimlogs \
		-fe-origin http://localhost:3000 \
		-tz Europe/Bratislava

.PHONY: test
test: ## run tests against a dockerized postgres database created with testcontainers
	go test ./... -v ${TEST_BUILD_ARGS} -coverprofile=c.out -coverpkg=github.com/Nesquiko/swimlogs/pkg/...

.PHONY: coverage-ui
coverage-ui: ## display coverage HTML report in browser
	go tool cover -html="c.out"

.PHONY: migrations-new
migrations-new: ## [name=$1] create a new database migration with the given name
	@echo 'Creating migration files for ${name}'
	migrate create -seq -ext=.sql -dir=./migrations/ ${name}


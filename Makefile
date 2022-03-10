# A Self-Documenting Makefile: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

.PHONY: up
up:  ## Spin up Temporal cluster
	docker compose up -d  --build --force-recreate

.PHONY: down
down: ## Destroy the Temporal cluster
	docker compose down -v

.PHONY: stop
stop: ## Stop the Temporal cluster
	docker compose stop

.PHONY: ps
ps: ## Check the status of Temporal services
	docker compose ps

.PHONY: shell
shell: ## Start a shell with the Temporal CLI
	docker compose exec temporal-admin-tools bash

.PHONY: db-init
db-init: ## Initialize the database
	go run ./db-init/

.PHONY: api
api: ## Start up API server
	go run ./api/

.PHONY: worker
worker: ## Start the worker
	go run ./worker/

.PHONY: bworker
bworker: ## Start the backfill worker
	go run ./backfill_worker/

.PHONY: fetch
fetch: ## Fetch latest after worker has started
	go run ./starter/

.PHONY: test
test: ## Run tests
	go test ./...

.PHONY: help
.DEFAULT_GOAL := help
help:
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-10s\033[0m %s\n", $$1, $$2}'

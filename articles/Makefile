# Default to dev profile
PROFILE ?= dev

.PHONY: up
up: # Start the service (usage: make up PROFILE=test)
	docker compose --profile $(PROFILE) up -d --build

.PHONY: down
down: # Stop all services across all profiles
	docker compose --profile dev --profile test down

.PHONY: clean
clean:
	docker compose --profile dev --profile test down --volumes

.PHONY: reload
reload:
	docker compose --profile $(PROFILE) up -d --build --no-deps

.PHONY: dev
dev: # Quickly spin up the development environment
	$(MAKE) up PROFILE=dev

.PHONY: test
test: # Quickly spin up the load-test environment
	$(MAKE) up PROFILE=test

.PHONY: logs
logs: # Follow logs for the active profile
	docker compose --profile $(PROFILE) logs -f
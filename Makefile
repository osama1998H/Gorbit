# Makefile
.PHONY: build run run-build remove stop clean test swagger

# Build the Docker containers
build:
	docker-compose build

# Run the Docker containers
run:
	docker-compose up -d

run-build:
	docker-compose up -d --build

remove:
	docker-compose down --remove-orphans

# Stop the Docker containers
stop:
	docker-compose down

# Remove all containers, networks, and volumes
clean:
	docker-compose down -v --rmi all

# Run tests
test:
	go test ./... -v

# Makefile
SWAGGER_VERSION := v1.16.1
SWAGGER_DIR := docs/swagger
SWAGGER_MAIN := cmd/api/main.go
GOBIN := $(shell go env GOPATH)/bin


swagger: ## Generate Swagger documentation
	@echo "Generating Swagger documentation..."
	@if ! command -v $(GOBIN)/swag >/dev/null; then \
		echo "Installing swag@$(SWAGGER_VERSION)..."; \
		go install github.com/swaggo/swag/cmd/swag@$(SWAGGER_VERSION); \
	fi
	@mkdir -p $(SWAGGER_DIR)
	@$(GOBIN)/swag init \
		--generalInfo $(SWAGGER_MAIN) \
		--output $(SWAGGER_DIR) \
		--outputTypes json,yaml \
		--parseDependency \
		--dir ./internal/api/v1
	@echo "Swagger docs generated at $(SWAGGER_DIR)"
	@echo "Version: $(SWAGGER_VERSION)"
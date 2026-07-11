# Main Makefile for the project

# Include example Makefile

define PRINT_INCLUDED_HELP
	@echo "$(1) targets:"
	@awk -F ':|##' '/^[^\t].+?:.*?##/ { printf "  %-20s %s\n", $$1, $$NF }' $(2)
endef

help: ## Display help message
	@echo "Usage:"
	@echo "  make <target>"
	@echo ""
	@echo "Makefile targets:"
	@awk '/^[a-zA-Z_-]+:.*?## .*$$/ { \
		printf "  %-21s %s\n", substr($$1, 1, index($$1, ":")-1), substr($$0, index($$0, "##")+3) \
	}' $(firstword $(MAKEFILE_LIST))
	@echo ""

.PHONY: help

sqlc: ## Generate SQLC queries and models
	@sqlc generate

seed: ## Seed the database with initial data
	@goose -dir ./database/seed -no-versioning up

resetseed: ## Reset the database seed
	@goose -dir ./database/seed -no-versioning reset

templ: ## Run Templ proxy for live reloading of templates
	@templ generate -watch -proxy=http://127.0.0.1:8080 -proxyport=8081

templ-build: ## Generate Templ templates
	@templ generate

tailwind: ## Generate Tailwind static CSS file
	@tailwindcss -i ./styles/tailwind.css -o ./static/styles.css

init: init-deps sqlc templ-build tailwind ## Initialize the project (fetch libraries, generate SQLC queries and models)
	@test -f .env || cp example.env .env


build: init sqlc templ-build tailwind ## Build the project
	@go build -o stock main.go

run: build ## Run the built binary
	@./stock

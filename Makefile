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

init-deps: ## Fetch all assets from assets.json
	@jq -r 'to_entries[] | "\(.value.destination)|\(.value.url)|\(.value.sha256)"' assets.json | \
	while IFS='|' read -r dest url sha256; do \
		if [ -f "$$dest" ]; then \
			echo "✓ $$dest already exists"; \
		else \
			echo "Fetching $$url to $$dest..."; \
			mkdir -p $$(dirname "$$dest"); \
			tmp="$$dest.tmp"; \
			curl -fsSL "$$url" -o "$$tmp" || (rm -f "$$tmp"; exit 1); \
			echo "$$sha256  $$tmp" | sha256sum -c - || (rm -f "$$tmp"; exit 1); \
			mv "$$tmp" "$$dest"; \
		fi; \
	done

verify-assets: ## Verify all assets match their checksums from assets.json
	@jq -r 'to_entries[] | "\(.value.destination)|\(.value.sha256)"' assets.json | \
	while IFS='|' read -r dest sha256; do \
		if [ ! -f "$$dest" ]; then \
			echo "✗ $$dest missing"; exit 1; \
		fi; \
		echo "$$sha256  $$dest" | sha256sum -c -; \
	done

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

test: init ## Run tests
	@go test ./...

build: init sqlc templ-build tailwind ## Build the project
	@go build -o stock main.go

run: build ## Run the built binary
	@./stock

docker-image:
	@nix build .#dockerImage
	@./result | docker load
	@rm -f ./result
	@image_name=$${IMAGE_NAME:-stock}; \
    image_tag=$${IMAGE_TAG:-localdev}; \
	image_tag_list=$${IMAGE_TAG_LIST:-}; \
	if [ -n "$$image_tag_list" ]; then \
		for tag in $$(echo "$$image_tag_list" | tr ',' ' '); do \
			echo "Tagging docker image as $$image_name:$$tag"; \
			docker tag stock:localdev "$$image_name:$$tag"; \
		done; \
	elif [ -n "$${IMAGE_NAME}$${IMAGE_TAG}" ]; then \
		echo "Tagging docker image as $$image_name:$$image_tag"; \
        docker tag stock:localdev "$$image_name:$$image_tag"; \
    fi

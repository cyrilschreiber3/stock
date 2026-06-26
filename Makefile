init-daisyui:
	@test -f ./styles/daisyui.mjs || curl -o ./styles/daisyui.mjs -fsSL https://github.com/saadeghi/daisyui/releases/download/v5.5.19/daisyui.mjs
	@test -f ./styles/daisyui-theme.mjs || curl -o ./styles/daisyui-theme.mjs -fsSL https://github.com/saadeghi/daisyui/releases/download/v5.5.19/daisyui-theme.mjs

init-htmx:
	@test -f ./static/htmx.min.js || curl -o static/htmx.min.js -fsSL https://unpkg.com/htmx.org@2.0.4

init-alpine:
	@test -f ./static/alpine.min.js || curl -o static/alpine.min.js -fsSL https://unpkg.com/alpinejs@3.15.11/dist/cdn.min.js

sqlc:
	@sqlc generate

init: init-daisyui init-htmx init-alpine sqlc
	@test -f .env || cp example.env .env

seed:
	@goose -dir ./database/seed -no-versioning up

resetseed:
	@goose -dir ./database/seed -no-versioning reset

templ:
	@templ generate -watch -proxy=http://127.0.0.1:8080 -proxyport=8081

templ-build:
	@templ generate

tailwind:
	@tailwindcss -i ./styles/tailwind.css -o ./static/styles.css

build: init-daisyui sqlc templ-build tailwind
	@go build -o stock main.go

run: build
	@./stock

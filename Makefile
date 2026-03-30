init-daisyui:
	@test -f ./styles/daisyui.css || curl -o ./styles/daisyui.css https://github.com/saadeghi/daisyui/releases/v5.5.19/download/daisyui.mjs
	@test -f ./styles/daisyui-theme.mjs || curl -o ./styles/daisyui-theme.mjs https://github.com/saadeghi/daisyui/releases/v5.5.19/download/daisyui-theme.mjs

templ:
	@templ generate -watch -proxy=http://127.0.0.1:8080 -proxyport=8081

templ-build:
	@templ generate

tailwind:
	@tailwindcss -i ./styles/tailwind.css -o ./static/styles.css

build: init-daisyui templ-build tailwind
	@go build -o spotify-playlist-manager main.go

run: build
	@./spotify-playlist-manager

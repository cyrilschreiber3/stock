run: build
	@./spotify-playlist-manager

templ:
	@templ generate -watch -proxy=http://127.0.0.1:8080 -proxyport=8081

tailwind:
	@tailwindcss -i ./tailwind.css -o ./static/styles.css

build: tailwind
	@templ generate
	@go build -o spotify-playlist-manager main.go
build:
	@go build -o bin/atom-fit cmd/app/main.go

run: build
	@./bin/atom-fit -env-path=.env
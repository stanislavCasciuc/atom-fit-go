build:
	@go build -o bin/atom-fit cmd/app/main.go

run: build
	@./bin/atom-fit -env-path=.env

migration:
	@migrate create -ext sql -dir migrations $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@go run cmd/migrator/main.go -migrations-path=migrations -env-path=.env up

migrate-down:
	@go run cmd/migrator/main.go -migrations-path=migrations -env-path=.env down

force-version:
	@go run cmd/migrator/main.go -migrations-path=migrations -env-path=.env -force=$(version)
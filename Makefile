build:
	@go build -o ./bin/blog-nest ./cmd/main.go

run: build
	@./bin/blog-nest

test:
	echo $(MAKECMDGOALS)

migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@, $(MAKECMDGOALS))

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down
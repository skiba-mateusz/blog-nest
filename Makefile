build:
	@go build -o ./bin/blog-nest ./main.go

run: build
	@./bin/blog-nest
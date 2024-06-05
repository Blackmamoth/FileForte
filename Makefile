build:
	@go build -o bin/fileforte cmd/main.go

run: build
	@./bin/fileforte
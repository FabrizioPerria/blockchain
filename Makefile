build:
	@go build -o bin/blockchain

run: build
	@./bin/blockchain

test:
	@go test -v ./...

cover:
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out

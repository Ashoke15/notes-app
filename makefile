build:
	@go build -0 notes-app .

run:
	@go run cmd/totion/main.go

test:
	@go test -v ./...
build:
	CGO_ENABLED=0 go build -o ./shellfire

test:
	go test ./...

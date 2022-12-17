build: $(find . -name "*.go")
	CGO_ENABLED=0 go build -o ./shellfire

install: build
	sudo mv -f ./shellfire /usr/local/bin

test:
	go test ./...

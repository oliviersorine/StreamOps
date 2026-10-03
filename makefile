run:
	go run ./cmd/streamops start

format:
	go fmt ./...

test:
	go test ./...

lint:
	go vet ./...

build:
	go build -o bin/streamops ./cmd/streamops

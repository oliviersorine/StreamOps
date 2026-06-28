run:
	go run ./cmd/streamops start

test:
	go test ./...

lint:
	go vet ./...

build:
	go build -o bin/streamops ./cmd/streamops

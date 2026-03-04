.PHONY: build test lint install clean

build:
	go build -o bin/vn ./cmd/vn

test:
	go test ./...

lint:
	golangci-lint run

install:
	go install ./cmd/vn

clean:
	rm -rf bin/

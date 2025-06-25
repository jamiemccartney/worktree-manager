BINARY_NAME := wt

build:
	go build -o $(BINARY_NAME)

run:
	go run main.go

test:
	go test ./...

clean:
	rm -f wt

deps:
	go mod tidy
	go mod download

install: build
	cp $(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)

.PHONY: build run test clean deps
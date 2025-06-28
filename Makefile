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
	go mod download

tidy:
	go mod tidy

lint:
	staticcheck ./...

install-linter:
	brew install staticcheck

install: build
	cp $(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)

.PHONY: build run test install-linter lint clean deps tidy
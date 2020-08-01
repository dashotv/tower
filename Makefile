all: test

test:
	go test -v ./...

generate:
	golem generate

build: generate
	go build

server:
	go run main.go server

.PHONY: server receiver test generate

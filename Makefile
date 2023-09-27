include .env
export $(shell sed 's/=.*//' .env)

all: test

test:
	go test -v ./...

generate:
	golem generate

build: generate
	go build

install: build
	go install

server: generate
	go run main.go server

deps:
	go install github.com/dashotv/golem@latest

docker:
	docker build -t tower .

docker-run:
	docker run -d --rm --name tower -p 9000:9000 tower

# this works on linux, for mac you use host.docker.internal
# DOCKER_HOST := `/sbin/ifconfig docker0 | grep 'inet addr:' | cut -d: -f2 | awk '{ print $1}'`
# to allow docker to talk to localhost mongo
# $ docker run -it my_application --add-host 'DOCKER_HOST:$DOCKER_HOST'

check-env:
ifndef TEST_MONGODB_URL
	$(error TEST_MONGODB_URL is undefined)
endif

.PHONY: test generate build install server deps docker docker-run check-env

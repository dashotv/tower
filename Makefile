all: test

test:
	go test -v ./...

generate:
	golem generate

build: generate
	go build

install: build
	go install

server:
	go run main.go server

deps:
	go get golang.org/x/tools/cmd/goimports
	go get github.com/dashotv/golem

production:
	sed -i \
		-e 's/seer_development/seer_production/g' \
		-e 's/torch_development/torch_production/g' \
		.golem/.golem.yaml

docker:
	docker build -t tower .

docker-run:
	docker run -d --rm --name tower -p 9000:9000 tower

# this works on linux, for mac you use host.docker.internal
# DOCKER_HOST := `/sbin/ifconfig docker0 | grep 'inet addr:' | cut -d: -f2 | awk '{ print $1}'`
# to allow docker to talk to localhost mongo
# $ docker run -it my_application --add-host 'DOCKER_HOST:$DOCKER_HOST'

.PHONY: server receiver test generate deps

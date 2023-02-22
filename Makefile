.PHONY: build clean fmt test proto run examples test

DOCKER_REPO ?= docker.io/vesoft
IMAGE_TAG ?= latest


default: build

clean:
	go mod tidy 

fmt:
	find . -type f -iname \*.go -exec go fmt {} \;

build: clean fmt
	go build -trimpath -ldflags "-X main.GitInfoSHA=`git rev-parse --short HEAD`" -o ./bin/agent cmd/agent.go
	go build -o ./bin/client examples/storage.go
	chmod +x ./bin/agent
	chmod +x ./bin/client

docker-build:
	docker build -t "${DOCKER_REPO}/nebula-agent:${IMAGE_TAG}" .

docker-push:
	docker push "${DOCKER_REPO}/nebula-agent:${IMAGE_TAG}"

proto:
	sh scripts/gen_proto.sh

run: build
	 ./bin/agent

test:
	go test -v ./...

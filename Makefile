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

build-analytics:build
	mkdir -p ./plugins/analytics
	cp ./packages/analytics/config.yaml ./plugins/analytics
	go build -trimpath -buildmode=plugin -ldflags "-X main.GitInfoSHA=`git rev-parse --short HEAD`" -o ./plugins/analytics/analytics.so packages/analytics/main.go

docker-build:
	docker build -t "${DOCKER_REPO}/nebula-agent:${IMAGE_TAG}" .

docker-push:
	docker push "${DOCKER_REPO}/nebula-agent:${IMAGE_TAG}"

proto:
	sh scripts/gen_proto.sh

run: build
	 ./bin/agent -f ./etc/config.test.yaml

test:
	go test -v ./...

test-shell:
	go build -o  ./bin/shell ./examples/task/shell.go
	bin/shell 127.0.0.1:8888 "ls"

test-streamShell:
	go build -o  ./bin/shell ./examples/streamshell/streamshell.go
	bin/shell 127.0.0.1:8888 "top"
 
run-analytics: build-analytics
	 ./bin/agent

test-analytics:
	go test -timeout 30s github.com/vesoft-inc/nebula-agent/v3/packages/analytics/pkg/service/task -v
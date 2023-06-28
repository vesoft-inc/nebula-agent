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

build-analytics:
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

testShell:
	go build -o  ./bin/shell ./examples/task/shell.go
	bin/shell 127.0.0.1:8888 "ls"

testStreamShell:
	go build -o  ./bin/shell ./examples/streamshell/streamshell.go
	bin/shell 127.0.0.1:8888 "top"

testHeartbeat:
	go build -o  ./bin/heartbeat-server ./examples/server/server.go
	bin/heartbeat-server

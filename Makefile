.PHONY: build clean fmt test proto run examples test

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

proto:
	sh scripts/gen_proto.sh

run: build
	 ./bin/agent

test:
	go test -v ./...

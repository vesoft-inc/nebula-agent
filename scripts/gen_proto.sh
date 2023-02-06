#!/usr/bin/env bash

function push() {
    pushd $1 >/dev/null 2>&1
}

function pop() {
    popd >/dev/null 2>&1
}

SCRIPTS_DIR=$(dirname "$0")

push $SCRIPTS_DIR/..
NEBULA_AGENT_ROOT=`pwd`
pop

PROGRAM=$(basename "$0")
GOPATH=$(go env GOPATH)

if [ -z $GOPATH ]; then
    printf "Error: the environment variable GOPATH is not set, please set it before running %s\n" $PROGRAM > /dev/stderr
    exit 1
fi

echo "install tools..."
GO111MODULE=off go get -u github.com/golang/protobuf/proto
GO111MODULE=off go get -u github.com/gogo/protobuf/protoc-gen-gofast
GO111MODULE=off go get google.golang.org/grpc


echo "collect proto..."
GO_PREFIX_PATH=github.com/vesoft-inc/nebula-agent/v3/pkg/proto
export PATH=$NEBULA_AGENT_ROOT/_tools/bin:$GOPATH/bin:$PATH

function collect() {
    file_name=$(basename $1)
    base_name=$(basename $file_name ".proto")
    if [ -z $GO_OUT_M ]; then
        GO_OUT_M="M$file_name=$GO_PREFIX_PATH"
    else
        GO_OUT_M="$GO_OUT_M,M$file_name=$GO_PREFIX_PATH"
    fi
}

mkdir -p pkg/proto
cd proto
for proto_file in `ls *.proto`
do
    collect $proto_file
done

echo "generate go code..."
ret=0

function gen() {
    base_name=$(basename $1 ".proto")
    protoc -I. --gofast_out=plugins=grpc,$GO_OUT_M:../pkg/proto $1 || ret=$?
    cd ../pkg/proto
    goimports -w *.pb*.go
    cd ../../proto
}

for proto_file in `ls *.proto`
do
    gen $proto_file
done
exit $ret

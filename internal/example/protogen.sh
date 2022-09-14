#!/bin/sh

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

cd "$SCRIPT_DIR"

protoc -I . \
    --go_out=. \
    --go_opt=paths=source_relative \
    --go-grpc_out=. \
    --go-grpc_opt=paths=source_relative \
    --grpc-transcode_out=. \
    --grpc-transcode_opt=paths=source_relative \
    --plugin=protoc-gen-grpc-transcode=./protoc-gen-grpc-transcode \
    *.proto

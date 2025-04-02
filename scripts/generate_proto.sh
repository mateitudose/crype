#!/bin/bash
protoc --proto_path=../api/proto --go_out=../api/generated \
       --go_opt=paths=source_relative --go-grpc_out=../api/generated \
       --go-grpc_opt=paths=source_relative "$(find ../api/proto -name '*.proto')"
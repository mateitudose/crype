#!/bin/bash
protoc --proto_path=../api/proto --go_out=../api/gen_proto \
       --go_opt=paths=source_relative --go-grpc_out=../api/gen_proto \
       --go-grpc_opt=paths=source_relative "$(find ../api/proto -name '*.proto')"
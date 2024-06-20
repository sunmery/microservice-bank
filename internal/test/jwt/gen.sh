#!/usr/bin/env bash
# 启用 POSIX 模式并设置严格的错误处理机制
set -o posix errexit -o pipefail

protoc --proto_path=. \
       --proto_path=./third_party \
       --go_out=. --go_opt=paths=source_relative \
       --go-http_out=. --go-http_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       ./internal/test/jwt/helloworld.proto

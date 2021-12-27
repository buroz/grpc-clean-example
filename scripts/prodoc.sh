#!/bin/bash

# for dict in ./internal/*; do
#     if [ -d "${dict}" ]; then
#         file="${dict//.\/internal\//}"
#         path="${dict}/${file}.proto"
#         protoc --go_out=plugins=grpc:. --go_opt=paths=source_relative ${path}
#     fi
# done

for dict in ./internal/*; do
    if [ -d "${dict}" ]; then
        protoc --go_out=plugins=grpc:. --go_opt=paths=source_relative ${dict}/*.proto
    fi
done

#!/usr/bin/env bash
RUN_NAME="sheepim-push-worker-service"

mkdir -p output/bin
cp script/* output/
chmod +x output/bootstrap.sh
go mod tidy
if [ "$IS_SYSTEM_TEST_ENV" != "1" ]; then
    go build -o output/bin/${RUN_NAME}
else
    go test -c -covermode=set -o output/bin/${RUN_NAME} -coverpkg=./...
fi

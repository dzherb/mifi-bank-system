#!/usr/bin/env bash
set -e

PROJECT_ROOT="$(dirname "$(realpath "$0")")/../"

go test "$PROJECT_ROOT/..." -coverprofile=./cover.out -covermode=atomic -coverpkg=./...
go-ignore-cov --file cover.out

if [[ "$1" == "html" ]]; then
    go tool cover -html=cover.out
fi

echo
go tool cover -func=cover.out | awk 'END { print "total: " $3 }'
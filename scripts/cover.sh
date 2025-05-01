#!/usr/bin/env bash
set -e

PROJECT_ROOT="$(dirname "$(realpath "$0")")/../"

go test "$PROJECT_ROOT/..." -coverprofile=./cover.out -covermode=atomic -coverpkg=./...
go tool cover -html=cover.out
echo
go tool cover -func=cover.out | tail -1 | awk '{ print "total: " $3 }'
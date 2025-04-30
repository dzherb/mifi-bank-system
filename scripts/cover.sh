#!/usr/bin/env bash
set -e

PROJECT_ROOT="$(dirname "$(realpath "$0")")/../"

go test "$PROJECT_ROOT/..." -coverprofile=./cover.out -covermode=atomic -coverpkg=./...
go tool cover -html=cover.out

GOBIN ?= $$(go env GOPATH)/bin

.PHONY: test
test:
	@go test ./... -count=1

.PHONY: cover
cover:
	@go test ./... -coverprofile=./cover.out -covermode=atomic -coverpkg=./...

.PHONY: install-go-test-coverage
install-go-test-coverage:
	@go install github.com/vladopajic/go-test-coverage/v2@latest

.PHONY: check-coverage
check-coverage: install-go-test-coverage cover
	@${GOBIN}/go-test-coverage --config=./.testcoverage.yml

.PHONY: clean
clean:
	@rm -f cover.out

.PHONY: schema
schema:
	@./scripts/schema.sh

.PHONY: lint
lint:
	@golangci-lint run ./...

.PHONY: fmt
fmt:
	@golangci-lint fmt ./...

.PHONY: check
check: lint cover check-coverage

.PHONY: migrations
migrations:
	@./scripts/migrations.sh
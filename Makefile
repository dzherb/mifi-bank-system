.PHONY: test cover clean schema lint fmt

test:
	@go test ./... -count=1

cover:
	@./scripts/cover.sh

clean:
	@rm -f cover.out

schema:
	@./scripts/schema.sh

lint:
	@golangci-lint run ./...

fmt:
	@golangci-lint fmt ./...
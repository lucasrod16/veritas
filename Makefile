MODULE := github.com/lucasrod16/veritas
PKGS=$(shell go list ./... | grep -v $(MODULE)/cmd/schema-validator | grep -v $(MODULE)/cmd/veritas | tr '\n' ' ')

build:
	CGO_ENABLED=0 go build -o ./bin/veritas ./cmd/veritas

build-schema-validator:
	CGO_ENABLED=0 go build -o ./bin/schema-validator ./cmd/schema-validator

test:
	@go test -v -failfast -coverprofile=c.out $(PKGS)

cover:
	go tool cover -func=c.out

validate-schema: build build-schema-validator
	./scripts/validate-schema.sh

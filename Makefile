default: build

build:
	CGO_ENABLED=0 go build -o ./bin/veritas .

build-schema-validator:
	CGO_ENABLED=0 go build -o ./bin/schema-validator ./cmd/schema-validator

test:
	@go test -v -failfast -coverprofile=c.out ./pkg/...

cover:
	go tool cover -func=c.out

validate-schema: build build-schema-validator
	./scripts/validate-schema.sh

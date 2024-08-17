build:
	CGO_ENABLED=0 go build -o veritas

test-unit:
	go test -v -failfast -coverprofile=c.out ./...

cover:
	go tool cover -func=c.out

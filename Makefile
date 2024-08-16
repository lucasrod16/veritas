build:
	CGO_ENABLED=0 go build -o veritas

test-unit:
	go test -coverprofile=c.out -failfast ./...

cover:
	go tool cover -func=c.out

GOLANGCI_LINT_VERSION = 1.58.0
ACTIONLINT_VERSION = 1.6.27

run:
	go run cmd/core-service/main.go -dev-logger

test:
	go test --timeout 10m -race ./...

coverage:
	go test -race -v -coverpkg=./... -coverprofile=profile.out ./...
	go tool cover -func profile.out

lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@v$(GOLANGCI_LINT_VERSION) run

actionlint:
	go run github.com/rhysd/actionlint/cmd/actionlint@v$(ACTIONLINT_VERSION)

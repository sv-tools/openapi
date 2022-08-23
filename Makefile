all: install tidy lint fmt test

install:
	@go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest

test:
	@go test -cover -race ./...

fmt:
	@go fmt ./...

lint:
	@go vet ./...
	@go vet -vettool=$(which fieldalignment) ./...

tidy:
	@go mod tidy

godoc:
	@go install golang.org/x/tools/cmd/godoc@latest
	@echo http://localhost:6060/pkg/github.com/sv-tools/openapi/
	@godoc -http=:6060 >/dev/null

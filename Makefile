BREWFILE=./.github/Brewfile

ifeq ($(shell uname), Darwin)
all: brew-install
endif

all: go-install tidy lint test

brew-install:
	@brew bundle --file $(BREWFILE)

go-install:
	@go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest

run-test:
	@go test -cover -race ./...

test: run-test

# @golangci-lint run --fix ./...
lint:
	@go vet ./...
	@go vet -vettool=$(which fieldalignment) ./...

tidy:
	@go mod tidy

godoc:
	@go install golang.org/x/tools/cmd/godoc@latest
	@echo http://localhost:6060/pkg/github.com/sv-tools/openapi/
	@godoc -http=:6060 >/dev/null

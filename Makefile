BREWFILE=./.github/Brewfile

ifeq ($(shell uname), Darwin)
all: brew-install
endif

all: tidy lint test

brew-install:
	@brew bundle --file $(BREWFILE)

run-test:
	@go test -cover -race ./...

test: run-test

# @golangci-lint run --fix ./...
lint:
	@go vet ./...

tidy:
	@go mod tidy

godoc:
	@go install golang.org/x/tools/cmd/godoc@latest
	@echo http://localhost:6060/pkg/github.com/sv-tools/openapi/
	@godoc -http=:6060 >/dev/null

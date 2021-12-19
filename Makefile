BREWFILE=./.github/Brewfile

ifeq ($(shell uname), Darwin)
all: brew-install
endif

# +lint
all: tidy test

brew-install:
	@brew bundle --file $(BREWFILE)

run-test:
	@go test -cover -race ./...

test: run-test

lint:
	@golangci-lint run --fix ./...

tidy:
	@go mod tidy

godoc:
	@go install golang.org/x/tools/cmd/godoc@latest
	@echo http://localhost:6060/pkg/github.com/sv-tools/openapi/
	@godoc -http=:6060 >/dev/null

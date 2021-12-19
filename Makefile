BREWFILE=./.github/Brewfile

NO_COLOR=\033[0m
OK_COLOR=\033[32;01m

ifeq ($(shell uname), Darwin)
all: brew-install
endif

all: tidy lint test done

done:
	@echo "$(OK_COLOR)==> Done.$(NO_COLOR)"

brew-install:
	@echo "$(OK_COLOR)==> Checking and installing dependencies using brew...$(NO_COLOR)"
	@brew bundle --file $(BREWFILE)

run-test:
	@echo "$(OK_COLOR)==> Testing...$(NO_COLOR)"
	@richgo test -cover -race

run-benchmark:
	@echo "$(OK_COLOR)==> Benchmarks...$(NO_COLOR)"
	@richgo test -benchmem -run=Bench -bench=. .

test: run-test run-benchmark

lint:
	@echo "$(OK_COLOR)==> Linting via golangci-lint...$(NO_COLOR)"
	@golangci-lint run --fix ./...

tidy:
	@echo "$(OK_COLOR)==> Updating go.mod...$(NO_COLOR)"
	@go mod tidy

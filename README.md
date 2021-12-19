# go-repo-template

[![Code Analysis](https://github.com/sv-tools/go-repo-template/actions/workflows/checks.yaml/badge.svg)](https://github.com/sv-tools/go-repo-template/actions/workflows/checks.yaml)
[![Go Reference](https://pkg.go.dev/badge/github.com/sv-tools/go-repo-template.svg)](https://pkg.go.dev/github.com/sv-tools/go-repo-template)
[![codecov](https://codecov.io/gh/sv-tools/go-repo-template/branch/main/graph/badge.svg?token=0XVOTDR1CW)](https://codecov.io/gh/sv-tools/go-repo-template)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/sv-tools/go-repo-template?style=flat)](https://github.com/sv-tools/go-repo-template/releases)

The template for new go repositories

## Features

1. `Makefile` to run some basic validations:
   1. `golangci-lint`
   2. `nancy`
   3. unit and benchmarks tests
   4. installing all needed tools on macOS using `brew`
   5. cleaning up the `go.sum` file by removing it and re-creating by `go mod tidy`
2. MIT License by default
3. GitHub Action workflows:
   1. testing all pull requests by running same tools and checking code coverage using `codecov` action
   2. making a new release, triggered by closed milestone
      1. creates a new tag using `bumptag` tool
      2. creates new `Next` milestone
      3. runs `goreleaser` to build a new release

## Usage

1. Create a repository using this repo as template
2. Replace in all files `go-repo-template` to the project's name
3. In case of:
   1. library
       1. Remove `.github/Dockerfile`, `.github/goreleaser-cli.yml` files
       2. Remove `release-cli` section in the `.github/workflows/release.yaml` file
   2. command line tool (cli)
      1. Remove `.github/goreleaser-lib.yml` file
      2. Remove `release-lib` section in the `.github/workflows/release.yaml` file
4. Modify `README.md` by removing this text
5. Feel free to modify any other files


## License

MIT licensed. See the bundled [LICENSE](LICENSE) file for more details.

# OpenAPI v3.1 Specification

[![Code Analysis](https://github.com/sv-tools/openapi/actions/workflows/checks.yaml/badge.svg)](https://github.com/sv-tools/openapi/actions/workflows/checks.yaml)
[![Go Reference](https://pkg.go.dev/badge/github.com/sv-tools/openapi.svg)](https://pkg.go.dev/github.com/sv-tools/openapi)
[![codecov](https://codecov.io/gh/sv-tools/openapi/branch/main/graph/badge.svg?token=0XVOTDR1CW)](https://codecov.io/gh/sv-tools/openapi)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/sv-tools/openapi?style=flat)](https://github.com/sv-tools/openapi/releases)

The implementation of OpenAPI v3.1 Specification for Go using generics.

## Supported Go versions:

* v1.19
* v1.18

## Features

* The `spec` folder contains full implementation of the v3.1 Specification using generics, so the minimum supported version of Go is `v1.18`.
* The `validate` folder uses [jsonschema/v5](https://github.com/santhosh-tekuri/jsonschema) with draft 2020-12 to validate the specifications.
* The official v3.0 and v3.1 [examples](https://github.com/OAI/OpenAPI-Specification/tree/main/examples) are tested.
  In most cases v3.0 specification can be converted to v3.1 by changing the version's parameter only.
  ```diff
  @@ -1,4 +1,4 @@
  -openapi: "3.0.0"
  +openapi: "3.1.0"
  ```

**NOTE**: The descriptions of most structures and their fields are taken from the official documentations.

## Links

* OpenAPI Specification: <https://github.com/OAI/OpenAPI-Specification> and <https://spec.openapis.org/oas/v3.1.0>
* JSON Schema: <https://json-schema.org/understanding-json-schema/index.html> and <https://json-schema.org/draft/2020-12/json-schema-core.html>
* The list of most popular alternatives: <https://github.com/OAI/OpenAPI-Specification/blob/main/IMPLEMENTATIONS.md#low-level-tooling>

## License

MIT licensed. See the bundled [LICENSE](LICENSE) file for more details.

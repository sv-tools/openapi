# OpenAPI v3.1 Specification

[![Code Analysis](https://github.com/sv-tools/openapi/actions/workflows/code.yaml/badge.svg)](https://github.com/sv-tools/openapi/actions/workflows/code.yaml)
[![Go Reference](https://pkg.go.dev/badge/github.com/sv-tools/openapi.svg)](https://pkg.go.dev/github.com/sv-tools/openapi)
[![codecov](https://codecov.io/gh/sv-tools/openapi/branch/main/graph/badge.svg?token=0XVOTDR1CW)](https://codecov.io/gh/sv-tools/openapi)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/sv-tools/openapi?style=flat)](https://github.com/sv-tools/openapi/releases)

The implementation of OpenAPI v3.1 Specification for Go using generics.

```shell
go get github.com/sv-tools/openapi
```

## Supported Go versions

* v1.25
* v1.24
* v1.23
* v1.22

## Versions

* v0 - **Deprecated**. The initial version with the full implementation of the v3.1 Specification using generics. See the `v0` branch.
* v1 - The current version with the in-place validation of the specification.
  * The minimum version of Go is `v1.22`.
  * Everything has been moved to the root folder. So, the import path is `github.com/sv-tools/openapi`.
  * Added `Validator` struct for validation of the specification and the data.
    * `Validator.ValidateSpec()` method validates the specification.
    * `Validator.ValidateData()` method validates the data.
    * `Validator.ValidateDataAsJSON()` method validates the data by converting it into `map[string]any` type first using `json.Marshal` and `json.Unmarshal`.
      **WARNING**: the function is slow due to double conversion.
  * Added `ParseObject` function to create `SchemaBuilder` by parsing an object.
    The function supports `json`, `yaml` and `openapi` field tags for the structs.
  * Use OpenAPI `v3.1.1` by default.

## Features

* The official v3.0 and v3.1 [examples](https://github.com/OAI/OpenAPI-Specification/tree/main/examples) are tested.
  In most cases, the v3.0 specification can be converted to v3.1 by changing only the version parameter.

  ```diff
  @@ -1,4 +1,4 @@
  -openapi: "3.0.0"
  +openapi: "3.1.0"
  ```

**NOTE**: The descriptions of most structures and their fields are taken from the official documentations.

## Links

* OpenAPI Specification: <https://github.com/OAI/OpenAPI-Specification> and <https://spec.openapis.org/oas/v3.1.0>
* JSON Schema: <https://json-schema.org/understanding-json-schema/index.html> and <https://json-schema.org/draft/2020-12/json-schema-core.html>
* The list of most popular alternatives: <https://tools.openapis.org>

## License

MIT licensed. See the bundled [LICENSE](LICENSE) file for more details.

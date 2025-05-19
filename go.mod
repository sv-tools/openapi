module github.com/sv-tools/openapi

go 1.22.0

retract v0.3.0 // due to a mistake, there is no real v0.3.0 release, it was pointed to v0.2.2 tag

require (
	github.com/santhosh-tekuri/jsonschema/v6 v6.0.1
	github.com/stretchr/testify v1.10.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/goccy/go-yaml v1.17.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)

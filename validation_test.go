package openapi_test

import (
	"encoding/json"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/sv-tools/openapi"
)

func TestValidator_ValidateSpec(t *testing.T) {
	info, err := os.ReadDir("testdata")
	require.NoError(t, err)

	for _, f := range info {
		if f.IsDir() {
			continue
		}
		name := f.Name()
		t.Run(name, func(t *testing.T) {
			data, err := os.ReadFile(path.Join("testdata", name))
			require.NoError(t, err)
			var o *openapi.Extendable[openapi.OpenAPI]
			switch path.Ext(name) {
			case ".yaml":
				require.NoError(t, yaml.Unmarshal(data, &o))
				newData, err := yaml.Marshal(&o)
				require.NoError(t, err)
				require.YAMLEq(t, string(data), string(newData))
			case ".json":
				require.NoError(t, json.Unmarshal(data, &o))
				newData, err := json.Marshal(&o)
				require.NoError(t, err)
				require.JSONEq(t, string(data), string(newData))
			default:
				t.Fatal("wrong file")
			}
			v, err := openapi.NewValidator(o)
			require.NoError(t, err)
			require.NoError(t, v.ValidateSpec(
				openapi.DoNotValidateExamples(),
				openapi.AllowUndefinedTagsInOperation(),
			))
		})
	}
}

func TestValidator_ValidateSpec_ManuallyCreated(t *testing.T) {
	for _, tt := range []struct {
		name string
		spec *openapi.Extendable[openapi.OpenAPI]
		opts []openapi.SpecValidationOption
		err  string
	}{
		{
			name: "empty",
			spec: openapi.NewExtendable(&openapi.OpenAPI{}),
			err:  "openapi: required",
		},
		{
			name: "minimal valid with empty paths",
			spec: openapi.NewExtendable(&openapi.OpenAPI{
				OpenAPI: "3.1.0",
				Info: openapi.NewExtendable(&openapi.Info{
					Title:   "Minimal Valid Spec",
					Version: "1.0.0",
				}),
				Paths: openapi.NewExtendable[openapi.Paths](&openapi.Paths{}),
			}),
		},
		{
			name: "minimal valid with empty components",
			spec: openapi.NewExtendable(&openapi.OpenAPI{
				OpenAPI: "3.1.0",
				Info: openapi.NewExtendable(&openapi.Info{
					Title:   "Minimal Valid Spec",
					Version: "1.0.0",
				}),
				Components: openapi.NewExtendable[openapi.Components](&openapi.Components{}),
			}),
		},
		{
			name: "minimal valid with empty webhooks",
			spec: openapi.NewExtendable(&openapi.OpenAPI{
				OpenAPI: "3.1.0",
				Info: openapi.NewExtendable(&openapi.Info{
					Title:   "Minimal Valid Spec",
					Version: "1.0.0",
				}),
				WebHooks: make(map[string]*openapi.RefOrSpec[openapi.Extendable[openapi.PathItem]]),
			}),
		},
		{
			name: "xml component",
			spec: openapi.NewExtendable(&openapi.OpenAPI{
				OpenAPI: "3.1.0",
				Info: openapi.NewExtendable(&openapi.Info{
					Title:   "Minimal Valid Spec",
					Version: "1.0.0",
				}),
				Components: openapi.NewExtendable[openapi.Components]((&openapi.Components{}).WithRefOrSpec(
					"Person",
					&openapi.Schema{
						Type: openapi.NewSingleOrArray[string]("object"),
						Properties: map[string]*openapi.RefOrSpec[openapi.Schema]{
							"id": openapi.NewRefOrSpec[openapi.Schema](&openapi.Schema{
								Type:   openapi.NewSingleOrArray[string]("integer"),
								Format: "int32",
								XML: openapi.NewExtendable(&openapi.XML{
									Attribute: true,
								}),
							}),
							"name": openapi.NewRefOrSpec[openapi.Schema](&openapi.Schema{
								Type: openapi.NewSingleOrArray[string]("string"),
								XML: openapi.NewExtendable(&openapi.XML{
									Namespace: "https://example.com/schema/sample",
									Prefix:    "sample",
								}),
							}),
						},
					},
				)),
			}),
			opts: []openapi.SpecValidationOption{openapi.AllowUnusedComponents()},
		},
		{
			name: "properties examples",
			spec: openapi.NewExtendable(&openapi.OpenAPI{
				OpenAPI: "3.1.0",
				Info: openapi.NewExtendable(&openapi.Info{
					Title:   "Minimal Valid Spec",
					Version: "1.0.0",
				}),
				Components: openapi.NewExtendable[openapi.Components]((&openapi.Components{}).WithRefOrSpec(
					"Person",
					&openapi.Schema{
						Type: openapi.NewSingleOrArray[string]("object"),
						Properties: map[string]*openapi.RefOrSpec[openapi.Schema]{
							"id": openapi.NewRefOrSpec[openapi.Schema](&openapi.Schema{
								Type:   openapi.NewSingleOrArray[string]("integer"),
								Format: "int32",
							}),
							"name": openapi.NewRefOrSpec[openapi.Schema](&openapi.Schema{
								Type: openapi.NewSingleOrArray[string]("string"),
							}),
						},
						Examples: []any{
							map[string]any{
								"id":   123,
								"name": "John Doe 1",
							},
							struct {
								ID   int    `json:"id"`
								Name string `json:"name"`
							}{
								ID:   124,
								Name: "John Doe 2",
							},
						},
					},
				)),
			}),
			opts: []openapi.SpecValidationOption{openapi.AllowUnusedComponents()},
		},
		{
			name: "properties default",
			spec: openapi.NewExtendable(&openapi.OpenAPI{
				OpenAPI: "3.1.0",
				Info: openapi.NewExtendable(&openapi.Info{
					Title:   "Minimal Valid Spec",
					Version: "1.0.0",
				}),
				Components: openapi.NewExtendable[openapi.Components]((&openapi.Components{}).WithRefOrSpec(
					"Person",
					&openapi.Schema{
						Type: openapi.NewSingleOrArray[string]("object"),
						Properties: map[string]*openapi.RefOrSpec[openapi.Schema]{
							"id": openapi.NewRefOrSpec[openapi.Schema](&openapi.Schema{
								Type:    openapi.NewSingleOrArray[string]("integer"),
								Format:  "int32",
								Default: 42,
							}),
							"name": openapi.NewRefOrSpec[openapi.Schema](&openapi.Schema{
								Type:    openapi.NewSingleOrArray[string]("string"),
								Default: "John Doe",
							}),
						},
					},
				)),
			}),
			opts: []openapi.SpecValidationOption{openapi.AllowUnusedComponents()},
		},
		{
			name: "properties example",
			spec: openapi.NewExtendable(&openapi.OpenAPI{
				OpenAPI: "3.1.0",
				Info: openapi.NewExtendable(&openapi.Info{
					Title:   "Minimal Valid Spec",
					Version: "1.0.0",
				}),
				Components: openapi.NewExtendable[openapi.Components]((&openapi.Components{}).WithRefOrSpec(
					"Person",
					&openapi.Schema{
						Type: openapi.NewSingleOrArray[string]("object"),
						Properties: map[string]*openapi.RefOrSpec[openapi.Schema]{
							"id": openapi.NewRefOrSpec[openapi.Schema](&openapi.Schema{
								Type:    openapi.NewSingleOrArray[string]("integer"),
								Format:  "int32",
								Default: 42,
							}),
							"name": openapi.NewRefOrSpec[openapi.Schema](&openapi.Schema{
								Type:    openapi.NewSingleOrArray[string]("string"),
								Default: "John Doe",
							}),
						},
						Example: map[string]any{
							"id":   123,
							"name": "John Doe",
						},
					},
				)),
			}),
			opts: []openapi.SpecValidationOption{openapi.AllowUnusedComponents()},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			v, err := openapi.NewValidator(tt.spec)
			require.NoError(t, err)

			err = v.ValidateSpec(tt.opts...)
			t.Log("error: ", err)

			if tt.err == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tt.err)
			}
		})
	}
}

func TestNewValidator(t *testing.T) {
	data, err := os.ReadFile(path.Join("testdata", "petstore.json"))
	require.NoError(t, err)
	var petStore openapi.Extendable[openapi.OpenAPI]
	require.NoError(t, json.Unmarshal(data, &petStore))

	for _, tt := range []struct {
		name string
		spec *openapi.Extendable[openapi.OpenAPI]
	}{
		{
			name: "nil",
		},
		{
			name: "empty",
			spec: openapi.NewExtendable(&openapi.OpenAPI{}),
		},
		{
			name: "petstore",
			spec: &petStore,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := openapi.NewValidator(tt.spec)
			require.NoError(t, err)
		})
	}
}

func TestValidator_ValidateData(t *testing.T) {
	data, err := os.ReadFile(path.Join("testdata", "petstore.json"))
	require.NoError(t, err)
	var spec openapi.Extendable[openapi.OpenAPI]
	require.NoError(t, json.Unmarshal(data, &spec))
	validator, err := openapi.NewValidator(&spec)
	require.NoError(t, err)

	for _, tt := range []struct {
		name          string
		ref           string
		data          string
		compileError  string
		validateError string
	}{
		{
			name: "by component",
			ref:  "#/components/schemas/Pet",
			data: `{"id": 123, "name": "foo", "tag": "bar"}`,
		},
		{
			name:          "by component failed",
			ref:           "/components/schemas/Pet",
			data:          `{"id": "123", "name": "foo", "tag": "bar"}`,
			validateError: "got string, want integer",
		},
		{
			name: "by route",
			ref:  "/paths/~1pets~1{petId}/get/responses/200/content/application~1json/schema",
			data: `{"id": 123, "name": "foo", "tag": "bar"}`,
		},
		{
			name:          "by route failed",
			ref:           "/paths/~1pets~1{petId}/get/responses/200/content/application~1json/schema",
			data:          `{"id": "123", "name": "foo", "tag": "bar"}`,
			validateError: "got string, want integer",
		},
		{
			name:         "component not found",
			ref:          "/components/schemas/Fake",
			data:         `{}`,
			compileError: "not found",
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var data any
			require.NoError(t, json.Unmarshal([]byte(tt.data), &data))
			err := validator.ValidateData(tt.ref, data)

			if tt.compileError != "" {
				require.ErrorContains(t, err, tt.compileError)
				return
			}

			if tt.validateError != "" {
				require.ErrorContains(t, err, tt.validateError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

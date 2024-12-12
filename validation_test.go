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

func TestValidation(t *testing.T) {
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
			require.NoError(t, openapi.ValidateSpec(o, openapi.DoNotValidateExamples(), openapi.DoNotValidateDefaultValues(), openapi.AllowUndefinedTagsInOperation()))
		})
	}
}

func TestManuallyCreatedSpec(t *testing.T) {
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
				Paths: openapi.NewExtendable(&openapi.Paths{
					Paths: map[string]*openapi.RefOrSpec[openapi.Extendable[openapi.PathItem]]{
						"/persons": openapi.NewRefOrExtSpec[openapi.PathItem](&openapi.PathItem{
							Get: openapi.NewExtendable(&openapi.Operation{
								Responses: openapi.NewExtendable(&openapi.Responses{
									Default: openapi.NewRefOrExtSpec[openapi.Response](&openapi.Response{
										Description: "A person",
										Content: map[string]*openapi.Extendable[openapi.MediaType]{
											"application/json": openapi.NewExtendable(&openapi.MediaType{
												Schema: openapi.NewRefOrSpec[openapi.Schema]("#/components/schemas/Person"),
											}),
										},
									}),
								}),
							}),
						}),
					},
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
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			err := openapi.ValidateSpec(tt.spec, tt.opts...)
			t.Log("error: ", err)
			if tt.err == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.err)
			}
		})
	}
}

func TestValidatePayload(t *testing.T) {
	data, err := os.ReadFile(path.Join("testdata", "petstore.json"))
	require.NoError(t, err)
	var spec openapi.Extendable[openapi.OpenAPI]
	require.NoError(t, json.Unmarshal(data, &spec))
	validator, err := openapi.NewDataValidator(&spec)
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
			err := validator.Validate(tt.ref, data)

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

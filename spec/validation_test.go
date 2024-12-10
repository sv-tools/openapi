package spec_test

import (
	"encoding/json"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/sv-tools/openapi/spec"
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
			var o *spec.Extendable[spec.OpenAPI]
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
			require.NoError(t, spec.ValidateSpec(o, spec.DoNotValidateExamples(), spec.DoNotValidateDefaultValues(), spec.AllowUndefinedTagsInOperation()))
		})
	}
}

func TestManuallyCreatedSpec(t *testing.T) {
	for _, tt := range []struct {
		name string
		spec *spec.Extendable[spec.OpenAPI]
		opts []spec.ValidationOption
		err  string
	}{
		{
			name: "empty",
			spec: spec.NewExtendable(&spec.OpenAPI{}),
			err:  "openapi: required",
		},
		{
			name: "minimal valid with empty paths",
			spec: spec.NewExtendable(&spec.OpenAPI{
				OpenAPI: "3.1.0",
				Info: spec.NewExtendable(&spec.Info{
					Title:   "Minimal Valid Spec",
					Version: "1.0.0",
				}),
				Paths: spec.NewExtendable[spec.Paths](&spec.Paths{}),
			}),
		},
		{
			name: "minimal valid with empty components",
			spec: spec.NewExtendable(&spec.OpenAPI{
				OpenAPI: "3.1.0",
				Info: spec.NewExtendable(&spec.Info{
					Title:   "Minimal Valid Spec",
					Version: "1.0.0",
				}),
				Components: spec.NewExtendable[spec.Components](&spec.Components{}),
			}),
		},
		{
			name: "minimal valid with empty webhooks",
			spec: spec.NewExtendable(&spec.OpenAPI{
				OpenAPI: "3.1.0",
				Info: spec.NewExtendable(&spec.Info{
					Title:   "Minimal Valid Spec",
					Version: "1.0.0",
				}),
				WebHooks: make(map[string]*spec.RefOrSpec[spec.Extendable[spec.PathItem]]),
			}),
		},
		{
			name: "xml component",
			spec: spec.NewExtendable(&spec.OpenAPI{
				OpenAPI: "3.1.0",
				Info: spec.NewExtendable(&spec.Info{
					Title:   "Minimal Valid Spec",
					Version: "1.0.0",
				}),
				Paths: spec.NewExtendable(&spec.Paths{
					Paths: map[string]*spec.RefOrSpec[spec.Extendable[spec.PathItem]]{
						"/persons": spec.NewRefOrExtSpec[spec.PathItem](&spec.PathItem{
							Get: spec.NewExtendable(&spec.Operation{
								Responses: spec.NewExtendable(&spec.Responses{
									Default: spec.NewRefOrExtSpec[spec.Response](&spec.Response{
										Description: "A person",
										Content: map[string]*spec.Extendable[spec.MediaType]{
											"application/json": spec.NewExtendable(&spec.MediaType{
												Schema: spec.NewRefOrSpec[spec.Schema]("#/components/schemas/Person"),
											}),
										},
									}),
								}),
							}),
						}),
					},
				}),
				Components: spec.NewExtendable[spec.Components]((&spec.Components{}).WithRefOrSpec(
					"Person",
					&spec.Schema{
						Type: spec.NewSingleOrArray[string]("object"),
						Properties: map[string]*spec.RefOrSpec[spec.Schema]{
							"id": spec.NewRefOrSpec[spec.Schema](&spec.Schema{
								Type:   spec.NewSingleOrArray[string]("integer"),
								Format: "int32",
								XML: spec.NewExtendable(&spec.XML{
									Attribute: true,
								}),
							}),
							"name": spec.NewRefOrSpec[spec.Schema](&spec.Schema{
								Type: spec.NewSingleOrArray[string]("string"),
								XML: spec.NewExtendable(&spec.XML{
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
			err := spec.ValidateSpec(tt.spec, tt.opts...)
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

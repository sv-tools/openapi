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
			spec: spec.NewOpenAPI(),
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
				Paths: spec.NewPaths(),
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
				Components: spec.NewComponents(),
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
				Components: spec.NewExtendable(&spec.Components{
					Schemas: map[string]*spec.RefOrSpec[spec.Schema]{
						"Person": spec.NewSchemaSpec((&spec.Schema{}).
							WithType("object").
							WithProperty("id", spec.NewSchemaSpec((&spec.Schema{}).
								WithType("integer").
								WithFormat("int32").
								WithXML(spec.NewXML((&spec.XML{}).
									WithAttribute(true),
								)),
							)).
							WithProperty("name", spec.NewSchemaSpec((&spec.Schema{}).
								WithType("string").
								WithXML(spec.NewXML((&spec.XML{}).
									WithNamespace("https://example.com/schema/sample").
									WithPrefix("sample"),
								)),
							)),
						),
					},
				}),
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

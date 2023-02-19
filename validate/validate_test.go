package validate_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path"
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v5"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/sv-tools/openapi/spec"
	"github.com/sv-tools/openapi/validate"
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
				require.NoError(t, validate.Yaml(data))
				require.NoError(t, yaml.Unmarshal(data, &o))
				newData, err := yaml.Marshal(&o)
				require.NoError(t, err)
				require.YAMLEq(t, string(data), string(newData))
			case ".json":
				require.NoError(t, validate.Json(data))
				require.NoError(t, json.Unmarshal(data, &o))
				newData, err := json.Marshal(&o)
				require.NoError(t, err)
				require.JSONEq(t, string(data), string(newData))
			default:
				t.Fatal("wrong file")
			}
			require.NoError(t, validate.Spec(o))
		})
	}
}

func TestManuallyCreatedSpec(t *testing.T) {
	minSpec := spec.NewOpenAPI()
	minSpec.Spec.OpenAPI = "3.1.0"
	minSpec.Spec.Info = spec.NewInfo()
	minSpec.Spec.Info.Spec.Title = "Minimal Valid Spec"
	minSpec.Spec.Info.Spec.Version = "1.0.0"
	minSpec.Spec.Paths = spec.NewPaths()

	for _, tt := range []struct {
		name string
		spec *spec.Extendable[spec.OpenAPI]
		err  string
	}{
		{
			name: "empty",
			spec: spec.NewOpenAPI(),
			err:  "does not validate",
		},
		{
			name: "minimal valid",
			spec: minSpec,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Spec(tt.spec)
			t.Log("report:", validate.Report(err, false))
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
	c := jsonschema.NewCompiler()
	err = c.AddResource("https://petstore.swagger.io/v1", bytes.NewBuffer(data))
	require.NoError(t, err)

	for _, tt := range []struct {
		name          string
		ref           string
		compileError  string
		validateError string
	}{
		{
			name:          "by component",
			ref:           "https://petstore.swagger.io/v1#/components/schemas/Pet",
			validateError: "expected integer, but got string",
		},
		{
			name:          "by route",
			ref:           "https://petstore.swagger.io/v1#/paths/%2fpets%2f%7BpetId%7D/get/responses/200/content/application%2fjson/schema",
			validateError: "expected integer, but got string",
		},
		{
			name:         "not found",
			ref:          "https://petstore.swagger.io/v1#/components/schemas/Fake",
			compileError: "Fake not found",
		},
		{
			name:         "wrong url",
			ref:          "https://petstore.swagger.io/v2#/components/schemas/Fake",
			compileError: "no Loader found for",
		},
		{
			name:         "not absolute url",
			ref:          "#/components/schemas/Fake",
			compileError: "no Loader found for",
		},
		{
			name:         "not absolute url",
			ref:          "/components/schemas/Fake",
			compileError: "no Loader found for",
		},
		{
			name:         "not absolute url",
			ref:          "components/schemas/Fake",
			compileError: "no Loader found for",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			s, err := c.Compile(tt.ref)
			if tt.compileError != "" {
				require.ErrorContains(t, err, tt.compileError)
				return
			}
			require.NoError(t, err)

			var v1 any
			err = json.Unmarshal([]byte(`{"id": 123, "name": "foo", "tag": "bar"}`), &v1)
			require.NoError(t, err)
			err = s.Validate(v1)
			require.NoError(t, err)

			var v2 any
			err = json.Unmarshal([]byte(`{"id": "123", "name": "foo", "tag": "bar"}`), &v2)
			require.NoError(t, err)
			err = s.Validate(v2)
			require.Error(t, err)
			require.ErrorContains(t, err, tt.validateError)
		})
	}
}

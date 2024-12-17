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
			v, err := openapi.NewValidator(
				o,
				openapi.AllowUndefinedTagsInOperation(),
				openapi.ValidateStringDataAsJSON(),
			)
			require.NoError(t, err)
			require.NoError(t, v.ValidateSpec())
		})
	}
}

func TestValidator_ValidateSpec_ManuallyCreated(t *testing.T) {
	for _, tt := range []struct {
		name string
		spec *openapi.Extendable[openapi.OpenAPI]
		opts []openapi.ValidationOption
		err  string
	}{
		{
			name: "info required",
			spec: openapi.NewOpenAPIBuilder().Build(),
			err:  "/info: required",
		},
		{
			name: "any of path or webhooks or components required",
			spec: openapi.NewOpenAPIBuilder().Build(),
			err:  "/paths||webhooks||components: required",
		},
		{
			name: "minimal valid with empty paths",
			spec: openapi.NewOpenAPIBuilder().Info(
				openapi.NewInfoBuilder().
					Title("Minimal Valid Spec").
					Version("1.0.0").
					Build(),
			).Paths(openapi.NewPaths()).Build(),
		},
		{
			name: "minimal valid with empty components",
			spec: openapi.NewOpenAPIBuilder().Info(
				openapi.NewInfoBuilder().
					Title("Minimal Valid Spec").
					Version("1.0.0").
					Build(),
			).Components(openapi.NewComponents()).Build(),
		},
		{
			name: "minimal valid with empty webhooks",
			spec: openapi.NewOpenAPIBuilder().Info(
				openapi.NewInfoBuilder().
					Title("Minimal Valid Spec").
					Version("1.0.0").
					Build(),
			).WebHooks(openapi.NewWebhooks()).Build(),
		},
		{
			name: "xml component",
			spec: openapi.NewOpenAPIBuilder().Info(
				openapi.NewInfoBuilder().
					Title("Minimal Valid Spec").
					Version("1.0.0").
					Build(),
			).AddComponent("Person", openapi.NewSchemaBuilder().
				AddType("object").
				AddProperty("id", openapi.NewSchemaBuilder().
					AddType("integer").
					Format("int32").
					XML(openapi.NewXMLBuilder().Attribute(true).Build()).
					Build(),
				).
				AddProperty("name", openapi.NewSchemaBuilder().
					AddType("string").
					XML(openapi.NewXMLBuilder().
						Namespace("https://example.com/schema/sample").
						Prefix("sample").
						Build(),
					).
					Build(),
				).
				Build(),
			).Build(),
			opts: []openapi.ValidationOption{openapi.AllowUnusedComponents()},
		},
		{
			name: "properties examples",
			spec: openapi.NewOpenAPIBuilder().Info(
				openapi.NewInfoBuilder().
					Title("Minimal Valid Spec").
					Version("1.0.0").
					Build(),
			).AddComponent("Person", openapi.NewSchemaBuilder().
				AddType("object").
				AddProperty("id", openapi.NewSchemaBuilder().
					AddType("integer").
					Format("int32").
					Build(),
				).
				AddProperty("name", openapi.NewSchemaBuilder().
					AddType("string").
					Build(),
				).
				AddExamples(
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
				).Build(),
			).Build(),
			opts: []openapi.ValidationOption{openapi.AllowUnusedComponents()},
		},
		{
			name: "properties examples error",
			spec: openapi.NewOpenAPIBuilder().Info(
				openapi.NewInfoBuilder().
					Title("Minimal Valid Spec").
					Version("1.0.0").
					Build(),
			).AddComponent("Person", openapi.NewSchemaBuilder().
				AddType("object").
				AddProperty("id", openapi.NewSchemaBuilder().
					AddType("integer").
					Format("int32").
					Build(),
				).
				AddProperty("name", openapi.NewSchemaBuilder().
					AddType("string").
					Build(),
				).
				AddExamples(
					map[string]any{
						"id":   "123",
						"name": false,
					},
				).Build(),
			).Build(),
			opts: []openapi.ValidationOption{openapi.AllowUnusedComponents()},
			err:  "at '/id': got string, want integer",
		},
		{
			name: "properties default",
			spec: openapi.NewOpenAPIBuilder().Info(
				openapi.NewInfoBuilder().
					Title("Minimal Valid Spec").
					Version("1.0.0").
					Build(),
			).AddComponent("Person", openapi.NewSchemaBuilder().
				AddType("object").
				AddProperty("id", openapi.NewSchemaBuilder().
					AddType("integer").
					Format("int32").
					Default(42).
					Build(),
				).
				AddProperty("name", openapi.NewSchemaBuilder().
					AddType("string").
					Default("John Doe").
					Build(),
				).Build(),
			).Build(),
			opts: []openapi.ValidationOption{openapi.AllowUnusedComponents()},
		},
		{
			name: "properties default error",
			spec: openapi.NewOpenAPIBuilder().Info(
				openapi.NewInfoBuilder().
					Title("Minimal Valid Spec").
					Version("1.0.0").
					Build(),
			).AddComponent("Person", openapi.NewSchemaBuilder().
				AddType("object").
				AddProperty("id", openapi.NewSchemaBuilder().
					AddType("integer").
					Format("int32").
					Default("42").
					Build(),
				).
				AddProperty("name", openapi.NewSchemaBuilder().
					AddType("string").
					Default(false).
					Build(),
				).Build(),
			).Build(),
			opts: []openapi.ValidationOption{openapi.AllowUnusedComponents()},
			err:  "at '': got string, want integer",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			v, err := openapi.NewValidator(tt.spec, tt.opts...)
			require.NoError(t, err)

			err = v.ValidateSpec()
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

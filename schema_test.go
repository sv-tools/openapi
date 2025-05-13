package openapi_test

import (
	"encoding/json"
	"testing"

	goyaml "github.com/goccy/go-yaml"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/sv-tools/openapi"
)

func TestSchema_Marshal_Unmarshal(t *testing.T) {
	for _, tt := range []struct {
		name            string
		data            string
		expected        *openapi.Schema
		emptyExtensions bool
	}{
		{
			name:            "spec only",
			data:            `{"title": "foo"}`,
			expected:        openapi.NewSchemaBuilder().Title("foo").Build().Spec,
			emptyExtensions: true,
		},
		{
			name:            "spec with extension field",
			data:            `{"title": "foo", "b": "bar"}`,
			expected:        openapi.NewSchemaBuilder().Title("foo").AddExt("b", "bar").Build().Spec,
			emptyExtensions: false,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Run("json", func(t *testing.T) {
				var v openapi.Schema
				require.NoError(t, json.Unmarshal([]byte(tt.data), &v))
				if tt.emptyExtensions {
					require.Empty(t, v.Extensions)
				} else {
					require.NotEmpty(t, v.Extensions)
				}
				data, err := json.Marshal(&v)
				require.NoError(t, err)
				require.JSONEq(t, tt.data, string(data))
				require.Equal(t, *tt.expected, v)
			})
			t.Run("yaml.v3", func(t *testing.T) {
				var v openapi.Schema
				require.NoError(t, yaml.Unmarshal([]byte(tt.data), &v))
				if tt.emptyExtensions {
					require.Empty(t, v.Extensions)
				} else {
					require.NotEmpty(t, v.Extensions)
				}
				data, err := yaml.Marshal(&v)
				require.NoError(t, err)
				require.YAMLEq(t, tt.data, string(data))
				require.Equal(t, *tt.expected, v)
			})
			t.Run("goccy/go-yaml", func(t *testing.T) {
				var v openapi.Schema
				require.NoError(t, goyaml.Unmarshal([]byte(tt.data), &v))
				if tt.emptyExtensions {
					require.Empty(t, v.Extensions)
				} else {
					require.NotEmpty(t, v.Extensions)
				}
				data, err := goyaml.Marshal(&v)
				require.NoError(t, err)
				require.YAMLEq(t, tt.data, string(data))
				require.Equal(t, *tt.expected, v)
			})
		})
	}
}

func TestSchema_AddExt(t *testing.T) {
	for _, tt := range []struct {
		name     string
		key      string
		value    any
		expected map[string]any
	}{
		{
			name:  "without prefix",
			key:   "foo",
			value: 42,
			expected: map[string]any{
				"foo": 42,
			},
		},
		{
			name:  "with prefix",
			key:   "x-foo",
			value: 43,
			expected: map[string]any{
				"x-foo": 43,
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ext := openapi.Schema{}
			ext.AddExt(tt.key, tt.value)
			require.Equal(t, tt.expected, ext.Extensions)
		})
	}
}

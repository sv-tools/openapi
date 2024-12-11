package openapi_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/sv-tools/openapi"
)

type testExtendable struct {
	A string `json:"a,omitempty" yaml:"a,omitempty"`
}

func TestExtendable_Marshal_Unmarshal(t *testing.T) {
	for _, tt := range []struct {
		name            string
		data            string
		expected        string
		emptyExtensions bool
	}{
		{
			name:            "spec only",
			data:            `{"a": "foo"}`,
			emptyExtensions: true,
		},
		{
			name:            "spec with extra non extension field",
			data:            `{"a": "foo", "b": "bar"}`,
			expected:        `{"a": "foo"}`,
			emptyExtensions: true,
		},
		{
			name:            "spec with extension field",
			data:            `{"a": "foo", "x-b": "bar"}`,
			emptyExtensions: false,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Run("json", func(t *testing.T) {
				var v *openapi.Extendable[testExtendable]
				require.NoError(t, json.Unmarshal([]byte(tt.data), &v))
				if tt.emptyExtensions {
					require.Empty(t, v.Extensions)
				} else {
					require.NotEmpty(t, v.Extensions)
				}
				data, err := json.Marshal(&v)
				require.NoError(t, err)
				if tt.expected == "" {
					tt.expected = tt.data
				}
				require.JSONEq(t, tt.expected, string(data))
			})
			t.Run("yaml", func(t *testing.T) {
				var v *openapi.Extendable[testExtendable]
				require.NoError(t, yaml.Unmarshal([]byte(tt.data), &v))
				if tt.emptyExtensions {
					require.Empty(t, v.Extensions)
				} else {
					require.NotEmpty(t, v.Extensions)
				}
				data, err := yaml.Marshal(&v)
				require.NoError(t, err)
				if tt.expected == "" {
					tt.expected = tt.data
				}
				require.YAMLEq(t, tt.expected, string(data))
			})
		})
	}
}

func TestExtendable_WithExt(t *testing.T) {
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
				"x-foo": 42,
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
			ext := openapi.NewExtendable(&testExtendable{})
			ext.WithExt(tt.key, tt.value)
			require.Equal(t, tt.expected, ext.Extensions)
		})
	}
}

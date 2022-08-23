package spec_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/sv-tools/openapi/spec"
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
				var v *spec.Extendable[testExtendable]
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
				var v *spec.Extendable[testExtendable]
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

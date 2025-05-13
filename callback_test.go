package openapi_test

import (
	"encoding/json"
	"testing"

	goyaml "github.com/goccy/go-yaml"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/sv-tools/openapi"
)

func TestCallback_Marshal_Unmarshal(t *testing.T) {
	for _, tt := range []struct {
		name     string
		data     string
		expected string
	}{
		{
			name: "spec",
			data: `{"example.com": {"get": {"summary": "foo"}}}`,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Run("json", func(t *testing.T) {
				var v openapi.Callback
				require.NoError(t, json.Unmarshal([]byte(tt.data), &v))
				data, err := json.Marshal(&v)
				require.NoError(t, err)
				if tt.expected == "" {
					tt.expected = tt.data
				}
				require.JSONEq(t, tt.expected, string(data))
			})
			t.Run("yaml.v3", func(t *testing.T) {
				var v openapi.Callback
				require.NoError(t, yaml.Unmarshal([]byte(tt.data), &v))
				data, err := yaml.Marshal(&v)
				require.NoError(t, err)
				if tt.expected == "" {
					tt.expected = tt.data
				}
				require.YAMLEq(t, tt.expected, string(data))
			})
			t.Run("goccy/go-yaml", func(t *testing.T) {
				var v openapi.Callback
				require.NoError(t, goyaml.Unmarshal([]byte(tt.data), &v))
				data, err := goyaml.Marshal(&v)
				require.NoError(t, err)
				if tt.expected == "" {
					tt.expected = tt.data
				}
				require.YAMLEq(t, tt.expected, string(data))
			})
		})
	}
}

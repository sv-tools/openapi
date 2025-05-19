package openapi_test

import (
	"encoding/json"
	"testing"

	goyaml "github.com/goccy/go-yaml"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/sv-tools/openapi"
)

func TestResponses_Marshal_Unmarshal(t *testing.T) {
	for _, tt := range []struct {
		name string
		data string
	}{
		{
			name: "response with default",
			data: `{"200": {"description": "foo"}, "default": {"description": "bar"}}`,
		},
		{
			name: "response with default only",
			data: `{"default": {"description": "bar"}}`,
		},
		{
			name: "response without default",
			data: `{"200": {"description": "foo"}}`,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Run("json", func(t *testing.T) {
				var v openapi.Responses
				require.NoError(t, json.Unmarshal([]byte(tt.data), &v))
				data, err := json.Marshal(&v)
				require.NoError(t, err)
				require.JSONEq(t, tt.data, string(data))
			})
			t.Run("yaml.v3", func(t *testing.T) {
				var v openapi.Responses
				require.NoError(t, yaml.Unmarshal([]byte(tt.data), &v))
				data, err := yaml.Marshal(&v)
				require.NoError(t, err)
				require.YAMLEq(t, tt.data, string(data))
			})
			t.Run("goccy/go-yaml", func(t *testing.T) {
				var v openapi.Responses
				require.NoError(t, goyaml.Unmarshal([]byte(tt.data), &v))
				data, err := goyaml.Marshal(&v)
				require.NoError(t, err)
				require.YAMLEq(t, tt.data, string(data))
			})
		})
	}
}

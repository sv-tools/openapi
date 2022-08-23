package spec_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/sv-tools/openapi/spec"
)

type testAD struct {
	AP   *spec.BoolOrSchema `json:"ap,omitempty" yaml:"ap,omitempty"`
	Name string             `json:"name,omitempty" yaml:"name,omitempty"`
}

func TestAdditionalPropertiesJSON(t *testing.T) {
	for _, tt := range []struct {
		name      string
		data      string
		nilAP     bool
		allowed   bool
		nilSchema bool
	}{
		{
			name:  "no AP",
			data:  `{"name": "foo"}`,
			nilAP: true,
		},
		{
			name:      "false",
			data:      `{"name": "foo", "ap": false}`,
			nilAP:     false,
			allowed:   false,
			nilSchema: true,
		},
		{
			name:      "true",
			data:      `{"name": "foo", "ap": true}`,
			nilAP:     false,
			allowed:   true,
			nilSchema: true,
		},
		{
			name:      "schema",
			data:      `{"name": "foo", "ap": {"title": "bar", "description": "test"}}`,
			nilAP:     false,
			allowed:   true,
			nilSchema: false,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Run("json", func(t *testing.T) {
				var v testAD
				require.NoError(t, json.Unmarshal([]byte(tt.data), &v))
				require.Equal(t, "foo", v.Name)
				if tt.nilAP {
					require.Nil(t, v.AP)
				} else {
					require.NotNil(t, v.AP)
					require.Equal(t, tt.allowed, v.AP.Allowed)
					require.Equal(t, tt.nilSchema, v.AP.Schema == nil)
				}
				newJson, err := json.Marshal(&v)
				require.NoError(t, err)
				require.JSONEq(t, tt.data, string(newJson))
			})

			t.Run("yaml", func(t *testing.T) {
				var v testAD
				require.NoError(t, yaml.Unmarshal([]byte(tt.data), &v))
				require.Equal(t, "foo", v.Name)
				if tt.nilAP {
					require.Nil(t, v.AP)
				} else {
					require.NotNil(t, v.AP)
					require.Equal(t, tt.allowed, v.AP.Allowed)
					require.Equal(t, tt.nilSchema, v.AP.Schema == nil)
				}
				newYaml, err := yaml.Marshal(&v)
				require.NoError(t, err)
				require.YAMLEq(t, tt.data, string(newYaml))
			})
		})
	}
}

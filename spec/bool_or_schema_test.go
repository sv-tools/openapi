package spec_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/sv-tools/openapi/spec"
)

type testAD struct {
	Name string             `json:"name,omitempty" yaml:"name,omitempty"`
	AP   *spec.BoolOrSchema `json:"ap,omitempty" yaml:"ap,omitempty"`
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
			var j testAD
			require.NoError(t, json.Unmarshal([]byte(tt.data), &j))
			require.Equal(t, "foo", j.Name)
			if tt.nilAP {
				require.Nil(t, j.AP)
			} else {
				require.NotNil(t, j.AP)
				require.Equal(t, tt.allowed, j.AP.Allowed)
				require.Equal(t, tt.nilSchema, j.AP.Schema == nil)
			}
			newJson, err := json.Marshal(&j)
			require.NoError(t, err)
			require.JSONEq(t, tt.data, string(newJson))

			var y testAD
			require.NoError(t, yaml.Unmarshal([]byte(tt.data), &y))
			require.Equal(t, "foo", y.Name)
			if tt.nilAP {
				require.Nil(t, y.AP)
			} else {
				require.NotNil(t, y.AP)
				require.Equal(t, tt.allowed, y.AP.Allowed)
				require.Equal(t, tt.nilSchema, y.AP.Schema == nil)
			}
			newYaml, err := yaml.Marshal(&y)
			require.NoError(t, err)
			require.YAMLEq(t, tt.data, string(newYaml))
		})
	}

}

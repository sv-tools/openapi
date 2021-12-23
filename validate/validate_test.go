package validate_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/sv-tools/openapi/spec"
	"github.com/sv-tools/openapi/validate"
)

func TestValidation(t *testing.T) {
	info, err := ioutil.ReadDir("testdata")
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

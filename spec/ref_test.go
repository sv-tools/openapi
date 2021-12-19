package spec_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/sv-tools/openapi/spec"
)

type testRefOrSpec struct {
	A string `json:"a,omitempty" yaml:"a,omitempty"`
	B string `json:"b,omitempty" yaml:"b,omitempty"`
}

func (o testRefOrSpec) OpenAPIConstraint() {}

func TestNewRefOrSpec(t *testing.T) {
	for _, tt := range []struct {
		name    string
		ref     *spec.Ref
		spec    *testRefOrSpec
		nilRef  bool
		nilSpec bool
	}{
		{
			name:    "empty",
			nilRef:  true,
			nilSpec: true,
		},
		{
			name:    "ref",
			ref:     spec.NewRef("foo"),
			nilRef:  false,
			nilSpec: true,
		},
		{
			name:    "spec",
			spec:    &testRefOrSpec{A: "foo", B: "bar"},
			nilRef:  true,
			nilSpec: false,
		},
		{
			name:    "both",
			ref:     spec.NewRef("foo"),
			spec:    &testRefOrSpec{A: "foo", B: "bar"},
			nilRef:  false,
			nilSpec: true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			o := spec.NewRefOrSpec(tt.ref, tt.spec)
			require.NotNil(t, o)
			if tt.nilRef {
				require.Nil(t, o.Ref)
			} else {
				require.Equal(t, tt.ref, o.Ref)
			}
			if tt.nilSpec {
				require.Nil(t, o.Spec)
			} else {
				require.Equal(t, tt.spec, o.Spec)
			}
		})
	}
}

func TestRefOrSpec_Marshal_Unmarshal(t *testing.T) {
	for _, tt := range []struct {
		name     string
		data     string
		expected string
		nilRef   bool
		nilSpec  bool
	}{
		{
			name:    "ref",
			data:    `{"$ref": "foo"}`,
			nilRef:  false,
			nilSpec: true,
		},
		{
			name:    "spec",
			data:    `{"a": "foo", "b": "bar"}`,
			nilRef:  true,
			nilSpec: false,
		},
		{
			name:     "both",
			data:     `{"$ref": "foo", "a": "foo", "b": "bar"}`,
			expected: `{"$ref": "foo"}`,
			nilRef:   false,
			nilSpec:  true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Run("json", func(t *testing.T) {
				var j *spec.RefOrSpec[testRefOrSpec]
				require.NoError(t, json.Unmarshal([]byte(tt.data), &j))
				if tt.nilRef {
					require.Nil(t, j.Ref)
				} else {
					require.NotNil(t, j.Ref)
				}
				if tt.nilSpec {
					require.Nil(t, j.Spec)
				} else {
					require.NotNil(t, j.Spec)
				}
				data, err := json.Marshal(&j)
				require.NoError(t, err)
				if tt.expected == "" {
					tt.expected = tt.data
				}
				require.JSONEq(t, tt.expected, string(data))
			})

			t.Run("yaml", func(t *testing.T) {
				var y *spec.RefOrSpec[testRefOrSpec]
				require.NoError(t, yaml.Unmarshal([]byte(tt.data), &y))
				if tt.nilRef {
					require.Nil(t, y.Ref)
				} else {
					require.NotNil(t, y.Ref)
				}
				if tt.nilSpec {
					require.Nil(t, y.Spec)
				} else {
					require.NotNil(t, y.Spec)
				}
				data, err := yaml.Marshal(&y)
				require.NoError(t, err)
				if tt.expected == "" {
					tt.expected = tt.data
				}
				require.YAMLEq(t, tt.expected, string(data))
			})
		})
	}
}

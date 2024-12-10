package spec_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/sv-tools/openapi/spec"
)

type testRefOrSpec struct {
	A string `json:"a,omitempty" yaml:"a,omitempty"`
	B string `json:"b,omitempty" yaml:"b,omitempty"`
}

func TestNewRefOrSpec(t *testing.T) {
	for _, tt := range []struct {
		ref_or_spec any
		name        string
		nilRef      bool
		nilSpec     bool
	}{
		{
			name:    "empty",
			nilRef:  true,
			nilSpec: true,
		},
		{
			name:        "ref by reference",
			ref_or_spec: &spec.Ref{Ref: "foo"},
			nilRef:      false,
			nilSpec:     true,
		},
		{
			name:        "ref by value",
			ref_or_spec: spec.Ref{Ref: "foo"},
			nilRef:      false,
			nilSpec:     true,
		},
		{
			name:        "string",
			ref_or_spec: "foo",
			nilRef:      false,
			nilSpec:     true,
		},
		{
			name:        "spec by reference",
			ref_or_spec: &testRefOrSpec{A: "foo", B: "bar"},
			nilRef:      true,
			nilSpec:     false,
		},
		{
			name:        "spec by value",
			ref_or_spec: testRefOrSpec{A: "foo", B: "bar"},
			nilRef:      true,
			nilSpec:     false,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			o := spec.NewRefOrSpec[testRefOrSpec](tt.ref_or_spec)
			require.NotNil(t, o)
			if tt.nilRef {
				require.Nil(t, o.Ref)
			} else {
				switch v := tt.ref_or_spec.(type) {
				case string:
					require.Equal(t, v, o.Ref.Ref)
				case *spec.Ref:
					require.Equal(t, v, o.Ref)
				case spec.Ref:
					require.Equal(t, &v, o.Ref)
				default:
					t.Fatal("unexpected ref type")
				}
			}
			if tt.nilSpec {
				require.Nil(t, o.Spec)
			} else {
				switch v := tt.ref_or_spec.(type) {
				case *testRefOrSpec:
					require.Equal(t, v, o.Spec)
				case testRefOrSpec:
					require.Equal(t, &v, o.Spec)
				default:
					t.Fatal("unexpected spec type")
				}
			}
		})
	}
}

func TestNewRefOrExtSpec(t *testing.T) {
	for _, tt := range []struct {
		ref_or_spec any
		name        string
		nilRef      bool
		nilSpec     bool
	}{
		{
			name:    "empty",
			nilRef:  true,
			nilSpec: true,
		},
		{
			name:        "ref by reference",
			ref_or_spec: &spec.Ref{Ref: "foo"},
			nilRef:      false,
			nilSpec:     true,
		},
		{
			name:        "ref by value",
			ref_or_spec: spec.Ref{Ref: "foo"},
			nilRef:      false,
			nilSpec:     true,
		},
		{
			name:        "string",
			ref_or_spec: "foo",
			nilRef:      false,
			nilSpec:     true,
		},
		{
			name:        "spec by reference",
			ref_or_spec: &testRefOrSpec{A: "foo", B: "bar"},
			nilRef:      true,
			nilSpec:     false,
		},
		{
			name:        "spec by value",
			ref_or_spec: testRefOrSpec{A: "foo", B: "bar"},
			nilRef:      true,
			nilSpec:     false,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			o := spec.NewRefOrExtSpec[testRefOrSpec](tt.ref_or_spec)
			require.NotNil(t, o)
			if tt.nilRef {
				require.Nil(t, o.Ref)
			} else {
				switch v := tt.ref_or_spec.(type) {
				case string:
					require.Equal(t, v, o.Ref.Ref)
				case *spec.Ref:
					require.Equal(t, v, o.Ref)
				case spec.Ref:
					require.Equal(t, &v, o.Ref)
				default:
					t.Fatal("unexpected ref type")
				}
			}
			if tt.nilSpec {
				require.Nil(t, o.Spec)
			} else {
				require.IsType(t, &spec.Extendable[testRefOrSpec]{}, o.Spec)
				switch v := tt.ref_or_spec.(type) {
				case *testRefOrSpec:
					require.Equal(t, v, o.Spec.Spec)
				case testRefOrSpec:
					require.Equal(t, &v, o.Spec.Spec)
				default:
					t.Fatal("unexpected spec type")
				}
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
				var v *spec.RefOrSpec[testRefOrSpec]
				require.NoError(t, json.Unmarshal([]byte(tt.data), &v))
				if tt.nilRef {
					require.Nil(t, v.Ref)
				} else {
					require.NotNil(t, v.Ref)
				}
				if tt.nilSpec {
					require.Nil(t, v.Spec)
				} else {
					require.NotNil(t, v.Spec)
				}
				data, err := json.Marshal(&v)
				require.NoError(t, err)
				if tt.expected == "" {
					tt.expected = tt.data
				}
				require.JSONEq(t, tt.expected, string(data))
			})

			t.Run("yaml", func(t *testing.T) {
				var v *spec.RefOrSpec[testRefOrSpec]
				require.NoError(t, yaml.Unmarshal([]byte(tt.data), &v))
				if tt.nilRef {
					require.Nil(t, v.Ref)
				} else {
					require.NotNil(t, v.Ref)
				}
				if tt.nilSpec {
					require.Nil(t, v.Spec)
				} else {
					require.NotNil(t, v.Spec)
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

func TestRefOrSpec_GetSpec(t *testing.T) {
	for _, tt := range []struct {
		name   string
		ref    any
		c      *spec.Extendable[spec.Components]
		exp    any
		expErr string
	}{
		{
			name: "with spec",
			ref:  spec.NewRefOrSpec[testRefOrSpec](&testRefOrSpec{A: "foo"}),
			exp:  &testRefOrSpec{A: "foo"},
		},
		{
			name:   "empty",
			ref:    spec.NewRefOrSpec[testRefOrSpec](nil),
			expErr: "not found",
		},
		{
			name:   "no components prefix if ref",
			ref:    spec.NewRefOrSpec[testRefOrSpec]("fooo"),
			expErr: "is not implemented",
		},
		{
			name:   "no components but with correct ref",
			ref:    spec.NewRefOrSpec[testRefOrSpec]("#/components/schemas/Pet"),
			expErr: "components is required",
		},
		{
			name: "correct ref and components",
			ref:  spec.NewRefOrSpec[spec.Schema]("#/components/schemas/Pet"),
			c: spec.NewExtendable((&spec.Components{}).
				WithRefOrSpec("Pet", spec.NewRefOrSpec[spec.Schema](&spec.Schema{Title: "foo"})),
			),
			exp: &spec.Schema{Title: "foo"},
		},
		{
			name: "ref to ref",
			ref:  spec.NewRefOrSpec[spec.Schema]("#/components/schemas/Pet"),
			c: spec.NewExtendable((&spec.Components{}).
				WithRefOrSpec("Pet", spec.NewRefOrSpec[spec.Schema]("#/components/schemas/Pet2")).
				WithRefOrSpec("Pet2", spec.NewRefOrSpec[spec.Schema](&spec.Schema{Title: "foo"})),
			),
			exp: &spec.Schema{Title: "foo"},
		},
		{
			name: "ref to incorrect ref",
			ref:  spec.NewRefOrSpec[spec.Schema]("#/components/schemas/Pet"),
			c: spec.NewExtendable((&spec.Components{}).
				WithRefOrSpec("Pet", spec.NewRefOrSpec[spec.Schema]("fooo")),
			),
			expErr: "is not implemented",
		},
		{
			name: "cycle ref",
			ref:  spec.NewRefOrSpec[spec.Schema]("#/components/schemas/Pet"),
			c: spec.NewExtendable((&spec.Components{}).
				WithRefOrSpec("Pet", spec.NewRefOrSpec[spec.Schema]("#/components/schemas/Pet2")).
				WithRefOrSpec("Pet2", spec.NewRefOrSpec[spec.Schema]("#/components/schemas/Pet")),
			),
			expErr: "cycle ref",
		},
		{
			name:   "ref to unexpected component",
			ref:    spec.NewRefOrSpec[testRefOrSpec]("#/components/test/Pet"),
			c:      spec.NewExtendable(&spec.Components{}),
			expErr: "unexpected component",
		},
		{
			name: "ref to unexpected component",
			ref:  spec.NewRefOrSpec[spec.Operation]("#/components/schemas/Pet"),
			c: spec.NewExtendable((&spec.Components{}).
				WithRefOrSpec("Pet", spec.NewRefOrSpec[spec.Schema](&spec.Schema{Title: "foo"})),
			),
			expErr: "expected spec of type",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			val := reflect.ValueOf(tt.ref).MethodByName("GetSpec").Call([]reflect.Value{reflect.ValueOf(tt.c)})
			require.Len(t, val, 2)
			if tt.expErr != "" {
				err, ok := val[1].Interface().(error)
				require.Truef(t, ok, "not error: %+v", val[1].Interface())
				require.ErrorContains(t, err, tt.expErr)
				require.Nil(t, val[0].Interface())
				return
			}
			require.Nil(t, val[1].Interface())
			require.Equal(t, tt.exp, val[0].Interface())
		})
	}
}

package openapi_test

import (
	"bytes"
	"encoding/json"
	"testing"

	goyaml "github.com/goccy/go-yaml"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/sv-tools/openapi"
)

type singleOrArrayCase[T any] struct {
	name     string
	data     []byte
	expected *openapi.SingleOrArray[T]
	wantErr  bool
}

func testSingleOrArray[T any](t *testing.T, tests []singleOrArrayCase[T]) {
	t.Helper()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run("json", func(t *testing.T) {
				t.Parallel()

				var o openapi.SingleOrArray[T]
				err := json.Unmarshal(tt.data, &o)
				if tt.wantErr {
					require.Error(t, err)
					return
				} else {
					require.NoError(t, err)
					require.Equal(t, *tt.expected, o)
				}
				newData, err := json.Marshal(&o)
				require.NoError(t, err)
				t.Log("orig: ", string(tt.data))
				t.Log(" new: ", string(newData))
				require.JSONEq(t, string(tt.data), string(newData))
			})
			t.Run("yaml.v3", func(t *testing.T) {
				t.Parallel()

				var o openapi.SingleOrArray[T]
				err := yaml.Unmarshal(tt.data, &o)
				if tt.wantErr {
					require.Error(t, err)
					return
				} else {
					require.NoError(t, err)
					require.Equal(t, *tt.expected, o)
				}
				newData, err := yaml.Marshal(&o)
				newData = bytes.TrimSpace(newData)
				require.NoError(t, err)
				t.Log("orig: ", string(tt.data))
				t.Log(" new: ", string(newData))
				require.YAMLEq(t, string(tt.data), string(newData))
			})
			t.Run("goccy/go-yaml", func(t *testing.T) {
				t.Parallel()

				var o openapi.SingleOrArray[T]
				err := goyaml.Unmarshal(tt.data, &o)
				if tt.wantErr {
					require.Error(t, err)
					return
				} else {
					require.NoError(t, err)
					require.Equal(t, *tt.expected, o)
				}
				newData, err := goyaml.Marshal(&o)
				newData = bytes.TrimSpace(newData)
				require.NoError(t, err)
				t.Log("orig: ", string(tt.data))
				t.Log(" new: ", string(newData))
				require.YAMLEq(t, string(tt.data), string(newData))
			})
		})
	}
}

func TestSingleOrArray_Marshal_Unmarshal(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		testSingleOrArray(t, []singleOrArrayCase[string]{
			{
				name:     "single",
				data:     []byte(`"single"`),
				expected: openapi.NewSingleOrArray("single"),
			},
			{
				name:     "multi",
				data:     []byte(`["first", "second"]`),
				expected: openapi.NewSingleOrArray("first", "second"),
			},
		})
	})

	t.Run("int", func(t *testing.T) {
		testSingleOrArray(t, []singleOrArrayCase[int]{
			{
				name:     "single",
				data:     []byte(`1`),
				expected: openapi.NewSingleOrArray(1),
			},
			{
				name:     "multi",
				data:     []byte(`[1, 2]`),
				expected: openapi.NewSingleOrArray(1, 2),
			},
			{
				name:    "string for int",
				data:    []byte(`"single"`),
				wantErr: true,
			},
			{
				name:    "array of string for int",
				data:    []byte(`["first", "second"]`),
				wantErr: true,
			},
		})
	})

	type Foo struct {
		A string `json:"a"`
		B int    `json:"b"`
	}
	t.Run("struct", func(t *testing.T) {
		testSingleOrArray(t, []singleOrArrayCase[Foo]{
			{
				name:     "single",
				data:     []byte(`{"a": "single", "b": 42}`),
				expected: openapi.NewSingleOrArray(Foo{A: "single", B: 42}),
			},
			{
				name:     "multi",
				data:     []byte(`[{"a": "first", "b": 1}, {"a": "second", "b": 2}]`),
				expected: openapi.NewSingleOrArray(Foo{A: "first", B: 1}, Foo{A: "second", B: 2}),
			},
			{
				name:    "string for struct",
				data:    []byte(`"single"`),
				wantErr: true,
			},
			{
				name:    "array of string for struct",
				data:    []byte(`["first", "second"]`),
				wantErr: true,
			},
		})
	})
}

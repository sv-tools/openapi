package spec_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/sv-tools/openapi/spec"
)

type singleOrArrayCase[T any] struct {
	name     string
	data     []byte
	expected spec.SingleOrArray[T]
	wantErr  bool
}

func testSingleOrArrayJSON[T any](t *testing.T, tests []singleOrArrayCase[T]) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var o spec.SingleOrArray[T]
			err := json.Unmarshal(tt.data, &o)
			if tt.wantErr {
				require.Error(t, err)
				return
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, o)
			}
			newData, err := json.Marshal(&o)
			require.NoError(t, err)
			t.Log("orig: ", string(tt.data))
			t.Log(" new: ", string(newData))
			require.JSONEq(t, string(tt.data), string(newData))
		})
	}
}

func TestSingleOrArrayJSON(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		testSingleOrArrayJSON(t, []singleOrArrayCase[string]{
			{
				name:     "single",
				data:     []byte(`"single"`),
				expected: spec.NewSingleOrArray("single"),
			},
			{
				name:     "multi",
				data:     []byte(`["first","second"]`),
				expected: spec.NewSingleOrArray("first", "second"),
			},
			{
				name: "null",
				data: []byte(`null`),
			},
			{
				name:    "int for string",
				data:    []byte(`42`),
				wantErr: true,
			},
			{
				name:    "array of int for string",
				data:    []byte(`[42, 103]`),
				wantErr: true,
			},
			{
				name:    "empty for string",
				data:    []byte(``),
				wantErr: true,
			},
		})
	})

	t.Run("int", func(t *testing.T) {
		testSingleOrArrayJSON(t, []singleOrArrayCase[int]{
			{
				name:     "single",
				data:     []byte(`1`),
				expected: spec.NewSingleOrArray(1),
			},
			{
				name:     "multi",
				data:     []byte(`[1,2]`),
				expected: spec.NewSingleOrArray(1, 2),
			},
			{
				name: "null",
				data: []byte(`null`),
			},
			{
				name:    "string for int",
				data:    []byte(`"single"`),
				wantErr: true,
			},
			{
				name:    "array of string for int",
				data:    []byte(`["first","second"]`),
				wantErr: true,
			},
			{
				name:    "empty for int",
				data:    []byte(``),
				wantErr: true,
			},
		})
	})

	type Foo struct {
		A string
		B int
	}
	t.Run("struct", func(t *testing.T) {
		testSingleOrArrayJSON(t, []singleOrArrayCase[Foo]{
			{
				name:     "single",
				data:     []byte(`{"A":"single","B":42}`),
				expected: spec.NewSingleOrArray(Foo{A: "single", B: 42}),
			},
			{
				name:     "multi",
				data:     []byte(`[{"A":"first","B":1},{"A":"second","B":2}]`),
				expected: spec.NewSingleOrArray(Foo{A: "first", B: 1}, Foo{A: "second", B: 2}),
			},
			{
				name: "null",
				data: []byte(`null`),
			},
			{
				name:    "string for struct",
				data:    []byte(`"single"`),
				wantErr: true,
			},
			{
				name:    "array of string for struct",
				data:    []byte(`["first","second"]`),
				wantErr: true,
			},
			{
				name:    "empty for struct",
				data:    []byte(``),
				wantErr: true,
			},
		})
	})
}

func testSingleOrArrayYAML[T any](t *testing.T, tests []singleOrArrayCase[T]) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//t.Parallel()

			var o spec.SingleOrArray[T]
			err := yaml.Unmarshal(tt.data, &o)
			if tt.wantErr {
				require.Error(t, err)
				return
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, o)
			}
			newData, err := yaml.Marshal(&o)
			newData = bytes.TrimSpace(newData)
			require.NoError(t, err)
			t.Log("orig: ", string(tt.data))
			t.Log(" new: ", string(newData))
			require.Equal(t, tt.data, newData)
		})
	}
}

func TestSingleOrArrayYAML(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		testSingleOrArrayYAML(t, []singleOrArrayCase[string]{
			{
				name:     "single",
				data:     []byte(`single`),
				expected: spec.NewSingleOrArray("single"),
			},
			{
				name: "multi",
				data: []byte(`- first
- second`),
				expected: spec.NewSingleOrArray("first", "second"),
			},
		})
	})

	t.Run("int", func(t *testing.T) {
		testSingleOrArrayYAML(t, []singleOrArrayCase[int]{
			{
				name:     "single",
				data:     []byte(`1`),
				expected: spec.NewSingleOrArray(1),
			},
			{
				name: "multi",
				data: []byte(`- 1
- 2`),
				expected: spec.NewSingleOrArray(1, 2),
			},
			{
				name:    "string for int",
				data:    []byte(`single`),
				wantErr: true,
			},
			{
				name: "array of string for int",
				data: []byte(`- first
- second`),
				wantErr: true,
			},
		})
	})

	type Foo struct {
		A string
		B int
	}
	t.Run("struct", func(t *testing.T) {
		testSingleOrArrayYAML(t, []singleOrArrayCase[Foo]{
			{
				name: "single",
				data: []byte(`a: single
b: 42`),
				expected: spec.NewSingleOrArray(Foo{A: "single", B: 42}),
			},
			{
				name: "multi",
				data: []byte(`- a: first
  b: 1
- a: second
  b: 2`),
				expected: spec.NewSingleOrArray(Foo{A: "first", B: 1}, Foo{A: "second", B: 2}),
			},
			{
				name:    "string for struct",
				data:    []byte(`single`),
				wantErr: true,
			},
			{
				name: "array of string for struct",
				data: []byte(`- first
- second`),
				wantErr: true,
			},
		})
	})
}

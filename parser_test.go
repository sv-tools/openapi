package openapi_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sv-tools/openapi"
)

type Simple struct {
	Fs        string `json:"fs,omitempty" openapi:",format:password" yaml:"FS,omitempty"` // json name should be used
	Fi        int    `yaml:"FI,omitempty"`                                                // yaml name shuld be used
	Fb        *bool
	Fbs       []byte            `json:"fbs,omitempty"                                  openapi:"fBS" yaml:"FS,omitempty"` // openapi name should be used
	Fm        map[string]string `openapi:",required,title:Map of strings,addtype:null"`                                   // default field name should be used
	Excluded1 map[string]string `openapi:"-"`
	Excluded2 map[string]string `json:"-"`
	Excluded3 map[string]string `yaml:"-"`
	Fa        any               `openapi:",deprecated"`

	fp string
}

type Complex struct {
	Simple          // anonymous field
	Next   *Complex `json:"Next"` // circular references
}

type SimpleByRef struct {
	S Simple `json:"s" openapi:"s,required,title:Simple By Ref,ref:#/components/schemas/github.com.sv-tools.openapi_test.Simple" yaml:"s"`
}

func TestParseObject(t *testing.T) {
	trueVar := true
	strVar := "foo"
	var nilBool *bool

	for _, tt := range []struct {
		name               string
		obj                any
		expected           *openapi.RefOrSpec[openapi.Schema]
		expectedComponents *openapi.Components
		err                string
	}{
		{
			name: "nil",
			obj:  nil,
			expected: openapi.NewSchemaBuilder().
				Type(openapi.NullType).
				GoType("nil").Build(),
		},
		{
			name: "bool true",
			obj:  true,
			expected: openapi.NewSchemaBuilder().
				Type(openapi.BooleanType).
				GoType("bool").Build(),
		},
		{
			name: "bool false",
			obj:  false,
			expected: openapi.NewSchemaBuilder().
				Type(openapi.BooleanType).
				GoType("bool").Build(),
		},
		{
			name: "ptr to bool true",
			obj:  &trueVar,
			expected: openapi.NewSchemaBuilder().
				Type(openapi.BooleanType, openapi.NullType).
				GoType("bool").Build(),
		},
		{
			name: "ptr to bool false",
			obj:  &trueVar,
			expected: openapi.NewSchemaBuilder().
				Type(openapi.BooleanType, openapi.NullType).
				GoType("bool").Build(),
		},
		{
			name: "ptr to *bool nil",
			obj:  nilBool,
			expected: openapi.NewSchemaBuilder().
				Type(openapi.BooleanType, openapi.NullType).
				GoType("bool").Build(),
		},
		{
			name: "int",
			obj:  42,
			expected: openapi.NewSchemaBuilder().
				Type(openapi.IntegerType).
				Format(openapi.Int64Format).
				GoType("int").Build(),
		},
		{
			name: "int8",
			obj:  int8(42),
			expected: openapi.NewSchemaBuilder().
				Type(openapi.IntegerType).
				Format(openapi.Int32Format).
				GoType("int8").Build(),
		},
		{
			name: "int32",
			obj:  int32(42),
			expected: openapi.NewSchemaBuilder().
				Type(openapi.IntegerType).
				Format(openapi.Int32Format).
				GoType("int32").Build(),
		},
		{
			name: "int64",
			obj:  int64(42),
			expected: openapi.NewSchemaBuilder().
				Type(openapi.IntegerType).
				Format(openapi.Int64Format).
				GoType("int64").Build(),
		},
		{
			name: "uint",
			obj:  uint(42),
			expected: openapi.NewSchemaBuilder().
				Type(openapi.IntegerType).
				Format(openapi.Int64Format).
				GoType("uint").Build(),
		},
		{
			name: "uint8",
			obj:  uint8(42),
			expected: openapi.NewSchemaBuilder().
				Type(openapi.IntegerType).
				Format(openapi.Int32Format).
				GoType("uint8").Build(),
		},
		{
			name: "uint32",
			obj:  uint32(42),
			expected: openapi.NewSchemaBuilder().
				Type(openapi.IntegerType).
				Format(openapi.Int32Format).
				GoType("uint32").Build(),
		},
		{
			name: "uint64",
			obj:  uint64(42),
			expected: openapi.NewSchemaBuilder().
				Type(openapi.IntegerType).
				Format(openapi.Int64Format).
				GoType("uint64").Build(),
		},
		{
			name: "float32",
			obj:  float32(42),
			expected: openapi.NewSchemaBuilder().
				Type(openapi.NumberType).
				Format(openapi.FloatFormat).
				GoType("float32").Build(),
		},
		{
			name: "float64",
			obj:  float64(42),
			expected: openapi.NewSchemaBuilder().
				Type(openapi.NumberType).
				Format(openapi.DoubleFormat).
				GoType("float64").Build(),
		},
		{
			name: "string",
			obj:  "foo",
			expected: openapi.NewSchemaBuilder().
				Type(openapi.StringType).
				GoType("string").Build(),
		},
		{
			name: "bytes",
			obj:  []byte("foo"),
			expected: openapi.NewSchemaBuilder().
				Type(openapi.StringType).
				ContentEncoding(openapi.Base64Encoding).
				GoType("[]byte").Build(),
		},
		{
			name: "map string",
			obj:  map[string]string{"foo": "bar"},
			expected: openapi.NewSchemaBuilder().Type(openapi.ObjectType).
				AdditionalProperties(openapi.NewBoolOrSchema(openapi.NewSchemaBuilder().
					Type(openapi.StringType).
					GoType("string").Build(),
				)).Build(),
		},
		{
			name: "map string ref string",
			obj:  map[string]*string{"foo": &strVar},
			expected: openapi.NewSchemaBuilder().Type(openapi.ObjectType).
				AdditionalProperties(openapi.NewBoolOrSchema(openapi.NewSchemaBuilder().
					Type(openapi.StringType, openapi.NullType).
					GoType("string").Build(),
				)).Build(),
		},
		{
			name: "map string int",
			obj:  map[string]int{"foo": 42},
			expected: openapi.NewSchemaBuilder().
				Type(openapi.ObjectType).
				AdditionalProperties(openapi.NewBoolOrSchema(openapi.NewSchemaBuilder().
					Type(openapi.IntegerType).
					Format(openapi.Int64Format).
					GoType("int").Build(),
				)).Build(),
		},
		{
			name: "map string any",
			obj:  map[string]any{"foo": 42, "bar": "baz"},
			expected: openapi.NewSchemaBuilder().
				Type(openapi.ObjectType).
				AdditionalProperties(openapi.NewBoolOrSchema(true)).
				Build(),
		},
		{
			name: "slice int",
			obj:  []int{42},
			expected: openapi.NewSchemaBuilder().
				Type(openapi.ArrayType).
				Items(openapi.NewBoolOrSchema(openapi.NewSchemaBuilder().
					Type(openapi.IntegerType).
					Format(openapi.Int64Format).
					GoType("int").Build(),
				)).Build(),
		},
		{
			name: "slice string",
			obj:  []string{"foo"},
			expected: openapi.NewSchemaBuilder().
				Type(openapi.ArrayType).
				Items(openapi.NewBoolOrSchema(openapi.NewSchemaBuilder().
					Type(openapi.StringType).
					GoType("string").Build(),
				)).Build(),
		},
		{
			name: "slice any",
			obj:  []any{"foo", 42},
			expected: openapi.NewSchemaBuilder().Type(openapi.ArrayType).
				Items(openapi.NewBoolOrSchema(true)).
				Build(),
		},
		{
			name: "double slice any",
			obj:  [][]any{{"foo", 42}},
			expected: openapi.NewSchemaBuilder().
				Type(openapi.ArrayType).
				Items(openapi.NewBoolOrSchema(openapi.NewSchemaBuilder().
					Type(openapi.ArrayType).
					Items(openapi.NewBoolOrSchema(true)).
					Build(),
				)).Build(),
		},
		{
			name: "triple slice any",
			obj:  [][][]any{{{"foo", 42}}},
			expected: openapi.NewSchemaBuilder().
				Type(openapi.ArrayType).
				Items(openapi.NewBoolOrSchema(openapi.NewSchemaBuilder().
					Type(openapi.ArrayType).
					Items(openapi.NewBoolOrSchema(
						openapi.NewSchemaBuilder().
							Type(openapi.ArrayType).
							Items(openapi.NewBoolOrSchema(true)).
							Build(),
					)).Build(),
				)).Build(),
		},
		{
			name: "map string any",
			obj:  map[string]map[string]any{"xyz": {"foo": 42, "bar": "baz"}},
			expected: openapi.NewSchemaBuilder().
				Type(openapi.ObjectType).
				AdditionalProperties(openapi.NewBoolOrSchema(openapi.NewSchemaBuilder().
					Type(openapi.ObjectType).
					AdditionalProperties(openapi.NewBoolOrSchema(true)).
					Build(),
				)).Build(),
		},
		{
			name: "slice map string any",
			obj:  []map[string]any{{"foo": 42, "bar": "baz"}},
			expected: openapi.NewSchemaBuilder().
				Type(openapi.ArrayType).
				Items(openapi.NewBoolOrSchema(openapi.NewSchemaBuilder().
					Type(openapi.ObjectType).
					AdditionalProperties(openapi.NewBoolOrSchema(true)).
					Build(),
				)).Build(),
		},
		{
			name: "json number",
			obj:  json.Number("42"),
			expected: openapi.NewSchemaBuilder().
				Type(openapi.NumberType).
				GoPackage("encoding/json").GoType("json.Number").Build(),
		},
		{
			name: "json raw",
			obj:  json.RawMessage(`"foo"`),
			expected: openapi.NewSchemaBuilder().
				Type(openapi.StringType).
				ContentMediaType("application/json").
				GoPackage("encoding/json").GoType("json.RawMessage").Build(),
		},
		{
			name: "simple struct",
			obj: Simple{
				Fs:  "foo",
				Fi:  42,
				Fb:  &trueVar,
				Fbs: []byte("bar"),
				Fm:  map[string]string{"baz": "qux"},
				Fa:  []any{"435", 42, false},
				fp:  "baz",
			},
			expected: openapi.NewSchemaBuilder().Ref("#/components/schemas/github.com.sv-tools.openapi_test.Simple").Build(),
			expectedComponents: openapi.NewComponents().Spec.Add(
				"github.com.sv-tools.openapi_test.Simple",
				openapi.NewSchemaBuilder().
					Type(openapi.ObjectType).
					AddProperty("fs", openapi.NewSchemaBuilder().Type(openapi.StringType).GoType("string").Format("password").Build()).
					AddProperty("FI", openapi.NewSchemaBuilder().Type(openapi.IntegerType).Format(openapi.Int64Format).GoType("int").Build()).
					AddProperty("Fb", openapi.NewSchemaBuilder().Type(openapi.BooleanType, openapi.NullType).GoType("bool").Build()).
					AddProperty("fBS", openapi.NewSchemaBuilder().Type(openapi.StringType).ContentEncoding(openapi.Base64Encoding).GoType("[]byte").Build()).
					AddProperty("Fm", openapi.NewSchemaBuilder().
						Type(openapi.ObjectType, openapi.NullType).
						Title("Map of strings").
						AdditionalProperties(openapi.NewBoolOrSchema(openapi.NewSchemaBuilder().
							Type(openapi.StringType).
							GoType("string").Build(),
						)).Build(),
					).
					AddProperty("Fa", openapi.NewSchemaBuilder().
						Deprecated(true).
						GoType("any").Build(),
					).
					AddRequired("Fm").
					GoPackage("github.com/sv-tools/openapi_test").GoType("openapi_test.Simple").Build(),
			),
		},
		{
			name: "complex struct",
			obj: Complex{
				Simple: Simple{
					Fs:  "foo",
					Fi:  42,
					Fb:  &trueVar,
					Fbs: []byte("bar"),
					Fm:  map[string]string{"baz": "qux"},
					Fa:  []any{"435", 42, false},
					fp:  "baz",
				},
				Next: &Complex{},
			},
			expected: openapi.NewSchemaBuilder().Ref("#/components/schemas/github.com.sv-tools.openapi_test.Complex").Build(),
			expectedComponents: openapi.NewComponents().Spec.Add(
				"github.com.sv-tools.openapi_test.Complex",
				openapi.NewSchemaBuilder().
					AllOf(
						openapi.NewSchemaBuilder().Ref("#/components/schemas/github.com.sv-tools.openapi_test.Simple").Build(),
						openapi.NewSchemaBuilder().
							Type(openapi.ObjectType).
							AddProperty("Next", openapi.NewSchemaBuilder().
								OneOf(
									openapi.NewSchemaBuilder().Ref("#/components/schemas/github.com.sv-tools.openapi_test.Complex").Build(),
									openapi.NewSchemaBuilder().Type(openapi.NullType).Build(),
								).
								Build(),
							).
							Build(),
					).
					GoPackage("github.com/sv-tools/openapi_test").GoType("openapi_test.Complex").Build(),
			).Add(
				"github.com.sv-tools.openapi_test.Simple",
				openapi.NewSchemaBuilder().
					Type(openapi.ObjectType).
					AddProperty("fs", openapi.NewSchemaBuilder().Type(openapi.StringType).GoType("string").Format("password").Build()).
					AddProperty("FI", openapi.NewSchemaBuilder().Type(openapi.IntegerType).Format(openapi.Int64Format).GoType("int").Build()).
					AddProperty("Fb", openapi.NewSchemaBuilder().Type(openapi.BooleanType, openapi.NullType).GoType("bool").Build()).
					AddProperty("fBS", openapi.NewSchemaBuilder().Type(openapi.StringType).ContentEncoding(openapi.Base64Encoding).GoType("[]byte").Build()).
					AddProperty("Fm", openapi.NewSchemaBuilder().
						Type(openapi.ObjectType, openapi.NullType).
						Title("Map of strings").
						AdditionalProperties(openapi.NewBoolOrSchema(openapi.NewSchemaBuilder().
							Type(openapi.StringType).
							GoType("string").Build(),
						)).Build(),
					).
					AddProperty("Fa", openapi.NewSchemaBuilder().
						Deprecated(true).
						GoType("any").Build(),
					).
					AddRequired("Fm").
					GoPackage("github.com/sv-tools/openapi_test").GoType("openapi_test.Simple").Build(),
			),
		},
		{
			name: "simple by ref",
			obj: SimpleByRef{
				S: Simple{},
			},
			expected: openapi.NewSchemaBuilder().Ref("#/components/schemas/github.com.sv-tools.openapi_test.SimpleByRef").Build(),
			expectedComponents: openapi.NewComponents().Spec.Add(
				"github.com.sv-tools.openapi_test.SimpleByRef",
				openapi.NewSchemaBuilder().
					Type(openapi.ObjectType).
					AddProperty("s", openapi.NewSchemaBuilder().Ref("#/components/schemas/github.com.sv-tools.openapi_test.Simple").Title("Simple By Ref").Build()).
					Required("s").
					GoPackage("github.com/sv-tools/openapi_test").GoType("openapi_test.SimpleByRef").Build(),
			).Add(
				"github.com.sv-tools.openapi_test.Simple",
				openapi.NewSchemaBuilder().
					Type(openapi.ObjectType).
					AddProperty("fs", openapi.NewSchemaBuilder().Type(openapi.StringType).GoType("string").Format("password").Build()).
					AddProperty("FI", openapi.NewSchemaBuilder().Type(openapi.IntegerType).Format(openapi.Int64Format).GoType("int").Build()).
					AddProperty("Fb", openapi.NewSchemaBuilder().Type(openapi.BooleanType, openapi.NullType).GoType("bool").Build()).
					AddProperty("fBS", openapi.NewSchemaBuilder().Type(openapi.StringType).ContentEncoding(openapi.Base64Encoding).GoType("[]byte").Build()).
					AddProperty("Fm", openapi.NewSchemaBuilder().
						Type(openapi.ObjectType, openapi.NullType).
						Title("Map of strings").
						AdditionalProperties(openapi.NewBoolOrSchema(openapi.NewSchemaBuilder().
							Type(openapi.StringType).
							GoType("string").Build(),
						)).Build(),
					).
					AddProperty("Fa", openapi.NewSchemaBuilder().
						Deprecated(true).
						GoType("any").Build(),
					).
					AddRequired("Fm").
					GoPackage("github.com/sv-tools/openapi_test").GoType("openapi_test.Simple").Build(),
			),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			spec := openapi.NewOpenAPIBuilder().
				Info(openapi.NewInfoBuilder().Title("Test").Version("1.0").Build()).
				Components(openapi.NewComponents()).
				Build()
			schema, err := openapi.ParseObject(tt.obj, spec.Spec.Components)
			if tt.err != "" {
				require.ErrorContains(t, err, tt.err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, schema)

			actual, err := schema.Build().MarshalJSON()
			require.NoError(t, err)

			expected, err := tt.expected.MarshalJSON()
			require.NoError(t, err)

			require.JSONEq(t, string(expected), string(actual))

			if tt.expectedComponents != nil {
				actualComponents, err := spec.Spec.Components.MarshalJSON()
				require.NoError(t, err)

				expectedComponents, err := json.Marshal(tt.expectedComponents)
				require.NoError(t, err)

				require.JSONEq(t, string(expectedComponents), string(actualComponents))
			}

			spec.Spec.Components.Spec.Add("test", schema.Build())
			validator, err := openapi.NewValidator(
				spec,
				openapi.AllowUnusedComponents(),
			)
			require.NoError(t, err)

			require.NoError(t, validator.ValidateSpec())

			value, err := openapi.ConvertToJSON(tt.obj)
			require.NoError(t, err)

			pretty, _ := json.MarshalIndent(tt.obj, "", "  ")
			t.Logf("obj: %s", pretty)

			require.NoError(t, validator.ValidateData("#/components/schemas/test", value))
		})
	}
}

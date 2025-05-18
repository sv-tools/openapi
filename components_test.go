package openapi_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sv-tools/openapi"
)

func TestComponents_Add(t *testing.T) {
	for _, tt := range []struct {
		name   string
		create func(tb testing.TB) (string, any)
		check  func(tb testing.TB, c *openapi.Components)
	}{
		{
			name: "schema ref or spec",
			create: func(tb testing.TB) (string, any) {
				tb.Helper()

				o := openapi.NewSchemaBuilder().Title("test").Build()
				return "testSchema", o
			},
			check: func(tb testing.TB, c *openapi.Components) {
				tb.Helper()

				require.Len(tb, c.Schemas, 1)
				require.NotNil(tb, c.Schemas["testSchema"])
				require.NotNil(tb, c.Schemas["testSchema"].Spec)
				require.Equal(tb, "test", c.Schemas["testSchema"].Spec.Title)
			},
		},
		{
			name: "response spec",
			create: func(tb testing.TB) (string, any) {
				tb.Helper()

				o := openapi.NewResponseBuilder().Description("test").Build()
				return "testResponse", o
			},
			check: func(tb testing.TB, c *openapi.Components) {
				tb.Helper()

				require.Len(tb, c.Responses, 1)
				require.NotNil(tb, c.Responses["testResponse"])
				require.NotNil(tb, c.Responses["testResponse"].Spec)
				require.NotNil(tb, c.Responses["testResponse"].Spec.Spec)
				require.Equal(tb, "test", c.Responses["testResponse"].Spec.Spec.Description)
			},
		},
		{
			name: "parameter spec",
			create: func(tb testing.TB) (string, any) {
				tb.Helper()

				o := openapi.NewParameterBuilder().Description("test").Build()
				return "testParameter", o
			},
			check: func(tb testing.TB, c *openapi.Components) {
				tb.Helper()

				require.Len(tb, c.Parameters, 1)
				require.NotNil(tb, c.Parameters["testParameter"])
				require.NotNil(tb, c.Parameters["testParameter"].Spec)
				require.NotNil(tb, c.Parameters["testParameter"].Spec.Spec)
				require.Equal(tb, "test", c.Parameters["testParameter"].Spec.Spec.Description)
			},
		},
		{
			name: "examples spec",
			create: func(tb testing.TB) (string, any) {
				tb.Helper()

				o := openapi.NewExampleBuilder().Description("test").Build()
				return "testExamples", o
			},
			check: func(tb testing.TB, c *openapi.Components) {
				tb.Helper()

				require.Len(tb, c.Examples, 1)
				require.NotNil(tb, c.Examples["testExamples"])
				require.NotNil(tb, c.Examples["testExamples"].Spec)
				require.NotNil(tb, c.Examples["testExamples"].Spec.Spec)
				require.Equal(tb, "test", c.Examples["testExamples"].Spec.Spec.Description)
			},
		},
		{
			name: "request body spec",
			create: func(tb testing.TB) (string, any) {
				tb.Helper()

				o := openapi.NewRequestBodyBuilder().Description("test").Build()
				return "testRequestBodies", o
			},
			check: func(tb testing.TB, c *openapi.Components) {
				tb.Helper()

				require.Len(tb, c.RequestBodies, 1)
				require.NotNil(tb, c.RequestBodies["testRequestBodies"])
				require.NotNil(tb, c.RequestBodies["testRequestBodies"].Spec)
				require.NotNil(tb, c.RequestBodies["testRequestBodies"].Spec.Spec)
				require.Equal(tb, "test", c.RequestBodies["testRequestBodies"].Spec.Spec.Description)
			},
		},
		{
			name: "headers spec",
			create: func(tb testing.TB) (string, any) {
				tb.Helper()

				o := openapi.NewHeaderBuilder().Description("test").Build()
				return "testHeader", o
			},
			check: func(tb testing.TB, c *openapi.Components) {
				tb.Helper()

				require.Len(tb, c.Headers, 1)
				require.NotNil(tb, c.Headers["testHeader"])
				require.NotNil(tb, c.Headers["testHeader"].Spec)
				require.NotNil(tb, c.Headers["testHeader"].Spec.Spec)
				require.Equal(tb, "test", c.Headers["testHeader"].Spec.Spec.Description)
			},
		},
		{
			name: "security schemes spec",
			create: func(tb testing.TB) (string, any) {
				tb.Helper()

				o := &openapi.SecurityScheme{
					Description: "test",
				}
				return "testSecurityScheme", openapi.NewRefOrExtSpec[openapi.SecurityScheme](o)
			},
			check: func(tb testing.TB, c *openapi.Components) {
				tb.Helper()

				require.Len(tb, c.SecuritySchemes, 1)
				require.NotNil(tb, c.SecuritySchemes["testSecurityScheme"])
				require.NotNil(tb, c.SecuritySchemes["testSecurityScheme"].Spec)
				require.NotNil(tb, c.SecuritySchemes["testSecurityScheme"].Spec.Spec)
				require.Equal(tb, "test", c.SecuritySchemes["testSecurityScheme"].Spec.Spec.Description)
			},
		},
		{
			name: "link spec",
			create: func(tb testing.TB) (string, any) {
				tb.Helper()

				o := openapi.NewLinkBuilder().Description("test").Build()
				return "testLink", o
			},
			check: func(tb testing.TB, c *openapi.Components) {
				tb.Helper()

				require.Len(tb, c.Links, 1)
				require.NotNil(tb, c.Links["testLink"])
				require.NotNil(tb, c.Links["testLink"].Spec)
				require.NotNil(tb, c.Links["testLink"].Spec.Spec)
				require.Equal(tb, "test", c.Links["testLink"].Spec.Spec.Description)
			},
		},
		{
			name: "callback spec",
			create: func(tb testing.TB) (string, any) {
				tb.Helper()

				o := openapi.NewCallbackBuilder().AddPathItem(
					"testPath",
					openapi.NewPathItemBuilder().Description("test").Build(),
				).Build()
				return "testCallback", o
			},
			check: func(tb testing.TB, c *openapi.Components) {
				tb.Helper()

				require.Len(tb, c.Callbacks, 1)
				require.NotNil(tb, c.Callbacks["testCallback"])
				require.NotNil(tb, c.Callbacks["testCallback"].Spec)
				require.NotNil(tb, c.Callbacks["testCallback"].Spec.Spec)
				paths := c.Callbacks["testCallback"].Spec.Spec.Paths
				require.Len(tb, paths, 1)
				require.NotNil(tb, paths["testPath"])
				require.NotNil(tb, paths["testPath"].Spec)
				require.NotNil(tb, paths["testPath"].Spec.Spec)
				require.Equal(tb, "test", paths["testPath"].Spec.Spec.Description)
			},
		},
		{
			name: "path item spec",
			create: func(tb testing.TB) (string, any) {
				tb.Helper()

				o := openapi.NewPathItemBuilder().Description("test").Build()
				return "testPathItem", o
			},
			check: func(tb testing.TB, c *openapi.Components) {
				tb.Helper()

				require.Len(tb, c.Paths, 1)
				require.NotNil(tb, c.Paths["testPathItem"])
				require.NotNil(tb, c.Paths["testPathItem"].Spec)
				require.NotNil(tb, c.Paths["testPathItem"].Spec.Spec)
				require.Equal(tb, "test", c.Paths["testPathItem"].Spec.Spec.Description)
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			name, obj := tt.create(t)
			tt.check(t, (&openapi.Components{}).Add(name, obj))
		})
	}
}

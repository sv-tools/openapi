package spec_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sv-tools/openapi/spec"
)

func TestComponents_WithRefOrSpec(t *testing.T) {
	for _, tt := range []struct {
		name   string
		create func() (string, any)
		check  func(tb testing.TB, c *spec.Components)
	}{
		{
			name: "schema spec",
			create: func() (string, any) {
				o := spec.NewSchemaSpec()
				o.Spec.Title = "test"
				return "testSchema", o.Spec
			},
			check: func(tb testing.TB, c *spec.Components) {
				require.Len(tb, c.Schemas, 1)
				require.NotNil(tb, c.Schemas["testSchema"])
				require.NotNil(tb, c.Schemas["testSchema"].Spec)
				require.Equal(tb, "test", c.Schemas["testSchema"].Spec.Title)
			},
		},
		{
			name: "schema ref or spec",
			create: func() (string, any) {
				o := spec.NewSchemaSpec()
				o.Spec.Title = "test"
				return "testSchema", o
			},
			check: func(tb testing.TB, c *spec.Components) {
				require.Len(tb, c.Schemas, 1)
				require.NotNil(tb, c.Schemas["testSchema"])
				require.NotNil(tb, c.Schemas["testSchema"].Spec)
				require.Equal(tb, "test", c.Schemas["testSchema"].Spec.Title)
			},
		},
		{
			name: "response spec",
			create: func() (string, any) {
				o := spec.NewResponseSpec()
				o.Spec.Spec.Description = "test"
				return "testResponse", o.Spec.Spec
			},
			check: func(tb testing.TB, c *spec.Components) {
				require.Len(tb, c.Responses, 1)
				require.NotNil(tb, c.Responses["testResponse"])
				require.NotNil(tb, c.Responses["testResponse"].Spec)
				require.NotNil(tb, c.Responses["testResponse"].Spec.Spec)
				require.Equal(tb, "test", c.Responses["testResponse"].Spec.Spec.Description)
			},
		},
		{
			name: "response ext spec",
			create: func() (string, any) {
				o := spec.NewResponseSpec()
				o.Spec.Spec.Description = "test"
				return "testResponse", o.Spec
			},
			check: func(tb testing.TB, c *spec.Components) {
				require.Len(tb, c.Responses, 1)
				require.NotNil(tb, c.Responses["testResponse"])
				require.NotNil(tb, c.Responses["testResponse"].Spec)
				require.NotNil(tb, c.Responses["testResponse"].Spec.Spec)
				require.Equal(tb, "test", c.Responses["testResponse"].Spec.Spec.Description)
			},
		},
		{
			name: "response ref or spec",
			create: func() (string, any) {
				o := spec.NewResponseSpec()
				o.Spec.Spec.Description = "test"
				return "testResponse", o
			},
			check: func(tb testing.TB, c *spec.Components) {
				require.Len(tb, c.Responses, 1)
				require.NotNil(tb, c.Responses["testResponse"])
				require.NotNil(tb, c.Responses["testResponse"].Spec)
				require.NotNil(tb, c.Responses["testResponse"].Spec.Spec)
				require.Equal(tb, "test", c.Responses["testResponse"].Spec.Spec.Description)
			},
		},
		{
			name: "parameter spec",
			create: func() (string, any) {
				o := spec.NewParameterSpec()
				o.Spec.Spec.Description = "test"
				return "testParameter", o.Spec.Spec
			},
			check: func(tb testing.TB, c *spec.Components) {
				require.Len(tb, c.Parameters, 1)
				require.NotNil(tb, c.Parameters["testParameter"])
				require.NotNil(tb, c.Parameters["testParameter"].Spec)
				require.NotNil(tb, c.Parameters["testParameter"].Spec.Spec)
				require.Equal(tb, "test", c.Parameters["testParameter"].Spec.Spec.Description)
			},
		},
		{
			name: "parameter ext spec",
			create: func() (string, any) {
				o := spec.NewParameterSpec()
				o.Spec.Spec.Description = "test"
				return "testParameter", o.Spec
			},
			check: func(tb testing.TB, c *spec.Components) {
				require.Len(tb, c.Parameters, 1)
				require.NotNil(tb, c.Parameters["testParameter"])
				require.NotNil(tb, c.Parameters["testParameter"].Spec)
				require.NotNil(tb, c.Parameters["testParameter"].Spec.Spec)
				require.Equal(tb, "test", c.Parameters["testParameter"].Spec.Spec.Description)
			},
		},
		{
			name: "parameter ref or spec",
			create: func() (string, any) {
				o := spec.NewParameterSpec()
				o.Spec.Spec.Description = "test"
				return "testParameter", o
			},
			check: func(tb testing.TB, c *spec.Components) {
				require.Len(tb, c.Parameters, 1)
				require.NotNil(tb, c.Parameters["testParameter"])
				require.NotNil(tb, c.Parameters["testParameter"].Spec)
				require.NotNil(tb, c.Parameters["testParameter"].Spec.Spec)
				require.Equal(tb, "test", c.Parameters["testParameter"].Spec.Spec.Description)
			},
		},
		{
			name: "examples spec",
			create: func() (string, any) {
				o := spec.NewExampleSpec()
				o.Spec.Spec.Description = "test"
				return "testExamples", o.Spec.Spec
			},
			check: func(tb testing.TB, c *spec.Components) {
				require.Len(tb, c.Examples, 1)
				require.NotNil(tb, c.Examples["testExamples"])
				require.NotNil(tb, c.Examples["testExamples"].Spec)
				require.NotNil(tb, c.Examples["testExamples"].Spec.Spec)
				require.Equal(tb, "test", c.Examples["testExamples"].Spec.Spec.Description)
			},
		},
		{
			name: "examples ext spec",
			create: func() (string, any) {
				o := spec.NewExampleSpec()
				o.Spec.Spec.Description = "test"
				return "testExamples", o.Spec
			},
			check: func(tb testing.TB, c *spec.Components) {
				require.Len(tb, c.Examples, 1)
				require.NotNil(tb, c.Examples["testExamples"])
				require.NotNil(tb, c.Examples["testExamples"].Spec)
				require.NotNil(tb, c.Examples["testExamples"].Spec.Spec)
				require.Equal(tb, "test", c.Examples["testExamples"].Spec.Spec.Description)
			},
		},
		{
			name: "examples ref or spec",
			create: func() (string, any) {
				o := spec.NewExampleSpec()
				o.Spec.Spec.Description = "test"
				return "testExamples", o
			},
			check: func(tb testing.TB, c *spec.Components) {
				require.Len(tb, c.Examples, 1)
				require.NotNil(tb, c.Examples["testExamples"])
				require.NotNil(tb, c.Examples["testExamples"].Spec)
				require.NotNil(tb, c.Examples["testExamples"].Spec.Spec)
				require.Equal(tb, "test", c.Examples["testExamples"].Spec.Spec.Description)
			},
		},
		{
			name: "request bodies spec",
			create: func() (string, any) {
				o := spec.NewRequestBodySpec()
				o.Spec.Spec.Description = "test"
				return "testRequestBodies", o.Spec.Spec
			},
			check: func(tb testing.TB, c *spec.Components) {
				require.Len(tb, c.RequestBodies, 1)
				require.NotNil(tb, c.RequestBodies["testRequestBodies"])
				require.NotNil(tb, c.RequestBodies["testRequestBodies"].Spec)
				require.NotNil(tb, c.RequestBodies["testRequestBodies"].Spec.Spec)
				require.Equal(tb, "test", c.RequestBodies["testRequestBodies"].Spec.Spec.Description)
			},
		},
		{
			name: "request bodies ext spec",
			create: func() (string, any) {
				o := spec.NewRequestBodySpec()
				o.Spec.Spec.Description = "test"
				return "testRequestBodies", o.Spec
			},
			check: func(tb testing.TB, c *spec.Components) {
				require.Len(tb, c.RequestBodies, 1)
				require.NotNil(tb, c.RequestBodies["testRequestBodies"])
				require.NotNil(tb, c.RequestBodies["testRequestBodies"].Spec)
				require.NotNil(tb, c.RequestBodies["testRequestBodies"].Spec.Spec)
				require.Equal(tb, "test", c.RequestBodies["testRequestBodies"].Spec.Spec.Description)
			},
		},
		{
			name: "request bodies ref or spec",
			create: func() (string, any) {
				o := spec.NewRequestBodySpec()
				o.Spec.Spec.Description = "test"
				return "testRequestBodies", o
			},
			check: func(tb testing.TB, c *spec.Components) {
				require.Len(tb, c.RequestBodies, 1)
				require.NotNil(tb, c.RequestBodies["testRequestBodies"])
				require.NotNil(tb, c.RequestBodies["testRequestBodies"].Spec)
				require.NotNil(tb, c.RequestBodies["testRequestBodies"].Spec.Spec)
				require.Equal(tb, "test", c.RequestBodies["testRequestBodies"].Spec.Spec.Description)
			},
		},
		{
			name: "headers spec",
			create: func() (string, any) {
				o := spec.NewHeaderSpec()
				o.Spec.Spec.Description = "test"
				return "testHeader", o.Spec.Spec
			},
			check: func(tb testing.TB, c *spec.Components) {
				require.Len(tb, c.Headers, 1)
				require.NotNil(tb, c.Headers["testHeader"])
				require.NotNil(tb, c.Headers["testHeader"].Spec)
				require.NotNil(tb, c.Headers["testHeader"].Spec.Spec)
				require.Equal(tb, "test", c.Headers["testHeader"].Spec.Spec.Description)
			},
		},
		{
			name: "headers ext spec",
			create: func() (string, any) {
				o := spec.NewHeaderSpec()
				o.Spec.Spec.Description = "test"
				return "testHeader", o.Spec
			},
			check: func(tb testing.TB, c *spec.Components) {
				require.Len(tb, c.Headers, 1)
				require.NotNil(tb, c.Headers["testHeader"])
				require.NotNil(tb, c.Headers["testHeader"].Spec)
				require.NotNil(tb, c.Headers["testHeader"].Spec.Spec)
				require.Equal(tb, "test", c.Headers["testHeader"].Spec.Spec.Description)
			},
		},
		{
			name: "headers ref or spec",
			create: func() (string, any) {
				o := spec.NewHeaderSpec()
				o.Spec.Spec.Description = "test"
				return "testHeader", o
			},
			check: func(tb testing.TB, c *spec.Components) {
				require.Len(tb, c.Headers, 1)
				require.NotNil(tb, c.Headers["testHeader"])
				require.NotNil(tb, c.Headers["testHeader"].Spec)
				require.NotNil(tb, c.Headers["testHeader"].Spec.Spec)
				require.Equal(tb, "test", c.Headers["testHeader"].Spec.Spec.Description)
			},
		},
		{
			name: "security schemes spec",
			create: func() (string, any) {
				o := spec.NewSecuritySchemeSpec()
				o.Spec.Spec.Description = "test"
				return "testSecurityScheme", o.Spec.Spec
			},
			check: func(tb testing.TB, c *spec.Components) {
				require.Len(tb, c.SecuritySchemes, 1)
				require.NotNil(tb, c.SecuritySchemes["testSecurityScheme"])
				require.NotNil(tb, c.SecuritySchemes["testSecurityScheme"].Spec)
				require.NotNil(tb, c.SecuritySchemes["testSecurityScheme"].Spec.Spec)
				require.Equal(tb, "test", c.SecuritySchemes["testSecurityScheme"].Spec.Spec.Description)
			},
		},
		{
			name: "security schemes ext spec",
			create: func() (string, any) {
				o := spec.NewSecuritySchemeSpec()
				o.Spec.Spec.Description = "test"
				return "testSecurityScheme", o.Spec
			},
			check: func(tb testing.TB, c *spec.Components) {
				require.Len(tb, c.SecuritySchemes, 1)
				require.NotNil(tb, c.SecuritySchemes["testSecurityScheme"])
				require.NotNil(tb, c.SecuritySchemes["testSecurityScheme"].Spec)
				require.NotNil(tb, c.SecuritySchemes["testSecurityScheme"].Spec.Spec)
				require.Equal(tb, "test", c.SecuritySchemes["testSecurityScheme"].Spec.Spec.Description)
			},
		},
		{
			name: "security schemes ref or spec",
			create: func() (string, any) {
				o := spec.NewSecuritySchemeSpec()
				o.Spec.Spec.Description = "test"
				return "testSecurityScheme", o
			},
			check: func(tb testing.TB, c *spec.Components) {
				require.Len(tb, c.SecuritySchemes, 1)
				require.NotNil(tb, c.SecuritySchemes["testSecurityScheme"])
				require.NotNil(tb, c.SecuritySchemes["testSecurityScheme"].Spec)
				require.NotNil(tb, c.SecuritySchemes["testSecurityScheme"].Spec.Spec)
				require.Equal(tb, "test", c.SecuritySchemes["testSecurityScheme"].Spec.Spec.Description)
			},
		},
		{
			name: "link spec",
			create: func() (string, any) {
				o := spec.NewLinkSpec()
				o.Spec.Spec.Description = "test"
				return "testLink", o.Spec.Spec
			},
			check: func(tb testing.TB, c *spec.Components) {
				require.Len(tb, c.Links, 1)
				require.NotNil(tb, c.Links["testLink"])
				require.NotNil(tb, c.Links["testLink"].Spec)
				require.NotNil(tb, c.Links["testLink"].Spec.Spec)
				require.Equal(tb, "test", c.Links["testLink"].Spec.Spec.Description)
			},
		},
		{
			name: "link ext spec",
			create: func() (string, any) {
				o := spec.NewLinkSpec()
				o.Spec.Spec.Description = "test"
				return "testLink", o.Spec
			},
			check: func(tb testing.TB, c *spec.Components) {
				require.Len(tb, c.Links, 1)
				require.NotNil(tb, c.Links["testLink"])
				require.NotNil(tb, c.Links["testLink"].Spec)
				require.NotNil(tb, c.Links["testLink"].Spec.Spec)
				require.Equal(tb, "test", c.Links["testLink"].Spec.Spec.Description)
			},
		},
		{
			name: "link ref or spec",
			create: func() (string, any) {
				o := spec.NewLinkSpec()
				o.Spec.Spec.Description = "test"
				return "testLink", o
			},
			check: func(tb testing.TB, c *spec.Components) {
				require.Len(tb, c.Links, 1)
				require.NotNil(tb, c.Links["testLink"])
				require.NotNil(tb, c.Links["testLink"].Spec)
				require.NotNil(tb, c.Links["testLink"].Spec.Spec)
				require.Equal(tb, "test", c.Links["testLink"].Spec.Spec.Description)
			},
		},
		{
			name: "callback spec",
			create: func() (string, any) {
				p := spec.NewPathItemSpec()
				p.Spec.Spec.Description = "test"
				o := spec.NewCallbackSpec()
				o.Spec.Spec.WithPathItem("testPath", p)
				return "testCallback", o.Spec.Spec
			},
			check: func(tb testing.TB, c *spec.Components) {
				require.Len(tb, c.Callbacks, 1)
				require.NotNil(tb, c.Callbacks["testCallback"])
				require.NotNil(tb, c.Callbacks["testCallback"].Spec)
				require.NotNil(tb, c.Callbacks["testCallback"].Spec.Spec)
				paths := c.Callbacks["testCallback"].Spec.Spec.Callback
				require.Len(tb, paths, 1)
				require.NotNil(tb, paths["testPath"])
				require.NotNil(tb, paths["testPath"].Spec)
				require.NotNil(tb, paths["testPath"].Spec.Spec)
				require.Equal(tb, "test", paths["testPath"].Spec.Spec.Description)
			},
		},
		{
			name: "callback spec",
			create: func() (string, any) {
				p := spec.NewPathItemSpec()
				p.Spec.Spec.Description = "test"
				o := spec.NewCallbackSpec()
				o.Spec.Spec.WithPathItem("testPath", p)
				return "testCallback", o.Spec.Spec
			},
			check: func(tb testing.TB, c *spec.Components) {
				require.Len(tb, c.Callbacks, 1)
				require.NotNil(tb, c.Callbacks["testCallback"])
				require.NotNil(tb, c.Callbacks["testCallback"].Spec)
				require.NotNil(tb, c.Callbacks["testCallback"].Spec.Spec)
				paths := c.Callbacks["testCallback"].Spec.Spec.Callback
				require.Len(tb, paths, 1)
				require.NotNil(tb, paths["testPath"])
				require.NotNil(tb, paths["testPath"].Spec)
				require.NotNil(tb, paths["testPath"].Spec.Spec)
				require.Equal(tb, "test", paths["testPath"].Spec.Spec.Description)
			},
		},
		{
			name: "callback ext spec",
			create: func() (string, any) {
				p := spec.NewPathItemSpec()
				p.Spec.Spec.Description = "test"
				o := spec.NewCallbackSpec()
				o.Spec.Spec.WithPathItem("testPath", p)
				return "testCallback", o.Spec
			},
			check: func(tb testing.TB, c *spec.Components) {
				require.Len(tb, c.Callbacks, 1)
				require.NotNil(tb, c.Callbacks["testCallback"])
				require.NotNil(tb, c.Callbacks["testCallback"].Spec)
				require.NotNil(tb, c.Callbacks["testCallback"].Spec.Spec)
				paths := c.Callbacks["testCallback"].Spec.Spec.Callback
				require.Len(tb, paths, 1)
				require.NotNil(tb, paths["testPath"])
				require.NotNil(tb, paths["testPath"].Spec)
				require.NotNil(tb, paths["testPath"].Spec.Spec)
				require.Equal(tb, "test", paths["testPath"].Spec.Spec.Description)
			},
		},
		{
			name: "callback ref or spec",
			create: func() (string, any) {
				p := spec.NewPathItemSpec()
				p.Spec.Spec.Description = "test"
				o := spec.NewCallbackSpec()
				o.Spec.Spec.WithPathItem("testPath", p)
				return "testCallback", o
			},
			check: func(tb testing.TB, c *spec.Components) {
				require.Len(tb, c.Callbacks, 1)
				require.NotNil(tb, c.Callbacks["testCallback"])
				require.NotNil(tb, c.Callbacks["testCallback"].Spec)
				require.NotNil(tb, c.Callbacks["testCallback"].Spec.Spec)
				paths := c.Callbacks["testCallback"].Spec.Spec.Callback
				require.Len(tb, paths, 1)
				require.NotNil(tb, paths["testPath"])
				require.NotNil(tb, paths["testPath"].Spec)
				require.NotNil(tb, paths["testPath"].Spec.Spec)
				require.Equal(tb, "test", paths["testPath"].Spec.Spec.Description)
			},
		},
		{
			name: "path item spec",
			create: func() (string, any) {
				o := spec.NewPathItemSpec()
				o.Spec.Spec.Description = "test"
				return "testPathItem", o.Spec.Spec
			},
			check: func(tb testing.TB, c *spec.Components) {
				require.Len(tb, c.Paths, 1)
				require.NotNil(tb, c.Paths["testPathItem"])
				require.NotNil(tb, c.Paths["testPathItem"].Spec)
				require.NotNil(tb, c.Paths["testPathItem"].Spec.Spec)
				require.Equal(tb, "test", c.Paths["testPathItem"].Spec.Spec.Description)
			},
		},
		{
			name: "path item ext spec",
			create: func() (string, any) {
				o := spec.NewPathItemSpec()
				o.Spec.Spec.Description = "test"
				return "testPathItem", o.Spec
			},
			check: func(tb testing.TB, c *spec.Components) {
				require.Len(tb, c.Paths, 1)
				require.NotNil(tb, c.Paths["testPathItem"])
				require.NotNil(tb, c.Paths["testPathItem"].Spec)
				require.NotNil(tb, c.Paths["testPathItem"].Spec.Spec)
				require.Equal(tb, "test", c.Paths["testPathItem"].Spec.Spec.Description)
			},
		},
		{
			name: "path item ref or spec",
			create: func() (string, any) {
				o := spec.NewPathItemSpec()
				o.Spec.Spec.Description = "test"
				return "testPathItem", o
			},
			check: func(tb testing.TB, c *spec.Components) {
				require.Len(tb, c.Paths, 1)
				require.NotNil(tb, c.Paths["testPathItem"])
				require.NotNil(tb, c.Paths["testPathItem"].Spec)
				require.NotNil(tb, c.Paths["testPathItem"].Spec.Spec)
				require.Equal(tb, "test", c.Paths["testPathItem"].Spec.Spec.Description)
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			c := spec.NewComponents().Spec
			name, obj := tt.create()
			c.WithRefOrSpec(name, obj)
			tt.check(t, c)
		})
	}

	t.Run("panic", func(t *testing.T) {
		c := spec.NewComponents().Spec
		require.Panics(t, func() {
			c.WithRefOrSpec("panic", 42)
		})
	})
}

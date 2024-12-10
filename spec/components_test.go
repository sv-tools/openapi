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
				o := &spec.Schema{
					Title: "test",
				}
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
			name: "schema ref or spec",
			create: func() (string, any) {
				o := &spec.Schema{
					Title: "test",
				}
				return "testSchema", spec.NewRefOrSpec[spec.Schema](o)
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
				o := &spec.Response{
					Description: "test",
				}
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
			name: "response ext spec",
			create: func() (string, any) {
				o := &spec.Response{
					Description: "test",
				}
				return "testResponse", spec.NewExtendable[spec.Response](o)
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
				o := &spec.Response{
					Description: "test",
				}
				return "testResponse", spec.NewRefOrExtSpec[spec.Response](o)
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
				o := &spec.Parameter{
					Description: "test",
				}
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
			name: "parameter ext spec",
			create: func() (string, any) {
				o := &spec.Parameter{
					Description: "test",
				}
				return "testParameter", spec.NewExtendable[spec.Parameter](o)
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
				o := &spec.Parameter{
					Description: "test",
				}
				return "testParameter", spec.NewRefOrExtSpec[spec.Parameter](o)
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
				o := &spec.Example{
					Description: "test",
				}
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
			name: "examples ext spec",
			create: func() (string, any) {
				o := &spec.Example{
					Description: "test",
				}
				return "testExamples", spec.NewExtendable[spec.Example](o)
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
				o := &spec.Example{
					Description: "test",
				}
				return "testExamples", spec.NewRefOrExtSpec[spec.Example](o)
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
			name: "request body spec",
			create: func() (string, any) {
				o := &spec.RequestBody{
					Description: "test",
				}
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
			name: "request body ext spec",
			create: func() (string, any) {
				o := &spec.RequestBody{
					Description: "test",
				}
				return "testRequestBodies", spec.NewExtendable[spec.RequestBody](o)
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
			name: "request body ref or spec",
			create: func() (string, any) {
				o := &spec.RequestBody{
					Description: "test",
				}
				return "testRequestBodies", spec.NewRefOrExtSpec[spec.RequestBody](o)
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
				o := &spec.Header{
					Description: "test",
				}
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
			name: "headers ext spec",
			create: func() (string, any) {
				o := &spec.Header{
					Description: "test",
				}
				return "testHeader", spec.NewExtendable[spec.Header](o)
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
				o := &spec.Header{
					Description: "test",
				}
				return "testHeader", spec.NewRefOrExtSpec[spec.Header](o)
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
				o := &spec.SecurityScheme{
					Description: "test",
				}
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
			name: "security schemes ext spec",
			create: func() (string, any) {
				o := &spec.SecurityScheme{
					Description: "test",
				}
				return "testSecurityScheme", spec.NewExtendable[spec.SecurityScheme](o)
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
				o := &spec.SecurityScheme{
					Description: "test",
				}
				return "testSecurityScheme", spec.NewRefOrExtSpec[spec.SecurityScheme](o)
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
				o := &spec.Link{
					Description: "test",
				}
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
			name: "link ext spec",
			create: func() (string, any) {
				o := &spec.Link{
					Description: "test",
				}
				return "testLink", spec.NewExtendable[spec.Link](o)
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
				o := &spec.Link{
					Description: "test",
				}
				return "testLink", spec.NewRefOrExtSpec[spec.Link](o)
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
				o := &spec.Callback{
					Callback: map[string]*spec.RefOrSpec[spec.Extendable[spec.PathItem]]{
						"testPath": spec.NewRefOrExtSpec[spec.PathItem](&spec.PathItem{
							Description: "test",
						}),
					},
				}
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
			name: "callback ext spec",
			create: func() (string, any) {
				o := &spec.Callback{
					Callback: map[string]*spec.RefOrSpec[spec.Extendable[spec.PathItem]]{
						"testPath": spec.NewRefOrExtSpec[spec.PathItem](&spec.PathItem{
							Description: "test",
						}),
					},
				}
				return "testCallback", spec.NewExtendable[spec.Callback](o)
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
				o := &spec.Callback{
					Callback: map[string]*spec.RefOrSpec[spec.Extendable[spec.PathItem]]{
						"testPath": spec.NewRefOrExtSpec[spec.PathItem](&spec.PathItem{
							Description: "test",
						}),
					},
				}
				return "testCallback", spec.NewRefOrExtSpec[spec.Callback](o)
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
				o := &spec.PathItem{
					Description: "test",
				}
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
		{
			name: "path item ext spec",
			create: func() (string, any) {
				o := &spec.PathItem{
					Description: "test",
				}
				return "testPathItem", spec.NewExtendable[spec.PathItem](o)
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
				o := &spec.PathItem{
					Description: "test",
				}
				return "testPathItem", spec.NewRefOrExtSpec[spec.PathItem](o)
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
			name, obj := tt.create()
			tt.check(t, (&spec.Components{}).WithRefOrSpec(name, obj))
		})
	}

	t.Run("panic", func(t *testing.T) {
		require.Panics(t, func() {
			(&spec.Components{}).WithRefOrSpec("panic", 42)
		})
	})
}

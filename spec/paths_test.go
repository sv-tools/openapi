package spec_test

import (
	"github.com/stretchr/testify/require"
	"github.com/sv-tools/openapi/spec"
	"testing"
)

func TestPaths_WithPathItem(t *testing.T) {
	for _, tt := range []struct {
		name   string
		create func() any
	}{
		{
			name: "path item spec",
			create: func() any {
				o := spec.NewPathItemSpec()
				o.Spec.Spec.Description = "test"
				return o.Spec.Spec
			},
		},
		{
			name: "path item ext spec",
			create: func() any {
				o := spec.NewPathItemSpec()
				o.Spec.Spec.Description = "test"
				return o.Spec
			},
		},
		{
			name: "path item ref or spec",
			create: func() any {
				o := spec.NewPathItemSpec()
				o.Spec.Spec.Description = "test"
				return o
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			c := spec.NewPaths().Spec
			c.WithPathItem("testPathItem", tt.create())
			require.Len(t, c.Paths, 1)
			require.NotNil(t, c.Paths["testPathItem"])
			require.NotNil(t, c.Paths["testPathItem"].Spec)
			require.NotNil(t, c.Paths["testPathItem"].Spec.Spec)
			require.Equal(t, "test", c.Paths["testPathItem"].Spec.Spec.Description)
		})
	}

	t.Run("panic", func(t *testing.T) {
		c := spec.NewPaths().Spec
		require.Panics(t, func() {
			c.WithPathItem("panic", 42)
		})
	})
}

package validate

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/santhosh-tekuri/jsonschema/v5"
	"gopkg.in/yaml.v3"

	"github.com/sv-tools/openapi/spec"
)

var (
	//go:embed schemas/*.json
	schemas   embed.FS
	resources = map[string]string{
		"schemas/v3.1.json": "https://spec.openapis.org/oas/3.1/schema/2021-09-28",
	}
	mainURL = "https://spec.openapis.org/oas/3.1/schema/2021-09-28"

	DefaultSchema   *jsonschema.Schema
	DefaultCompiler *jsonschema.Compiler
)

func init() {
	DefaultCompiler = jsonschema.NewCompiler()
	for filename, url := range resources {
		r, err := schemas.Open(filename)
		if err != nil {
			panic(err)
		}
		if err := DefaultCompiler.AddResource(url, r); err != nil {
			panic(err)
		}
		if err := r.Close(); err != nil {
			panic(err)
		}
	}
	DefaultSchema = DefaultCompiler.MustCompile(mainURL)
}

// Json validates a given data in JSON format against the DefaultSchema.
func Json(data []byte) error {
	var v any
	if err := json.Unmarshal(data, &v); err != nil {
		return fmt.Errorf("unable to unmarshal json: %w", err)
	}
	if err := DefaultSchema.Validate(v); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}
	return nil
}

// Yaml validates a given data in YAML format against the DefaultSchema.
func Yaml(data []byte) error {
	var v any
	if err := yaml.Unmarshal(data, &v); err != nil {
		return fmt.Errorf("unable to unmarshal yaml: %w", err)
	}
	if err := DefaultSchema.Validate(v); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}
	return nil
}

// Spec validates a given spec against the DefaultSchema.
func Spec(o *spec.Extendable[spec.OpenAPI]) error {
	data, err := json.Marshal(&o)
	if err != nil {
		return fmt.Errorf("unable to marshal json: %w", err)
	}
	return Json(data)
}

// Report converts given error to Json format.
func Report(err error, detailed bool) string {
	var (
		target *jsonschema.ValidationError
		out    any
	)
	if errors.As(err, &target) {
		if detailed {
			out = target.DetailedOutput()
		} else {
			out = target.BasicOutput()
		}
	} else {
		out = err
	}

	data, err := json.MarshalIndent(&out, "", "  ")
	if err != nil {
		return "unable to prepare report: " + err.Error()
	}
	return string(data)
}

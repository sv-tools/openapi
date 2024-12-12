package openapi

import (
	"bytes"
	"errors"
	"fmt"
	"net/mail"
	"net/url"
	"strings"
	"sync"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

// Validatable is an interface for validating the specification.
type validatable interface {
	// an unexported method to be used by ValidateSpec function
	validateSpec(path string, opts *specValidationOptions) []*validationError
}

type visitedObjects map[string]bool

func (o visitedObjects) String() string {
	keys := make([]string, 0, len(o))
	for k := range o {
		keys = append(keys, k)
	}
	return strings.Join(keys, ", ")
}

type specValidationOptions struct {
	spec                            *Extendable[OpenAPI]
	visited                         visitedObjects
	linkToOperationID               map[string]string
	allowExtensionNameWithoutPrefix bool
	allowRequestBodyForGet          bool
	allowRequestBodyForHead         bool
	allowRequestBodyForDelete       bool
	allowUndefinedTagsInOperation   bool
	doNotValidateExamples           bool
	doNotValidateDefaultValues      bool
}

func newSpecValidationOptions(spec *Extendable[OpenAPI], opts ...SpecValidationOption) *specValidationOptions {
	options := &specValidationOptions{
		spec:              spec,
		visited:           make(visitedObjects),
		linkToOperationID: make(map[string]string),
	}
	for _, opt := range opts {
		opt(options)
	}

	return options
}

// SpecValidationOption is a type for validation options.
type SpecValidationOption func(*specValidationOptions)

// AllowExtensionNameWithoutPrefix is a validation option to allow extension name without `x-` prefix.
func AllowExtensionNameWithoutPrefix() SpecValidationOption {
	return func(v *specValidationOptions) {
		v.allowExtensionNameWithoutPrefix = true
	}
}

func AllowRequestBodyForGet() SpecValidationOption {
	return func(v *specValidationOptions) {
		v.allowRequestBodyForGet = true
	}
}

func AllowRequestBodyForHead() SpecValidationOption {
	return func(v *specValidationOptions) {
		v.allowRequestBodyForHead = true
	}
}

func AllowRequestBodyForDelete() SpecValidationOption {
	return func(v *specValidationOptions) {
		v.allowRequestBodyForDelete = true
	}
}

func AllowUndefinedTagsInOperation() SpecValidationOption {
	return func(v *specValidationOptions) {
		v.allowUndefinedTagsInOperation = true
	}
}

func DoNotValidateExamples() SpecValidationOption {
	return func(v *specValidationOptions) {
		v.doNotValidateExamples = true
	}
}

func DoNotValidateDefaultValues() SpecValidationOption {
	return func(v *specValidationOptions) {
		v.doNotValidateDefaultValues = true
	}
}

type validationError struct {
	path string
	err  error
}

func newValidationError(path string, err any, args ...any) *validationError {
	switch e := err.(type) {
	case error:
		return &validationError{path: path, err: e}
	case string:
		return &validationError{path: path, err: fmt.Errorf(e, args...)}
	default:
		// unreachable
		panic(fmt.Sprintf("unsupported error type: %T", e))
	}
}

func joinDot(path ...string) string {
	switch len(path) {
	case 0:
		return ""
	case 1:
		return path[0]
	default:
	}
	if path[0] == "" {
		path = path[1:]
	}
	return strings.Join(path, ".")
}

func joinArrayItem(path string, v any) string {
	switch t := v.(type) {
	case int:
		return fmt.Sprintf("%s[%d]", path, t)
	case string:
		return fmt.Sprintf("%s['%s']", path, t)
	default:
		return fmt.Sprintf("%s[%v]", path, t)
	}
}

func (e *validationError) Error() string {
	return fmt.Sprintf("%s: %s", e.path, e.err)
}

func (e *validationError) Unwrap() error {
	return e.err
}

var (
	ErrRequired          = errors.New("required")
	ErrMutuallyExclusive = errors.New("mutually exclusive")
	ErrUnused            = errors.New("unused")
)

func checkURL(value string) error {
	if value == "" {
		return nil
	}
	if _, err := url.Parse(value); err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}
	return nil
}

func checkEmail(value string) error {
	if value == "" {
		return nil
	}
	if _, err := mail.ParseAddress(value); err != nil {
		return fmt.Errorf("invalid email: %w", err)
	}
	return nil
}

// ValidateData validates given object against the given schema.
// Warning: this function is not implemented yet.
func ValidateData(value any, schema *Schema, spec *Extendable[OpenAPI]) error {
	if spec == nil {
		return errors.New("given Schema cannot be nil")
	}

	return errors.New("data validation not implemented")
}

func ValidateSpec(spec *Extendable[OpenAPI], opts ...SpecValidationOption) error {
	options := newSpecValidationOptions(spec, opts...)

	if errs := spec.validateSpec("", options); len(errs) > 0 {
		joinErrors := make([]error, len(errs))
		for i := range errs {
			joinErrors[i] = errs[i]
		}
		return errors.Join(joinErrors...)
	}

	return nil
}

// CompilerOption is a type to modify the jsonschema.Compiler.
type CompilerOption func(*jsonschema.Compiler)

type DataValidator struct {
	compiler *jsonschema.Compiler
	schemas  sync.Map
	mu       sync.Mutex
}

const specPrefix = "http://spec#"

func NewDataValidator(spec *Extendable[OpenAPI], opts ...CompilerOption) (*DataValidator, error) {
	data, err := spec.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("marshaling spec failed: %w", err)
	}
	doc, err := jsonschema.UnmarshalJSON(bytes.NewReader(data))
	c := jsonschema.NewCompiler()
	c.DefaultDraft(jsonschema.Draft2020)
	for _, opt := range opts {
		opt(c)
	}
	if err := c.AddResource(specPrefix, doc); err != nil {
		return nil, fmt.Errorf("adding spec to compiler failed: %w", err)
	}
	return &DataValidator{
		compiler: c,
		schemas:  sync.Map{},
	}, nil
}

func (v *DataValidator) Validate(loc string, value any) error {
	var schema *jsonschema.Schema
	if s, ok := v.schemas.Load(loc); ok {
		schema = s.(*jsonschema.Schema)
	} else {
		var err error
		schema, err = func() (*jsonschema.Schema, error) {
			v.mu.Lock()
			defer v.mu.Unlock()
			if s, ok := v.schemas.Load(loc); ok {
				return s.(*jsonschema.Schema), nil
			} else {
				schema, err := v.compiler.Compile(specPrefix + loc)
				if err != nil {
					return nil, fmt.Errorf("compiling spec for given location %q failed: %w", loc, err)
				}
				v.schemas.Store(loc, schema)
				return schema, nil
			}
		}()
		if err != nil {
			return err
		}
	}
	return schema.Validate(value)
}

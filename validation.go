package openapi

import (
	"bytes"
	"encoding/json"
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
	validateSpec(location string, opts *specValidationOptions) []*validationError
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
	validator                       *Validator
	visited                         visitedObjects
	linkToOperationID               map[string]string
	allowExtensionNameWithoutPrefix bool
	allowRequestBodyForGet          bool
	allowRequestBodyForHead         bool
	allowRequestBodyForDelete       bool
	allowUndefinedTagsInOperation   bool
	allowUnusedComponents           bool
	doNotValidateExamples           bool
	doNotValidateDefaultValues      bool
}

func newSpecValidationOptions(validator *Validator, opts ...SpecValidationOption) *specValidationOptions {
	options := &specValidationOptions{
		validator:         validator,
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

// AllowRequestBodyForGet is a validation option to allow request body for GET operation.
func AllowRequestBodyForGet() SpecValidationOption {
	return func(v *specValidationOptions) {
		v.allowRequestBodyForGet = true
	}
}

// AllowRequestBodyForHead is a validation option to allow request body for HEAD operation.
func AllowRequestBodyForHead() SpecValidationOption {
	return func(v *specValidationOptions) {
		v.allowRequestBodyForHead = true
	}
}

// AllowRequestBodyForDelete is a validation option to allow request body for DELETE operation.
func AllowRequestBodyForDelete() SpecValidationOption {
	return func(v *specValidationOptions) {
		v.allowRequestBodyForDelete = true
	}
}

// AllowUndefinedTagsInOperation is a validation option to allow undefined tags in operation.
func AllowUndefinedTagsInOperation() SpecValidationOption {
	return func(v *specValidationOptions) {
		v.allowUndefinedTagsInOperation = true
	}
}

// AllowUnusedComponents is a validation option to allow unused components.
func AllowUnusedComponents() SpecValidationOption {
	return func(v *specValidationOptions) {
		v.allowUnusedComponents = true
	}
}

// DoNotValidateExamples is a validation option to skip examples validation.
func DoNotValidateExamples() SpecValidationOption {
	return func(v *specValidationOptions) {
		v.doNotValidateExamples = true
	}
}

// DoNotValidateDefaultValues is a validation option to skip default values validation.
func DoNotValidateDefaultValues() SpecValidationOption {
	return func(v *specValidationOptions) {
		v.doNotValidateDefaultValues = true
	}
}

type validationError struct {
	location string
	err      error
}

func newValidationError(location string, err any, args ...any) *validationError {
	switch e := err.(type) {
	case error:
		return &validationError{location: location, err: e}
	case string:
		return &validationError{location: location, err: fmt.Errorf(e, args...)}
	default:
		// unreachable
		panic(fmt.Sprintf("unsupported error type: %T", e))
	}
}

var jsonPointerEscaper = strings.NewReplacer("~", "~0", "/", "~1")

func joinLoc(base string, parts ...any) string {
	if len(parts) == 0 {
		return base
	}

	elems := append(make([]string, 0, len(parts)+1), base)
	for _, v := range parts {
		elems = append(elems, jsonPointerEscaper.Replace(fmt.Sprintf("%v", v)))
	}

	return strings.Join(elems, "/")
}

func (e *validationError) Error() string {
	return fmt.Sprintf("%s: %s", e.location, e.err)
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

// CompilerOption is a type to modify the jsonschema.Compiler.
type CompilerOption func(*jsonschema.Compiler)

// Validator is a struct for validating the OpenAPI specification and a data.
type Validator struct {
	spec     *Extendable[OpenAPI]
	compiler *jsonschema.Compiler
	schemas  sync.Map
	mu       sync.Mutex
}

const specPrefix = "http://spec"

// NewValidator creates an instance of Validator struct.
//
// The function creates new jsonschema comppiler and adds the given spec to the compiler.
func NewValidator(spec *Extendable[OpenAPI], opts ...CompilerOption) (*Validator, error) {
	data, err := json.Marshal(spec)
	if err != nil {
		return nil, fmt.Errorf("marshaling spec failed: %w", err)
	}
	doc, err := jsonschema.UnmarshalJSON(bytes.NewReader(data))
	c := jsonschema.NewCompiler()
	c.DefaultDraft(jsonschema.Draft2020)
	if err := c.AddResource(specPrefix, doc); err != nil {
		return nil, fmt.Errorf("adding spec to compiler failed: %w", err)
	}
	for _, opt := range opts {
		opt(c)
	}
	return &Validator{
		spec:     spec,
		compiler: c,
		schemas:  sync.Map{},
	}, nil
}

// ValidateSpec validates the specification.
func (v *Validator) ValidateSpec(opts ...SpecValidationOption) error {
	if errs := v.spec.validateSpec("", newSpecValidationOptions(v, opts...)); len(errs) > 0 {
		joinErrors := make([]error, len(errs))
		for i := range errs {
			joinErrors[i] = errs[i]
		}
		return errors.Join(joinErrors...)
	}

	return nil
}

// ValidateData validates the given value against the schema located at the given location.
// The location should be in form of JSON Pointer.
func (v *Validator) ValidateData(location string, value any) error {
	var schema *jsonschema.Schema
	if s, ok := v.schemas.Load(location); ok {
		schema = s.(*jsonschema.Schema)
	} else {
		var err error
		schema, err = func() (*jsonschema.Schema, error) {
			v.mu.Lock()
			defer v.mu.Unlock()
			if s, ok := v.schemas.Load(location); ok {
				return s.(*jsonschema.Schema), nil
			} else {
				if !strings.HasPrefix(location, "#") {
					location = "#" + location
				}
				schema, err := v.compiler.Compile(specPrefix + location)
				if err != nil {
					return nil, fmt.Errorf("compiling spec for given location %q failed: %w", location, err)
				}
				v.schemas.Store(location, schema)
				return schema, nil
			}
		}()
		if err != nil {
			return err
		}
	}
	return schema.Validate(value)
}

// ValidateDataAsJSON marshal and unmarshals the given value to JSON and
// validates it against the schema located at the given location.
func (v *Validator) ValidateDataAsJSON(location string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshaling value failed: %w", err)
	}
	doc, err := jsonschema.UnmarshalJSON(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("unmarshaling value failed: %w", err)
	}
	return v.ValidateData(location, doc)
}

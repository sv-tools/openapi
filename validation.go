package openapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/mail"
	"net/url"
	"reflect"
	"strings"
	"sync"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

// Validatable is an interface for validating the specification.
type validatable interface {
	// an unexported method to be used by ValidateSpec function
	validateSpec(location string, validator *Validator) []*validationError
}

type visitedObjects map[string]bool

func (o visitedObjects) String() string {
	keys := make([]string, 0, len(o))
	for k := range o {
		keys = append(keys, k)
	}
	return strings.Join(keys, ", ")
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

// Validator is a struct for validating the OpenAPI specification and a data.
type Validator struct {
	spec *Extendable[OpenAPI]

	compiler *jsonschema.Compiler
	schemas  sync.Map
	mu       sync.Mutex

	opts              *validationOptions
	visited           visitedObjects
	linkToOperationID map[string]string
}

const specPrefix = "http://spec"

// NewValidator creates an instance of Validator struct.
//
// The function creates new jsonschema comppiler and adds the given spec to the compiler.
func NewValidator(spec *Extendable[OpenAPI], opts ...ValidationOption) (*Validator, error) {
	options := &validationOptions{}
	for _, opt := range opts {
		opt(options)
	}
	validator := &Validator{
		spec:    spec,
		schemas: sync.Map{},
		opts:    options,
	}
	data, err := json.Marshal(spec)
	if err != nil {
		return nil, fmt.Errorf("marshaling spec failed: %w", err)
	}
	doc, err := jsonschema.UnmarshalJSON(bytes.NewReader(data))
	compiler := jsonschema.NewCompiler()
	compiler.DefaultDraft(jsonschema.Draft2020)
	if err := compiler.AddResource(specPrefix, doc); err != nil {
		return nil, fmt.Errorf("adding spec to compiler failed: %w", err)
	}
	for _, f := range validator.opts.updateCompiler {
		f(compiler)
	}
	validator.compiler = compiler
	return validator, nil
}

// ValidateSpec validates the specification.
func (v *Validator) ValidateSpec() error {
	// clear visited objects
	v.visited = make(visitedObjects)
	v.linkToOperationID = make(map[string]string)

	if errs := v.spec.validateSpec("", v); len(errs) > 0 {
		joinErrors := make([]error, len(errs))
		for i := range errs {
			joinErrors[i] = errs[i]
		}
		return errors.Join(joinErrors...)
	}

	return nil
}

// ValidateData validates the given value against the schema located at the given location.
//
// The location should be in form of JSON Pointer.
// The value can be a struct, a string containing JSON, or any other types.
// If the value is a struct, it will be marshaled and unmarshaled to JSON.
func (v *Validator) ValidateData(location string, value any) error {
	var schema *jsonschema.Schema
	if s, ok := v.schemas.Load(location); ok {
		schema = s.(*jsonschema.Schema)
	} else {
		var err error
		// use lambda to simplify the mutex unlocking code after the schema is compiled
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

	switch getKind(value) {
	case reflect.Struct:
		// jsonschema does not support struct, so we need to marshal and unmarshal
		// the value to JSON representation (map[any]struct).
		data, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("marshaling value failed: %w", err)
		}
		value, err = jsonschema.UnmarshalJSON(bytes.NewReader(data))
		if err != nil {
			return fmt.Errorf("unmarshaling value failed: %w", err)
		}
	case reflect.String:
		if v.opts.validateDataAsJSON {
			// check if the value is already a JSON, if not keep it as is.
			s, err := jsonschema.UnmarshalJSON(strings.NewReader(value.(string)))
			if err == nil {
				value = s
			}
		}
	}
	return schema.Validate(value)
}

// ValidateDataAsJSON marshal and unmarshals the given value to JSON and
// validates it against the schema located at the given location.
//
// If the value is a string, it will be unmarshaled to JSON first, if failed it will be kept as is.
func (v *Validator) ValidateDataAsJSON(location string, value any) error {
	switch getKind(value) {
	// marshal and unmarshal the value to JSON representation (map[any]struct).
	case reflect.Struct:
		data, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("marshaling value failed: %w", err)
		}
		value, err = jsonschema.UnmarshalJSON(bytes.NewReader(data))
		if err != nil {
			return fmt.Errorf("unmarshaling value failed: %w", err)
		}
	// check if the value is already a JSON, if not keep it as is.
	case reflect.String:
		s, err := jsonschema.UnmarshalJSON(strings.NewReader(value.(string)))
		if err == nil {
			value = s
		}
	}
	return v.ValidateData(location, value)
}

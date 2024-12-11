package openapi

import (
	"errors"
	"fmt"
	"net/mail"
	"net/url"
	"strings"
)

// Validatable is an interface for validating the specification.
type validatable interface {
	// an unexported method to be used by ValidateSpec function
	validateSpec(path string, opts *validationOptions) []*validationError
}

type visitedObjects map[string]bool

func (o visitedObjects) String() string {
	keys := make([]string, 0, len(o))
	for k := range o {
		keys = append(keys, k)
	}
	return strings.Join(keys, ", ")
}

type validationOptions struct {
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

func newValidationOptions(spec *Extendable[OpenAPI], opts ...ValidationOption) *validationOptions {
	options := &validationOptions{
		spec:              spec,
		visited:           make(visitedObjects),
		linkToOperationID: make(map[string]string),
	}
	for _, opt := range opts {
		opt(options)
	}

	return options
}

// ValidationOption is a type for validation options.
type ValidationOption func(*validationOptions)

// AllowExtensionNameWithoutPrefix is a validation option to allow extension name without `x-` prefix.
func AllowExtensionNameWithoutPrefix() ValidationOption {
	return func(v *validationOptions) {
		v.allowExtensionNameWithoutPrefix = true
	}
}

func AllowRequestBodyForGet() ValidationOption {
	return func(v *validationOptions) {
		v.allowRequestBodyForGet = true
	}
}

func AllowRequestBodyForHead() ValidationOption {
	return func(v *validationOptions) {
		v.allowRequestBodyForHead = true
	}
}

func AllowRequestBodyForDelete() ValidationOption {
	return func(v *validationOptions) {
		v.allowRequestBodyForDelete = true
	}
}

func AllowUndefinedTagsInOperation() ValidationOption {
	return func(v *validationOptions) {
		v.allowUndefinedTagsInOperation = true
	}
}

func DoNotValidateExamples() ValidationOption {
	return func(v *validationOptions) {
		v.doNotValidateExamples = true
	}
}

func DoNotValidateDefaultValues() ValidationOption {
	return func(v *validationOptions) {
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

// ValidateSpec validates only Extendable[OpenAPI] object.
func ValidateSpec(spec *Extendable[OpenAPI], opts ...ValidationOption) error {
	options := newValidationOptions(spec, opts...)

	if errs := spec.validateSpec("", options); len(errs) > 0 {
		joinErrors := make([]error, len(errs))
		for i := range errs {
			joinErrors[i] = errs[i]
		}
		return errors.Join(joinErrors...)
	}

	return nil
}

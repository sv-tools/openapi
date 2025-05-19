package openapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

const ExtensionPrefix = "x-"

// Extendable allows extensions to the OpenAPI Schema.
// The field name MUST begin with `x-`, for example, `x-internal-id`.
// Field names beginning `x-oai-` and `x-oas-` are reserved for uses defined by the OpenAPI Initiative.
// The value can be null, a primitive, an array or an object.
//
// https://spec.openapis.org/oas/v3.1.1#specification-extensions
//
// Example:
//
//	  openapi: 3.1.1
//	  info:
//	    title: Sample Pet Store App
//	    summary: A pet store manager.
//	    description: This is a sample server for a pet store.
//	    version: 1.0.1
//	    x-build-data: 2006-01-02T15:04:05Z07:00
//		x-build-commit-id: dac33af14d0d4a5f1c226141042ca7cefc6aeb75
type Extendable[T any] struct {
	Spec       *T             `json:"-" yaml:"-"`
	Extensions map[string]any `json:"-" yaml:"-"`
}

// NewExtendable creates new Extendable object for given spec
func NewExtendable[T any](spec *T) *Extendable[T] {
	ext := Extendable[T]{
		Spec:       spec,
		Extensions: make(map[string]any),
	}
	return &ext
}

// AddExt sets the extension and returns the current object.
// The `x-` prefix will be added automatically to given name.
func (o *Extendable[T]) AddExt(name string, value any) *Extendable[T] {
	if o.Extensions == nil {
		o.Extensions = make(map[string]any, 1)
	}
	if !strings.HasPrefix(name, ExtensionPrefix) {
		name = ExtensionPrefix + name
	}
	o.Extensions[name] = value
	return o
}

// GetExt returns the extension value by name.
// The `x-` prefix will be added automatically to given name.
func (o *Extendable[T]) GetExt(name string) any {
	if o.Extensions == nil {
		return nil
	}
	if !strings.HasPrefix(name, ExtensionPrefix) {
		name = ExtensionPrefix + name
	}
	return o.Extensions[name]
}

// MarshalJSON implements json.Marshaler interface.
func (o *Extendable[T]) MarshalJSON() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	var raw map[string]json.RawMessage
	if len(o.Extensions) > 0 {
		exts, err := json.Marshal(&o.Extensions)
		if err != nil {
			return nil, fmt.Errorf("%T.Extensions: %w", o.Spec, err)
		}
		if err := json.Unmarshal(exts, &raw); err != nil {
			return nil, fmt.Errorf("%T(raw extensions): %w", o.Spec, err)
		}
	}
	if o.Spec != nil {
		fields, err := json.Marshal(o.Spec)
		if err != nil {
			return nil, fmt.Errorf("%T: %w", o.Spec, err)
		}
		if err := json.Unmarshal(fields, &raw); err != nil {
			return nil, fmt.Errorf("%T(raw fields): %w", o.Spec, err)
		}
	}
	if len(raw) == 0 {
		return nil, nil
	}
	data, err := json.Marshal(&raw)
	if err != nil {
		return nil, fmt.Errorf("%T(raw): %w", o.Spec, err)
	}
	return data, nil
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (o *Extendable[T]) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("%T: %w", o.Spec, err)
	}
	exts := make(map[string]any)
	for name, value := range raw {
		if strings.HasPrefix(name, ExtensionPrefix) {
			var v any
			if err := json.Unmarshal(value, &v); err != nil {
				return fmt.Errorf("%T.Extensions.%s: %w", o.Spec, name, err)
			}
			exts[name] = v
		}
	}
	if len(exts) > 0 {
		o.Extensions = exts
		for name := range exts {
			delete(raw, name)
		}
	}
	fields, err := json.Marshal(&raw)
	if err != nil {
		return fmt.Errorf("%T(raw): %w", o.Spec, err)
	}
	if err := json.Unmarshal(fields, &o.Spec); err != nil {
		return fmt.Errorf("%T: %w", o.Spec, err)
	}

	return nil
}

// MarshalYAML implements yaml.Marshaler interface.
func (o *Extendable[T]) MarshalYAML() (any, error) {
	data, err := json.Marshal(o)
	if err != nil {
		return nil, fmt.Errorf("%T: %w", o.Spec, err)
	}
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("%T(raw): %w", o.Spec, err)
	}
	return raw, nil
}

// UnmarshalYAML implements yaml.obsoleteUnmarshaler and goyaml.InterfaceUnmarshaler interfaces.
func (o *Extendable[T]) UnmarshalYAML(unmarshal func(any) error) error {
	var raw map[string]any
	if err := unmarshal(&raw); err != nil {
		return fmt.Errorf("%T: %w", o.Spec, err)
	}
	o.Extensions = make(map[string]any)
	for name, value := range raw {
		if strings.HasPrefix(name, ExtensionPrefix) {
			o.Extensions[name] = value
			delete(raw, name)
		}
	}
	if o.Spec == nil {
		o.Spec = new(T)
	}
	if err := unmarshal(o.Spec); err != nil {
		o.Spec = nil
		return fmt.Errorf("%T: %w", o.Spec, err)
	}
	return nil
}

var ErrExtensionNameMustStartWithPrefix = errors.New("extension name must start with `" + ExtensionPrefix + "`")

const unsupportedSpecTypePrefix = "unsupported spec type: "

type UnsupportedSpecTypeError string

func (e UnsupportedSpecTypeError) Error() string {
	return unsupportedSpecTypePrefix + string(e)
}

func (e UnsupportedSpecTypeError) Is(target error) bool {
	return strings.HasPrefix(target.Error(), unsupportedSpecTypePrefix)
}

func NewUnsupportedSpecTypeError(spec any) error {
	return UnsupportedSpecTypeError(fmt.Sprintf("%T", spec))
}

func (o *Extendable[T]) validateSpec(location string, validator *Validator) []*validationError {
	var errs []*validationError
	if o.Spec != nil {
		if spec, ok := any(o.Spec).(validatable); ok {
			errs = append(errs, spec.validateSpec(location, validator)...)
		} else {
			errs = append(errs, newValidationError(location, NewUnsupportedSpecTypeError(o.Spec)))
		}
	}
	if validator.opts.allowExtensionNameWithoutPrefix {
		return errs
	}

	for name := range o.Extensions {
		if !strings.HasPrefix(name, ExtensionPrefix) {
			errs = append(errs, newValidationError(joinLoc(location, name), ErrExtensionNameMustStartWithPrefix))
		}
	}
	return errs
}

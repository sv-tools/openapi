package openapi

import (
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// Ref is a simple object to allow referencing other components in the OpenAPI document, internally and externally.
// The $ref string value contains a URI [RFC3986], which identifies the location of the value being referenced.
// See the rules for resolving Relative References.
//
// https://spec.openapis.org/oas/v3.1.1#reference-object
//
// Example:
//
//	$ref: '#/components/schemas/Pet'
type Ref struct {
	// REQUIRED.
	// The reference identifier.
	// This MUST be in the form of a URI.
	Ref string `json:"$ref" yaml:"$ref"`
	// A short summary which by default SHOULD override that of the referenced component.
	// If the referenced object-type does not allow a summary field, then this field has no effect.
	Summary string `json:"summary,omitempty" yaml:"summary,omitempty"`
	// A description which by default SHOULD override that of the referenced component.
	// CommonMark syntax MAY be used for rich text representation.
	// If the referenced object-type does not allow a description field, then this field has no effect.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
}

// RefOrSpec holds either Ref or any OpenAPI spec type.
//
// NOTE: The Ref object takes precedent over Spec if using json or yaml Marshal and Unmarshal functions.
type RefOrSpec[T any] struct {
	Ref  *Ref `json:"-" yaml:"-"`
	Spec *T   `json:"-" yaml:"-"`
}

// NewRefOrSpec creates an object of RefOrSpec type from given Ref or string or any form of Spec.
func NewRefOrSpec[T any](v any) *RefOrSpec[T] {
	o := RefOrSpec[T]{}
	switch t := v.(type) {
	case *Ref:
		o.Ref = t
	case Ref:
		o.Ref = &t
	case string:
		o.Ref = &Ref{Ref: t}
	case *T:
		o.Spec = t
	case T:
		o.Spec = &t
	case nil:
	}
	return &o
}

// NewRefOrExtSpec creates an object of RefOrSpec[Extendable[any]] type from given Ref or string or any form of Spec.
func NewRefOrExtSpec[T any](v any) *RefOrSpec[Extendable[T]] {
	o := RefOrSpec[Extendable[T]]{}
	switch t := v.(type) {
	case *Ref:
		o.Ref = t
	case Ref:
		o.Ref = &t
	case string:
		o.Ref = &Ref{Ref: t}
	case *T:
		o.Spec = NewExtendable[T](t)
	case T:
		o.Spec = NewExtendable[T](&t)
	case nil:
	}
	return &o
}

func (o *RefOrSpec[T]) getLocationOrRef(location string) string {
	if o.Ref != nil {
		return o.Ref.Ref
	}
	return location
}

// GetSpec return a Spec if it is set or loads it from Components in case of Ref or an error
func (o *RefOrSpec[T]) GetSpec(c *Extendable[Components]) (*T, error) {
	return o.getSpec(c, make(visitedObjects))
}

func (o *RefOrSpec[T]) getSpec(c *Extendable[Components], visited visitedObjects) (*T, error) {
	// some guards
	switch {
	case o.Spec != nil:
		return o.Spec, nil
	case o.Ref == nil:
		return nil, fmt.Errorf("spect not found; all visited refs: %s", visited)
	case visited[o.Ref.Ref]:
		return nil, fmt.Errorf("cycle ref %q detected; all visited refs: %s", o.Ref.Ref, visited)
	case !strings.HasPrefix(o.Ref.Ref, "#/components/"):
		// TODO: support loading by url
		return nil, fmt.Errorf("loading outside of components is not implemented for the ref %q; all visited refs: %s", o.Ref.Ref, visited)
	case c == nil:
		return nil, fmt.Errorf("components is required, but got nil; all visited refs: %s", visited)
	}
	visited[o.Ref.Ref] = true

	parts := strings.SplitN(o.Ref.Ref[13:], "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("incorrect ref %q; all visited refs: %s", o.Ref.Ref, visited)
	}
	objName := parts[1]
	var ref any
	switch parts[0] {
	case "schemas":
		ref = c.Spec.Schemas[objName]
	case "responses":
		ref = c.Spec.Responses[objName]
	case "parameters":
		ref = c.Spec.Parameters[objName]
	case "examples":
		ref = c.Spec.Examples[objName]
	case "requestBodies":
		ref = c.Spec.RequestBodies[objName]
	case "headers":
		ref = c.Spec.Headers[objName]
	case "links":
		ref = c.Spec.Links[objName]
	case "callbacks":
		ref = c.Spec.Callbacks[objName]
	case "paths":
		ref = c.Spec.Paths[objName]
	default:
		return nil, fmt.Errorf("unexpected component %q; all visited refs: %s", ref, visited)
	}
	obj, ok := ref.(*RefOrSpec[T])
	if !ok {
		return nil, fmt.Errorf("expected spec of type %T, but got %T; all visited refs: %s", RefOrSpec[T]{}, ref, visited)
	}
	if obj.Spec != nil {
		return obj.Spec, nil
	}
	return obj.getSpec(c, visited)
}

// MarshalJSON implements json.Marshaler interface.
func (o *RefOrSpec[T]) MarshalJSON() ([]byte, error) {
	var v any
	if o.Ref != nil {
		v = o.Ref
	} else {
		v = o.Spec
	}
	data, err := json.Marshal(&v)
	if err != nil {
		return nil, fmt.Errorf("%T: %w", o.Spec, err)
	}
	return data, nil
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (o *RefOrSpec[T]) UnmarshalJSON(data []byte) error {
	if json.Unmarshal(data, &o.Ref) == nil && o.Ref.Ref != "" {
		o.Spec = nil
		return nil
	}

	o.Ref = nil
	if err := json.Unmarshal(data, &o.Spec); err != nil {
		return fmt.Errorf("%T: %w", o.Spec, err)
	}
	return nil
}

// MarshalYAML implements yaml.Marshaler interface.
func (o *RefOrSpec[T]) MarshalYAML() (any, error) {
	var v any
	if o.Ref != nil {
		v = o.Ref
	} else {
		v = o.Spec
	}
	return v, nil
}

// UnmarshalYAML implements yaml.Unmarshaler interface.
func (o *RefOrSpec[T]) UnmarshalYAML(node *yaml.Node) error {
	if node.Decode(&o.Ref) == nil && o.Ref.Ref != "" {
		return nil
	}

	o.Ref = nil
	if err := node.Decode(&o.Spec); err != nil {
		return fmt.Errorf("%T: %w", o.Spec, err)
	}
	return nil
}

func (o *RefOrSpec[T]) validateSpec(location string, validator *Validator) []*validationError {
	var errs []*validationError
	if o.Spec != nil {
		if spec, ok := any(o.Spec).(validatable); ok {
			errs = append(errs, spec.validateSpec(location, validator)...)
		} else {
			errs = append(errs, newValidationError(location, fmt.Errorf("unsupported spec type: %T", o.Spec)))
		}
	} else {
		// do not validate already visited refs
		if validator.visited[o.Ref.Ref] {
			return errs
		}
		validator.visited[o.Ref.Ref] = true
		spec, err := o.GetSpec(validator.spec.Spec.Components)
		if err != nil {
			errs = append(errs, newValidationError(location, err))
		} else if spec != nil {
			errs = append(errs, o.validateSpec(location, validator)...)
		}
	}
	return errs
}

package spec

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

// Ref is a simple object to allow referencing other components in the OpenAPI document, internally and externally.
// The $ref string value contains a URI [RFC3986], which identifies the location of the value being referenced.
// See the rules for resolving Relative References.
//
// https://spec.openapis.org/oas/v3.1.0#reference-object
//
// Example:
//   $ref: '#/components/schemas/Pet'
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

// NewRef creates an object of Ref type.
func NewRef(ref string) *Ref {
	return &Ref{
		Ref: ref,
	}
}

// RefOrSpec holds either Ref or any OpenAPI spec type.
//
// NOTE: The Ref object takes precedent over Spec if using json or yaml Marshal and Unmarshal functions.
type RefOrSpec[T openAPIConstraint] struct {
	Ref  *Ref `json:"-" yaml:"-"`
	Spec *T   `json:"-" yaml:"-"`
}

// NewRefOrSpec creates an object of RefOrSpec type for either Ref or Spec
func NewRefOrSpec[T openAPIConstraint](ref *Ref, spec *T) *RefOrSpec[T] {
	o := RefOrSpec[T]{}
	switch {
	case ref != nil:
		o.Ref = ref
	case spec != nil:
		o.Spec = spec
	}
	return &o
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

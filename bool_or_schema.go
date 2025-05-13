package openapi

import (
	"encoding/json"
)

// BoolOrSchema handles Boolean or Schema type.
//
// It MUST be used as a pointer,
// otherwise the `false` can be omitted by json or yaml encoders in case of `omitempty` tag is set.
type BoolOrSchema struct {
	Schema  *RefOrSpec[Schema] `json:"-" yaml:"-"`
	Allowed bool               `json:"-" yaml:"-"`
}

// MarshalJSON implements json.Marshaler interface.
func (o *BoolOrSchema) MarshalJSON() ([]byte, error) {
	var v any
	if o.Schema != nil {
		v = o.Schema
	} else {
		v = o.Allowed
	}
	return json.Marshal(&v)
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (o *BoolOrSchema) UnmarshalJSON(data []byte) error {
	if json.Unmarshal(data, &o.Allowed) == nil {
		o.Schema = nil
		return nil
	}
	if err := json.Unmarshal(data, &o.Schema); err != nil {
		return err
	}
	o.Allowed = true
	return nil
}

// MarshalYAML implements yaml.Marshaler interface.
func (o *BoolOrSchema) MarshalYAML() (any, error) {
	var v any
	if o.Schema != nil {
		v = o.Schema
	} else {
		v = o.Allowed
	}

	return v, nil
}

// UnmarshalYAML implements yaml.obsoleteUnmarshaler and goyaml.InterfaceUnmarshaler interfaces.
func (o *BoolOrSchema) UnmarshalYAML(unmarshal func(any) error) error {
	if unmarshal(&o.Allowed) == nil {
		o.Schema = nil
		return nil
	}
	if o.Schema == nil {
		o.Schema = &RefOrSpec[Schema]{}
	}
	if err := unmarshal(o.Schema); err != nil {
		return err
	}
	o.Allowed = true
	return nil
}

func (o *BoolOrSchema) validateSpec(path string, validator *Validator) []*validationError {
	var errs []*validationError
	if o.Schema != nil {
		errs = append(errs, o.Schema.validateSpec(path, validator)...)
	}
	return errs
}

func NewBoolOrSchema(v any) *BoolOrSchema {
	switch v := v.(type) {
	case bool:
		return &BoolOrSchema{Allowed: v}
	case *RefOrSpec[Schema]:
		return &BoolOrSchema{Schema: v}
	case *SchemaBulder:
		return &BoolOrSchema{Schema: v.Build()}
	default:
		return nil
	}
}

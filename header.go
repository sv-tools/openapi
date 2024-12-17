package openapi

// Header Object follows the structure of the Parameter Object with the some changes.
//
// https://spec.openapis.org/oas/v3.1.1#header-object
//
// Example:
//
//	description: The number of allowed requests in the current period
//	schema:
//	  type: integer
//
// All fields are copied from Parameter Object as is, except name and in fields.
type Header struct {
	// The schema defining the type used for the header.
	Schema *RefOrSpec[Schema] `json:"schema,omitempty" yaml:"schema,omitempty"`
	// A map containing the representations for the header.
	// The key is the media type and the value describes it.
	// The map MUST only contain one entry.
	Content map[string]*Extendable[MediaType] `json:"content,omitempty" yaml:"content,omitempty"`
	// A brief description of the header.
	// This could contain examples of use.
	// CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// Describes how the header value will be serialized.
	Style string `json:"style,omitempty" yaml:"style,omitempty"`
	// When this is true, header values of type array or object generate separate headers
	// for each value of the array or key-value pair of the map.
	// For other types of parameters this property has no effect.
	// When style is form, the default value is true.
	// For all other styles, the default value is false.
	Explode bool `json:"explode,omitempty" yaml:"explode,omitempty"`
	// Determines whether this header is mandatory.
	// The property MAY be included and its default value is false.
	Required bool `json:"required,omitempty" yaml:"required,omitempty"`
	// Specifies that a header is deprecated and SHOULD be transitioned out of usage.
	// Default value is false.
	Deprecated bool `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
}

func (o *Header) validateSpec(location string, validator *Validator) []*validationError {
	var errs []*validationError
	if o.Schema != nil && o.Content != nil {
		errs = append(errs, newValidationError(joinLoc(location, "schema&content"), ErrMutuallyExclusive))
	}

	if l := len(o.Content); l > 0 {
		if l != 1 {
			errs = append(errs, newValidationError(joinLoc(location, "content"), "must be only one item, but got '%d'", l))
		}
		for k, v := range o.Content {
			errs = append(errs, v.validateSpec(joinLoc(location, "content", k), validator)...)
		}
	}
	if o.Schema != nil {
		errs = append(errs, o.Schema.validateSpec(joinLoc(location, "schema"), validator)...)
	}

	switch o.Style {
	case "", StyleSimple:
	default:
		errs = append(errs, newValidationError(joinLoc(location, "style"), "invalid value, expected one of [%s], but got '%s'", StyleSimple, o.Style))
	}

	return errs
}

type HeaderBuilder struct {
	spec *RefOrSpec[Extendable[Header]]
}

func NewHeaderBuilder() *HeaderBuilder {
	return &HeaderBuilder{
		spec: NewRefOrExtSpec[Header](&Header{}),
	}
}

func (b *HeaderBuilder) Build() *RefOrSpec[Extendable[Header]] {
	return b.spec
}

func (b *HeaderBuilder) Extensions(v map[string]any) *HeaderBuilder {
	b.spec.Spec.Extensions = v
	return b
}

func (b *HeaderBuilder) AddExt(name string, value any) *HeaderBuilder {
	b.spec.Spec.AddExt(name, value)
	return b
}

func (b *HeaderBuilder) Schema(v *RefOrSpec[Schema]) *HeaderBuilder {
	b.spec.Spec.Spec.Schema = v
	return b
}

func (b *HeaderBuilder) Content(v map[string]*Extendable[MediaType]) *HeaderBuilder {
	b.spec.Spec.Spec.Content = v
	return b
}

func (b *HeaderBuilder) AddContent(name string, value *Extendable[MediaType]) *HeaderBuilder {
	if b.spec.Spec.Spec.Content == nil {
		b.spec.Spec.Spec.Content = make(map[string]*Extendable[MediaType], 1)
	}
	b.spec.Spec.Spec.Content[name] = value
	return b
}

func (b *HeaderBuilder) Description(v string) *HeaderBuilder {
	b.spec.Spec.Spec.Description = v
	return b
}

func (b *HeaderBuilder) Style(v string) *HeaderBuilder {
	b.spec.Spec.Spec.Style = v
	return b
}

func (b *HeaderBuilder) Explode(v bool) *HeaderBuilder {
	b.spec.Spec.Spec.Explode = v
	return b
}

func (b *HeaderBuilder) Required(v bool) *HeaderBuilder {
	b.spec.Spec.Spec.Required = v
	return b
}

func (b *HeaderBuilder) Deprecated(v bool) *HeaderBuilder {
	b.spec.Spec.Spec.Deprecated = v
	return b
}

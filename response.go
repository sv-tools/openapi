package openapi

// Response describes a single response from an API Operation, including design-time, static links to operations based on the response.
//
// https://spec.openapis.org/oas/v3.1.1#response-object
//
// Example:
//
//	description: A complex object array response
//	content:
//	  application/json:
//	    schema:
//	      type: array
//	      items:
//	        $ref: '#/components/schemas/VeryComplexType'
type Response struct {
	// Maps a header name to its definition.
	// [RFC7230] states header names are case insensitive.
	// If a response header is defined with the name "Content-Type", it SHALL be ignored.
	Headers map[string]*RefOrSpec[Extendable[Header]] `json:"headers,omitempty" yaml:"headers,omitempty"`
	// A map containing descriptions of potential response payloads.
	// The key is a media type or [media type range](appendix-D) and the value describes it.
	// For responses that match multiple keys, only the most specific key is applicable. e.g. text/plain overrides text/*
	Content map[string]*Extendable[MediaType] `json:"content,omitempty" yaml:"content,omitempty"`
	// A map of operations links that can be followed from the response.
	// The key of the map is a short name for the link, following the naming constraints of the names for Component Objects.
	Links map[string]*RefOrSpec[Extendable[Link]] `json:"links,omitempty" yaml:"links,omitempty"`
	// REQUIRED.
	// A description of the response.
	// CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
}

func (o *Response) validateSpec(location string, opts *specValidationOptions) []*validationError {
	errs := make([]*validationError, 0)
	if o.Description == "" {
		errs = append(errs, newValidationError(joinLoc(location, "description"), ErrRequired))
	}
	if o.Content != nil {
		for k, v := range o.Content {
			errs = append(errs, v.validateSpec(joinLoc(location, "content", k), opts)...)
		}
	}
	if o.Links != nil {
		for k, v := range o.Links {
			errs = append(errs, v.validateSpec(joinLoc(location, "links", k), opts)...)
		}
	}
	if o.Headers != nil {
		for k, v := range o.Headers {
			errs = append(errs, v.validateSpec(joinLoc(location, "headers", k), opts)...)
		}
	}
	return errs
}

type ResponseBuilder struct {
	spec *RefOrSpec[Extendable[Response]]
}

func NewResponseBuilder() *ResponseBuilder {
	return &ResponseBuilder{
		spec: NewRefOrExtSpec[Response](&Response{}),
	}
}

func (b *ResponseBuilder) Build() *RefOrSpec[Extendable[Response]] {
	return b.spec
}

func (b *ResponseBuilder) Extensions(v map[string]any) *ResponseBuilder {
	b.spec.Spec.Extensions = v
	return b
}

func (b *ResponseBuilder) AddExt(name string, value any) *ResponseBuilder {
	b.spec.Spec.AddExt(name, value)
	return b
}

func (b *ResponseBuilder) Headers(v map[string]*RefOrSpec[Extendable[Header]]) *ResponseBuilder {
	b.spec.Spec.Spec.Headers = v
	return b
}

func (b *ResponseBuilder) AddHeader(key string, value *RefOrSpec[Extendable[Header]]) *ResponseBuilder {
	if b.spec.Spec.Spec.Headers == nil {
		b.spec.Spec.Spec.Headers = make(map[string]*RefOrSpec[Extendable[Header]], 1)
	}
	b.spec.Spec.Spec.Headers[key] = value
	return b
}

func (b *ResponseBuilder) Content(v map[string]*Extendable[MediaType]) *ResponseBuilder {
	b.spec.Spec.Spec.Content = v
	return b
}

func (b *ResponseBuilder) AddContent(key string, value *Extendable[MediaType]) *ResponseBuilder {
	if b.spec.Spec.Spec.Content == nil {
		b.spec.Spec.Spec.Content = make(map[string]*Extendable[MediaType], 1)
	}
	b.spec.Spec.Spec.Content[key] = value
	return b
}

func (b *ResponseBuilder) Links(v map[string]*RefOrSpec[Extendable[Link]]) *ResponseBuilder {
	b.spec.Spec.Spec.Links = v
	return b
}

func (b *ResponseBuilder) AddLink(key string, value *RefOrSpec[Extendable[Link]]) *ResponseBuilder {
	if b.spec.Spec.Spec.Links == nil {
		b.spec.Spec.Spec.Links = make(map[string]*RefOrSpec[Extendable[Link]], 1)
	}
	b.spec.Spec.Spec.Links[key] = value
	return b
}

func (b *ResponseBuilder) Description(v string) *ResponseBuilder {
	b.spec.Spec.Spec.Description = v
	return b
}

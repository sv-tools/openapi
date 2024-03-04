package spec

// Response describes a single response from an API Operation, including design-time, static links to operations based on the response.
//
// https://spec.openapis.org/oas/v3.1.0#response-object
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

func (o *Response) validateSpec(path string, opts *validationOptions) []*validationError {
	errs := make([]*validationError, 0)
	if o.Description == "" {
		errs = append(errs, newValidationError(joinDot(path, "description"), ErrRequired))
	}
	if o.Content != nil {
		for k, v := range o.Content {
			errs = append(errs, v.validateSpec(joinArrayItem(joinDot(path, "content"), k), opts)...)
		}
	}
	if o.Links != nil {
		for k, v := range o.Links {
			errs = append(errs, v.validateSpec(joinArrayItem(joinDot(path, "links"), k), opts)...)
		}
	}
	if o.Headers != nil {
		for k, v := range o.Headers {
			errs = append(errs, v.validateSpec(joinArrayItem(joinDot(path, "headers"), k), opts)...)
		}
	}
	return errs
}

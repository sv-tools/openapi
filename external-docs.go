package openapi

// ExternalDocs allows referencing an external resource for extended documentation.
//
// https://spec.openapis.org/oas/v3.1.1#external-documentation-object
//
// Example:
//
//	description: Find more info here
//	url: https://example.com
type ExternalDocs struct {
	// A description of the target documentation.
	// CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description" yaml:"description"`
	// REQUIRED.
	// The URL for the target documentation.
	// This MUST be in the form of a URL.
	URL string `json:"url" yaml:"url"`
}

func (o *ExternalDocs) validateSpec(location string, opts *specValidationOptions) []*validationError {
	var errs []*validationError
	if o.URL == "" {
		errs = append(errs, newValidationError(joinLoc(location, "url"), ErrRequired))
	}
	if err := checkURL(o.URL); err != nil {
		errs = append(errs, newValidationError(joinLoc(location, "url"), err))
	}
	return errs
}

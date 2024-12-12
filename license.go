package openapi

// License information for the exposed API.
//
// https://spec.openapis.org/oas/v3.1.0#license-object
//
// Example:
//
//	name: Apache 2.0
//	identifier: Apache-2.0
type License struct {
	// REQUIRED.
	// The license name used for the API.
	Name string `json:"name" yaml:"name"`
	// An SPDX license expression for the API.
	// The identifier field is mutually exclusive of the url field.
	Identifier string `json:"identifier,omitempty" yaml:"identifier,omitempty"`
	// A URL to the license used for the API.
	// This MUST be in the form of a URL.
	// The url field is mutually exclusive of the identifier field.
	URL string `json:"url,omitempty" yaml:"url,omitempty"`
}

func (o *License) validateSpec(path string, opts *specValidationOptions) []*validationError {
	var errs []*validationError
	if o.Name == "" {
		errs = append(errs, newValidationError(joinDot(path, "name"), ErrRequired))
	}
	if o.Identifier != "" && o.URL != "" {
		errs = append(errs, newValidationError(joinDot(path, "identifier&url"), ErrMutuallyExclusive))
	}
	if err := checkURL(o.URL); err != nil {
		errs = append(errs, newValidationError(joinDot(path, "url"), err))
	}
	return errs
}

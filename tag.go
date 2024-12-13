package openapi

// Tag adds metadata to a single tag that is used by the Operation Object.
// It is not mandatory to have a Tag Object per tag defined in the Operation Object instances.
//
// https://spec.openapis.org/oas/v3.1.0#tag-object
//
// Example:
//
//	name: pet
//	description: Pets operations
type Tag struct {
	// Additional external documentation for this tag.
	ExternalDocs *Extendable[ExternalDocs] `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	// REQUIRED.
	// The name of the tag.
	Name string `json:"name" yaml:"name"`
	// A description for the tag.
	// CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
}

func (o *Tag) validateSpec(location string, opts *specValidationOptions) []*validationError {
	var errs []*validationError
	if o.Name == "" {
		errs = append(errs, newValidationError(joinLoc(location, "name"), ErrRequired))
	}
	if o.ExternalDocs != nil {
		errs = append(errs, o.ExternalDocs.validateSpec(joinLoc(location, "externalDocs"), opts)...)
	}
	opts.visited[joinLoc("tags", o.Name)] = true
	return errs
}

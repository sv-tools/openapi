package openapi

// Info provides metadata about the API.
// The metadata MAY be used by the clients if needed, and MAY be presented in editing or documentation generation tools for convenience.
//
// https://spec.openapis.org/oas/v3.1.1#info-object
//
// Example:
//
//	title: Sample Pet Store App
//	summary: A pet store manager.
//	description: This is a sample server for a pet store.
//	termsOfService: https://example.com/terms/
//	contact:
//	  name: API Support
//	  url: https://www.example.com/support
//	  email: support@example.com
//	license:
//	  name: Apache 2.0
//	  url: https://www.apache.org/licenses/LICENSE-2.0.html
//	version: 1.0.1
type Info struct {
	// REQUIRED.
	// The title of the API.
	Title string `json:"title" yaml:"title"`
	// A short summary of the API.
	Summary string `json:"summary,omitempty" yaml:"summary,omitempty"`
	// A description of the API.
	// CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// A URL to the Terms of Service for the API.
	// This MUST be in the form of a URL.
	TermsOfService string `json:"termsOfService,omitempty" yaml:"termsOfService,omitempty"`
	// The contact information for the exposed API.
	Contact *Extendable[Contact] `json:"contact,omitempty" yaml:"contact,omitempty"`
	// The license information for the exposed API.
	License *Extendable[License] `json:"license,omitempty" yaml:"license,omitempty"`
	// REQUIRED.
	// The version of the OpenAPI document (which is distinct from the OpenAPI Specification version or the API implementation version).
	Version string `json:"version" yaml:"version"`
}

func (o *Info) validateSpec(location string, validator *Validator) []*validationError {
	var errs []*validationError
	if o.Title == "" {
		errs = append(errs, newValidationError(joinLoc(location, "title"), ErrRequired))
	}
	if o.Version == "" {
		errs = append(errs, newValidationError(joinLoc(location, "version"), ErrRequired))
	}
	if o.Contact != nil {
		errs = append(errs, o.Contact.validateSpec(joinLoc(location, "contact"), validator)...)
	}
	if o.License != nil {
		errs = append(errs, o.License.validateSpec(joinLoc(location, "license"), validator)...)
	}
	if err := checkURL(o.TermsOfService); err != nil {
		errs = append(errs, newValidationError(joinLoc(location, "termsOfService"), err))
	}
	return errs
}

type InfoBuilder struct {
	spec *Extendable[Info]
}

func NewInfoBuilder() *InfoBuilder {
	return &InfoBuilder{
		spec: NewExtendable[Info](&Info{}),
	}
}

func (b *InfoBuilder) Build() *Extendable[Info] {
	return b.spec
}

func (b *InfoBuilder) Extensions(v map[string]any) *InfoBuilder {
	b.spec.Extensions = v
	return b
}

func (b *InfoBuilder) AddExt(name string, value any) *InfoBuilder {
	b.spec.AddExt(name, value)
	return b
}

func (b *InfoBuilder) Title(v string) *InfoBuilder {
	b.spec.Spec.Title = v
	return b
}

func (b *InfoBuilder) Summary(v string) *InfoBuilder {
	b.spec.Spec.Summary = v
	return b
}

func (b *InfoBuilder) Description(v string) *InfoBuilder {
	b.spec.Spec.Description = v
	return b
}

func (b *InfoBuilder) TermsOfService(v string) *InfoBuilder {
	b.spec.Spec.TermsOfService = v
	return b
}

func (b *InfoBuilder) Contact(v *Extendable[Contact]) *InfoBuilder {
	b.spec.Spec.Contact = v
	return b
}

func (b *InfoBuilder) License(v *Extendable[License]) *InfoBuilder {
	b.spec.Spec.License = v
	return b
}

func (b *InfoBuilder) Version(v string) *InfoBuilder {
	b.spec.Spec.Version = v
	return b
}

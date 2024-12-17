package openapi

// Example is expected to be compatible with the type schema of its associated value.
// Tooling implementations MAY choose to validate compatibility automatically, and reject the example value(s) if incompatible.
//
// https://spec.openapis.org/oas/v3.1.1#example-object
//
// Example:
//
//	requestBody:
//	  content:
//	    'application/json':
//	      schema:
//	        $ref: '#/components/schemas/Address'
//	      examples:
//	        foo:
//	          summary: A foo example
//	          value: {"foo": "bar"}
//	        bar:
//	          summary: A bar example
//	          value: {"bar": "baz"}
//	    'application/xml':
//	      examples:
//	        xmlExample:
//	          summary: This is an example in XML
//	          externalValue: 'https://example.org/examples/address-example.xml'
//	    'text/plain':
//	      examples:
//	        textExample:
//	          summary: This is a text example
//	          externalValue: 'https://foo.bar/examples/address-example.txt'
type Example struct {
	// Short description for the example.
	Summary string `json:"summary,omitempty" yaml:"summary,omitempty"`
	// Long description for the example.
	// CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// Embedded literal example.
	// The value field and externalValue field are mutually exclusive.
	// To represent examples of media types that cannot naturally represented in JSON or YAML,
	// use a string value to contain the example, escaping where necessary.
	Value any `json:"value,omitempty" yaml:"value,omitempty"`
	// A URI that points to the literal example.
	// This provides the capability to reference examples that cannot easily be included in JSON or YAML documents.
	// The value field and externalValue field are mutually exclusive.
	// See the rules for resolving Relative References.
	ExternalValue string `json:"externalValue,omitempty" yaml:"externalValue,omitempty"`
}

func (o *Example) validateSpec(location string, validator *Validator) []*validationError {
	var errs []*validationError
	if o.Value != nil && o.ExternalValue != "" {
		errs = append(errs, newValidationError(joinLoc(location, "value&externalValue"), ErrMutuallyExclusive))
	}
	if err := checkURL(o.ExternalValue); err != nil {
		errs = append(errs, newValidationError(joinLoc(location, "externalValue"), err))
	}
	// no validation of Value field, because it needs a schema and
	// should be validated in the object that defines the example and a schema
	return errs
}

type ExampleBuilder struct {
	spec *RefOrSpec[Extendable[Example]]
}

func NewExampleBuilder() *ExampleBuilder {
	return &ExampleBuilder{
		spec: NewRefOrExtSpec[Example](&Example{}),
	}
}

func (b *ExampleBuilder) Build() *RefOrSpec[Extendable[Example]] {
	return b.spec
}

func (b *ExampleBuilder) Extensions(v map[string]any) *ExampleBuilder {
	b.spec.Spec.Extensions = v
	return b
}

func (b *ExampleBuilder) AddExt(name string, value any) *ExampleBuilder {
	b.spec.Spec.AddExt(name, value)
	return b
}

func (b *ExampleBuilder) Summary(v string) *ExampleBuilder {
	b.spec.Spec.Spec.Summary = v
	return b
}

func (b *ExampleBuilder) Description(v string) *ExampleBuilder {
	b.spec.Spec.Spec.Description = v
	return b
}

func (b *ExampleBuilder) Value(v any) *ExampleBuilder {
	b.spec.Spec.Spec.Value = v
	return b
}

func (b *ExampleBuilder) ExternalValue(v string) *ExampleBuilder {
	b.spec.Spec.Spec.ExternalValue = v
	return b
}

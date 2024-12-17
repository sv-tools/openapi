package openapi

// MediaType provides schema and examples for the media type identified by its key.
//
// https://spec.openapis.org/oas/v3.1.1#media-type-object
//
// Example:
//
//	application/json:
//	  schema:
//	    $ref: "#/components/schemas/Pet"
//	  examples:
//	    cat:
//	      summary: An example of a cat
//	      value:
//	        name: Fluffy
//	        petType: Cat
//	        color: White
//	        gender: male
//	        breed: Persian
//	    dog:
//	      summary: An example of a dog with a cat's name
//	      value:
//	        name: Puma
//	        petType: Dog
//	        color: Black
//	        gender: Female
//	        breed: Mixed
//	    frog:
//	      $ref: "#/components/examples/frog-example"
type MediaType struct {
	// The schema defining the content of the request, response, or parameter.
	Schema *RefOrSpec[Schema] `json:"schema,omitempty" yaml:"schema,omitempty"`
	// Example of the media type. The example object SHOULD be in the correct format as specified by the media type.
	// The example field is mutually exclusive of the examples field.
	// Furthermore, if referencing a schema which contains an example, the example value SHALL override the example provided by the schema.
	Example any `json:"example,omitempty" yaml:"example,omitempty"`
	// Examples of the parameter’s potential value.
	// Each example SHOULD contain a value in the correct format as specified in the parameter encoding.
	// The examples field is mutually exclusive of the example field.
	// Furthermore, if referencing a schema that contains an example, the examples value SHALL override the example provided by the schema.
	Examples map[string]*RefOrSpec[Extendable[Example]] `json:"examples,omitempty" yaml:"examples,omitempty"`
	// A map between a property name and its encoding information.
	// The key, being the property name, MUST exist in the schema as a property.
	// The encoding object SHALL only apply to requestBody objects when the media type is multipart or application/x-www-form-urlencoded.
	Encoding map[string]*Extendable[Encoding] `json:"encoding,omitempty" yaml:"encoding,omitempty"`
}

func (o *MediaType) validateSpec(location string, opts *specValidationOptions) []*validationError {
	var errs []*validationError
	if o.Schema != nil {
		errs = append(errs, o.Schema.validateSpec(joinLoc(location, "schema"), opts)...)
	}
	if len(o.Encoding) > 0 {
		for k, v := range o.Encoding {
			errs = append(errs, v.validateSpec(joinLoc(location, "encoding", k), opts)...)
		}
	}
	if o.Example != nil && len(o.Examples) > 0 {
		errs = append(errs, newValidationError(joinLoc(location, "example&examples"), ErrMutuallyExclusive))
	}
	if len(o.Examples) > 0 {
		for k, v := range o.Examples {
			errs = append(errs, v.validateSpec(joinLoc(location, "examples", k), opts)...)
		}
	}

	if opts.doNotValidateExamples {
		return errs
	}
	if o.Schema == nil {
		return append(errs, newValidationError(location, "unable to validate examples without schema"))
	}
	if o.Example != nil {
		if e := opts.validator.ValidateDataAsJSON(location, o.Example); e != nil {
			errs = append(errs, newValidationError(joinLoc(location, "example"), e))
		}
	}
	if len(o.Examples) > 0 {
		for k, v := range o.Examples {
			example, err := v.GetSpec(opts.validator.spec.Spec.Components)
			if err != nil {
				// do not add the error, because it is already validated earlier
				continue
			}
			if value := example.Spec.Value; value != nil {
				if e := opts.validator.ValidateDataAsJSON(location, value); e != nil {
					errs = append(errs, newValidationError(joinLoc(location, "examples", k), e))
				}
			}
		}
	}

	return errs
}

type MediaTypeBuilder struct {
	spec *Extendable[MediaType]
}

func NewMediaTypeBuilder() *MediaTypeBuilder {
	return &MediaTypeBuilder{
		spec: NewExtendable[MediaType](&MediaType{}),
	}
}

func (b *MediaTypeBuilder) Build() *Extendable[MediaType] {
	return b.spec
}

func (b *MediaTypeBuilder) Extensions(v map[string]any) *MediaTypeBuilder {
	b.spec.Extensions = v
	return b
}

func (b *MediaTypeBuilder) AddExt(name string, value any) *MediaTypeBuilder {
	b.spec.AddExt(name, value)
	return b
}

func (b *MediaTypeBuilder) Schema(v *RefOrSpec[Schema]) *MediaTypeBuilder {
	b.spec.Spec.Schema = v
	return b
}

func (b *MediaTypeBuilder) Example(v any) *MediaTypeBuilder {
	b.spec.Spec.Example = v
	return b
}

func (b *MediaTypeBuilder) Examples(v map[string]*RefOrSpec[Extendable[Example]]) *MediaTypeBuilder {
	b.spec.Spec.Examples = v
	return b
}

func (b *MediaTypeBuilder) AddExample(name string, value *RefOrSpec[Extendable[Example]]) *MediaTypeBuilder {
	if b.spec.Spec.Examples == nil {
		b.spec.Spec.Examples = make(map[string]*RefOrSpec[Extendable[Example]], 1)
	}
	b.spec.Spec.Examples[name] = value
	return b
}

func (b *MediaTypeBuilder) Encoding(v map[string]*Extendable[Encoding]) *MediaTypeBuilder {
	b.spec.Spec.Encoding = v
	return b
}

func (b *MediaTypeBuilder) AddEncoding(name string, value *Extendable[Encoding]) *MediaTypeBuilder {
	if b.spec.Spec.Encoding == nil {
		b.spec.Spec.Encoding = make(map[string]*Extendable[Encoding], 1)
	}
	b.spec.Spec.Encoding[name] = value
	return b
}

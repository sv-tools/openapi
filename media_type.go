package openapi

// MediaType provides schema and examples for the media type identified by its key.
//
// https://spec.openapis.org/oas/v3.1.0#media-type-object
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

func (o *MediaType) validateSpec(path string, opts *specValidationOptions) []*validationError {
	var errs []*validationError
	if o.Schema != nil {
		errs = append(errs, o.Schema.validateSpec(joinDot(path, "schema"), opts)...)
	}
	if len(o.Encoding) > 0 {
		for k, v := range o.Encoding {
			errs = append(errs, v.validateSpec(joinArrayItem(joinDot(path, "encoding"), k), opts)...)
		}
	}
	if o.Example != nil && len(o.Examples) > 0 {
		errs = append(errs, newValidationError(joinDot(path, "example&examples"), ErrMutuallyExclusive))
	}
	if len(o.Examples) > 0 {
		for k, v := range o.Examples {
			errs = append(errs, v.validateSpec(joinArrayItem(joinDot(path, "examples"), k), opts)...)
		}
	}

	if opts.doNotValidateExamples {
		return errs
	}
	if o.Schema == nil {
		return append(errs, newValidationError(path, "unable to validate examples without schema"))
	}
	spec, err := o.Schema.GetSpec(opts.spec.Spec.Components)
	if err != nil {
		// do not add the error, because it is already validated earlier
		return errs
	}
	if o.Example != nil {
		if e := ValidateData(o.Example, spec, opts.spec); e != nil {
			errs = append(errs, newValidationError(joinDot(path, "example"), e))
		}
	}
	if len(o.Examples) > 0 {
		for k, v := range o.Examples {
			example, err := v.GetSpec(opts.spec.Spec.Components)
			if err != nil {
				// do not add the error, because it is already validated earlier
				continue
			}
			if value := example.Spec.Value; value != nil {
				if e := ValidateData(o.Example, spec, opts.spec); e != nil {
					errs = append(errs, newValidationError(joinArrayItem(joinDot(path, "examples"), k), e))
				}
			}
		}
	}

	return errs
}

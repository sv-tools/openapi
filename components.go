package openapi

// Components holds a set of reusable objects for different aspects of the OAS.
// All objects defined within the components object will have no effect on the API unless they are explicitly referenced
// from properties outside the components object.
//
// https://spec.openapis.org/oas/v3.1.1#components-object
//
// Example:
//
//	components:
//	  schemas:
//	    GeneralError:
//	      type: object
//	      properties:
//	        code:
//	          type: integer
//	          format: int32
//	        message:
//	          type: string
//	    Category:
//	      type: object
//	      properties:
//	        id:
//	          type: integer
//	          format: int64
//	        name:
//	          type: string
//	    Tag:
//	      type: object
//	      properties:
//	        id:
//	          type: integer
//	          format: int64
//	        name:
//	          type: string
//	  parameters:
//	    skipParam:
//	      name: skip
//	      in: query
//	      description: number of items to skip
//	      required: true
//	      schema:
//	        type: integer
//	        format: int32
//	    limitParam:
//	      name: limit
//	      in: query
//	      description: max records to return
//	      required: true
//	      schema:
//	        type: integer
//	        format: int32
//	  responses:
//	    NotFound:
//	      description: Entity not found.
//	    IllegalInput:
//	      description: Illegal input for operation.
//	    GeneralError:
//	      description: General Error
//	      content:
//	        application/json:
//	          schema:
//	            $ref: '#/components/schemas/GeneralError'
//	  securitySchemes:
//	    api_key:
//	      type: apiKey
//	      name: api_key
//	      in: header
//	    petstore_auth:
//	      type: oauth2
//	      flows:
//	        implicit:
//	          authorizationUrl: https://example.org/api/oauth/dialog
//	          scopes:
//	            write:pets: modify pets in your account
//	            read:pets: read your pets
type Components struct {
	// An object to hold reusable Schema Objects.
	Schemas map[string]*RefOrSpec[Schema] `json:"schemas,omitempty" yaml:"schemas,omitempty"`
	// An object to hold reusable Response Objects.
	Responses map[string]*RefOrSpec[Extendable[Response]] `json:"responses,omitempty" yaml:"responses,omitempty"`
	// An object to hold reusable Parameter Objects.
	Parameters map[string]*RefOrSpec[Extendable[Parameter]] `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	// An object to hold reusable Example Objects.
	Examples map[string]*RefOrSpec[Extendable[Example]] `json:"examples,omitempty" yaml:"examples,omitempty"`
	// An object to hold reusable Request Body Objects.
	RequestBodies map[string]*RefOrSpec[Extendable[RequestBody]] `json:"requestBodies,omitempty" yaml:"requestBodies,omitempty"`
	// An object to hold reusable Header Objects.
	Headers map[string]*RefOrSpec[Extendable[Header]] `json:"headers,omitempty" yaml:"headers,omitempty"`
	// An object to hold reusable Security Scheme Objects.
	SecuritySchemes map[string]*RefOrSpec[Extendable[SecurityScheme]] `json:"securitySchemes,omitempty" yaml:"securitySchemes,omitempty"`
	// An object to hold reusable Link Objects.
	Links map[string]*RefOrSpec[Extendable[Link]] `json:"links,omitempty" yaml:"links,omitempty"`
	// An object to hold reusable Callback Objects.
	Callbacks map[string]*RefOrSpec[Extendable[Callback]] `json:"callbacks,omitempty" yaml:"callbacks,omitempty"`
	// An object to hold reusable Path Item Object.
	Paths map[string]*RefOrSpec[Extendable[PathItem]] `json:"paths,omitempty" yaml:"paths,omitempty"`
}

// Add adds the given object to the appropriate list based on a type and returns the current object (self|this).
func (o *Components) Add(name string, v any) *Components {
	if v == nil {
		return o
	}
	switch spec := v.(type) {
	case *RefOrSpec[Schema]:
		if o.Schemas == nil {
			o.Schemas = make(map[string]*RefOrSpec[Schema], 1)
		}
		o.Schemas[name] = spec
	case *RefOrSpec[Extendable[Response]]:
		if o.Responses == nil {
			o.Responses = make(map[string]*RefOrSpec[Extendable[Response]], 1)
		}
		o.Responses[name] = spec
	case *RefOrSpec[Extendable[Parameter]]:
		if o.Parameters == nil {
			o.Parameters = make(map[string]*RefOrSpec[Extendable[Parameter]], 1)
		}
		o.Parameters[name] = spec
	case *RefOrSpec[Extendable[Example]]:
		if o.Examples == nil {
			o.Examples = make(map[string]*RefOrSpec[Extendable[Example]], 1)
		}
		o.Examples[name] = spec
	case *RefOrSpec[Extendable[RequestBody]]:
		if o.RequestBodies == nil {
			o.RequestBodies = make(map[string]*RefOrSpec[Extendable[RequestBody]], 1)
		}
		o.RequestBodies[name] = spec
	case *RefOrSpec[Extendable[Header]]:
		if o.Headers == nil {
			o.Headers = make(map[string]*RefOrSpec[Extendable[Header]], 1)
		}
		o.Headers[name] = spec
	case *RefOrSpec[Extendable[SecurityScheme]]:
		if o.SecuritySchemes == nil {
			o.SecuritySchemes = make(map[string]*RefOrSpec[Extendable[SecurityScheme]], 1)
		}
		o.SecuritySchemes[name] = spec
	case *RefOrSpec[Extendable[Link]]:
		if o.Links == nil {
			o.Links = make(map[string]*RefOrSpec[Extendable[Link]], 1)
		}
		o.Links[name] = spec
	case *RefOrSpec[Extendable[Callback]]:
		if o.Callbacks == nil {
			o.Callbacks = make(map[string]*RefOrSpec[Extendable[Callback]], 1)
		}
		o.Callbacks[name] = spec
	case *RefOrSpec[Extendable[PathItem]]:
		if o.Paths == nil {
			o.Paths = make(map[string]*RefOrSpec[Extendable[PathItem]], 1)
		}
		o.Paths[name] = spec
	default:
		// ignore to avoid panic
	}
	return o
}

func (o *Components) validateSpec(location string, validator *Validator) []*validationError {
	var errs []*validationError
	if o.Schemas != nil {
		for k, v := range o.Schemas {
			errs = append(errs, v.validateSpec(joinLoc(location, "schemas", k), validator)...)
		}
	}
	if o.Responses != nil {
		for k, v := range o.Responses {
			errs = append(errs, v.validateSpec(joinLoc(location, "responses", k), validator)...)
		}
	}
	if o.Parameters != nil {
		for k, v := range o.Parameters {
			errs = append(errs, v.validateSpec(joinLoc(location, "parameters", k), validator)...)
		}
	}
	if o.Examples != nil {
		for k, v := range o.Examples {
			errs = append(errs, v.validateSpec(joinLoc(location, "examples", k), validator)...)
		}
	}
	if o.RequestBodies != nil {
		for k, v := range o.RequestBodies {
			errs = append(errs, v.validateSpec(joinLoc(location, "requestBodies", k), validator)...)
		}
	}
	if o.Headers != nil {
		for k, v := range o.Headers {
			errs = append(errs, v.validateSpec(joinLoc(location, "headers", k), validator)...)
		}
	}
	if o.SecuritySchemes != nil {
		for k, v := range o.SecuritySchemes {
			errs = append(errs, v.validateSpec(joinLoc(location, "securitySchemes", k), validator)...)
		}
	}
	if o.Links != nil {
		for k, v := range o.Links {
			errs = append(errs, v.validateSpec(joinLoc(location, "links", k), validator)...)
		}
	}
	if o.Callbacks != nil {
		for k, v := range o.Callbacks {
			errs = append(errs, v.validateSpec(joinLoc(location, "callbacks", k), validator)...)
		}
	}
	if o.Paths != nil {
		for k, v := range o.Paths {
			errs = append(errs, v.validateSpec(joinLoc(location, "paths", k), validator)...)
		}
	}

	return errs
}

func NewComponents() *Extendable[Components] {
	return NewExtendable[Components](&Components{})
}

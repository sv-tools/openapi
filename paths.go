package openapi

import (
	"encoding/json"
	"strings"

	"gopkg.in/yaml.v3"
)

// Paths holds the relative paths to the individual endpoints and their operations.
// The path is appended to the URL from the Server Object in order to construct the full URL.
// The Paths MAY be empty, due to Access Control List (ACL) constraints.
//
// https://spec.openapis.org/oas/v3.1.0#paths-object
//
// Example:
//
//	/pets:
//	  get:
//	    description: Returns all pets from the system that the user has access to
//	    responses:
//	      '200':
//	        description: A list of pets.
//	        content:
//	          application/json:
//	            schema:
//	              type: array
//	              items:
//	                $ref: '#/components/schemas/pet'
type Paths struct {
	// A relative path to an individual endpoint.
	// The field name MUST begin with a forward slash (/).
	// The path is appended (no relative URL resolution) to the expanded URL
	// from the Server Object’s url field in order to construct the full URL.
	// Path templating is allowed.
	// When matching URLs, concrete (non-templated) paths would be matched before their templated counterparts.
	// Templated paths with the same hierarchy but different templated names MUST NOT exist as they are identical.
	// In case of ambiguous matching, it’s up to the tooling to decide which one to use.
	Paths map[string]*RefOrSpec[Extendable[PathItem]] `json:"-" yaml:"-"`
}

// MarshalJSON implements json.Marshaler interface.
func (o *Paths) MarshalJSON() ([]byte, error) {
	return json.Marshal(&o.Paths)
}

// UnmarshalYAML implements yaml.Unmarshaler interface.
func (o *Paths) UnmarshalYAML(node *yaml.Node) error {
	return node.Decode(&o.Paths)
}

// MarshalYAML implements yaml.Marshaler interface.
func (o *Paths) MarshalYAML() (any, error) {
	return o.Paths, nil
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (o *Paths) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &o.Paths)
}

func (o *Paths) validateSpec(loc string, opts *specValidationOptions) []*validationError {
	var errs []*validationError
	for k, v := range o.Paths {
		if !strings.HasPrefix(k, "/") {
			errs = append(errs, newValidationError(joinLoc(loc, k), "path must start with a forward slash (`/`)"))
		}
		if v == nil {
			errs = append(errs, newValidationError(joinLoc(loc, k), "path item cannot be empty"))
		} else {
			errs = append(errs, v.validateSpec(joinLoc(loc, k), opts)...)
		}
	}
	return errs
}

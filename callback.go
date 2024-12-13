package openapi

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
)

// Callback is a map of possible out-of band callbacks related to the parent operation.
// Each value in the map is a Path Item Object that describes a set of requests that may be initiated by
// the API provider and the expected responses.
// The key value used to identify the path item object is an expression, evaluated at runtime,
// that identifies a URL to use for the callback operation.
// To describe incoming requests from the API provider independent from another API call, use the webhooks field.
//
// https://spec.openapis.org/oas/v3.1.1#callback-object
//
// Example:
//
//	myCallback:
//	  '{$request.query.queryUrl}':
//	    post:
//	      requestBody:
//	        description: Callback payload
//	        content:
//	          'application/json':
//	            schema:
//	              $ref: '#/components/schemas/SomePayload'
//	      responses:
//	        '200':
//	          description: callback successfully processed
type Callback struct {
	Callback map[string]*RefOrSpec[Extendable[PathItem]]
}

// MarshalJSON implements json.Marshaler interface.
func (o *Callback) MarshalJSON() ([]byte, error) {
	return json.Marshal(&o.Callback)
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (o *Callback) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &o.Callback)
}

// MarshalYAML implements yaml.Marshaler interface.
func (o *Callback) MarshalYAML() (any, error) {
	return o.Callback, nil
}

// UnmarshalYAML implements yaml.Unmarshaler interface.
func (o *Callback) UnmarshalYAML(node *yaml.Node) error {
	return node.Decode(&o.Callback)
}

func (o *Callback) validateSpec(location string, opts *specValidationOptions) []*validationError {
	var errs []*validationError
	for k, v := range o.Callback {
		errs = append(errs, v.validateSpec(joinLoc(location, k), opts)...)
	}
	return nil
}

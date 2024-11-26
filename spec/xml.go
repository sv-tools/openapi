package spec

import "errors"

// XML is a metadata object that allows for more fine-tuned XML model definitions.
// When using arrays, XML element names are not inferred (for singular/plural forms) and the name property SHOULD
// be used to add that information.
// See examples for expected behavior.
//
// https://spec.openapis.org/oas/v3.1.0#xml-object
//
// Example:
//
//	Person:
//	  type: object
//	  properties:
//	    id:
//	      type: integer
//	      format: int32
//	      xml:
//	        attribute: true
//	    name:
//	      type: string
//	      xml:
//	        namespace: https://example.com/schema/sample
//	        prefix: sample
//
//	<Person id="123">
//	    <sample:name xmlns:sample="https://example.com/schema/sample">example</sample:name>
//	</Person>
type XML struct {
	// Replaces the name of the element/attribute used for the described schema property.
	// When defined within items, it will affect the name of the individual XML elements within the list.
	// When defined alongside type being array (outside the items), it will affect the wrapping element and only if wrapped is true.
	// If wrapped is false, it will be ignored.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// The URI of the namespace definition.
	// This MUST be in the form of an absolute URI.
	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	// The prefix to be used for the name.
	Prefix string `json:"prefix,omitempty" yaml:"prefix,omitempty"`
	// Declares whether the property definition translates to an attribute instead of an element. Default value is false.
	Attribute bool `json:"attribute,omitempty" yaml:"attribute,omitempty"`
	// MAY be used only for an array definition.
	// Signifies whether the array is wrapped (for example, <books><book/><book/></books>) or unwrapped (<book/><book/>).
	// Default value is false.
	// The definition takes effect only when defined alongside type being array (outside the items).
	Wrapped bool `json:"wrapped,omitempty" yaml:"wrapped,omitempty"`
}

// NewXML creates XML object.
func NewXML(spec ...*XML) *Extendable[XML] {
	switch len(spec) {
	case 0:
		return NewExtendable(&XML{})
	case 1:
		return NewExtendable(spec[0])
	default:
		panic(errors.New("NewXML is called with more than one spec"))
	}
}

func (o *XML) WithName(name string) *XML {
	o.Name = name
	return o
}

func (o *XML) WithNamespace(namespace string) *XML {
	o.Namespace = namespace
	return o
}

func (o *XML) WithPrefix(prefix string) *XML {
	o.Prefix = prefix
	return o
}

func (o *XML) WithAttribute(attribute bool) *XML {
	o.Attribute = attribute
	return o
}

func (o *XML) WithWrapped(wrapped bool) *XML {
	o.Wrapped = wrapped
	return o
}

func (o *XML) validateSpec(path string, opts *validationOptions) []*validationError {
	return nil // nothing to validate
}

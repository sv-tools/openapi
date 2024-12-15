package openapi

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	// Draft202012 is the default value for the schema field.
	Draft202012 = "https://json-schema.org/draft/2020-12/schema"
)

// The Schema Object allows the definition of input and output data types.
// These types can be objects, but also primitives and arrays.
// This object is a superset of the JSON Schema Specification Draft 2020-12.
// For more information about the properties, see JSON Schema Core and JSON Schema Validation.
// Unless stated otherwise, the property definitions follow those of JSON Schema and do not add any additional semantics.
// Where JSON Schema indicates that behavior is defined by the application (e.g. for annotations),
// OAS also defers the definition of semantics to the application consuming the OpenAPI document.
//
// https://spec.openapis.org/oas/v3.1.1#schema-object
// https://json-schema.org/understanding-json-schema/index.html
type Schema struct {
	// *** Core Fields ***

	// The $schema keyword is used to declare which dialect of JSON Schema the schema was written for.
	// The value of the $schema keyword is also the identifier for a schema that can be used to verify
	// that the schema is valid according to the dialect $schema identifies.
	// A schema that describes another schema is called a "meta-schema".
	// $schema applies to the entire document and must be at the root level.
	// It does not apply to externally referenced ($ref, $dynamicRef) documents.
	// Those schemas need to declare their own $schema.
	// If $schema is not used, an implementation might allow you to specify a value externally or
	// it might make assumptions about which specification version should be used to evaluate the schema.
	// It's recommended that all JSON Schemas have a $schema keyword to communicate to readers and
	// tooling which specification version is intended.
	//
	// https://json-schema.org/understanding-json-schema/reference/schema#schema
	Schema string `json:"$schema,omitempty" yaml:"$schema,omitempty"`
	// The value of $id is a URI-reference without a fragment that resolves against the retrieval-uri.
	// The resulting URI is the base URI for the schema.
	//
	// https://json-schema.org/understanding-json-schema/structuring#id
	ID string `json:"$id,omitempty" yaml:"$id,omitempty"`
	// https://json-schema.org/understanding-json-schema/structuring#dollardefs
	Defs          map[string]*RefOrSpec[Schema] `json:"$defs,omitempty" yaml:"$defs,omitempty"`
	DynamicRef    string                        `json:"$dynamicRef,omitempty" yaml:"$dynamicRef,omitempty"`
	Vocabulary    map[string]bool               `json:"$vocabulary,omitempty" yaml:"$vocabulary,omitempty"`
	DynamicAnchor string                        `json:"$dynamicAnchor,omitempty" yaml:"dynamicAnchor,omitempty"`
	// https://json-schema.org/understanding-json-schema/reference/type#type-specific-keywords
	Type *SingleOrArray[string] `json:"type,omitempty" yaml:"type,omitempty"`

	// *** Generic Fields ***
	//
	// https://json-schema.org/understanding-json-schema/reference/generic.html

	// The default keyword specifies a default value.
	// This value is not used to fill in missing values during the validation process.
	// Non-validation tools such as documentation generators or form generators may use this value
	// to give hints to users about how to use a value.
	// However, default is typically used to express that if a value is missing,
	// then the value is semantically the same as if the value was present with the default value.
	// The value of default should validate against the schema in which it resides, but that isn't required.
	Default any `json:"default,omitempty" yaml:"default,omitempty"`
	// The title keyword preferably be short.
	Title string `json:"title,omitempty" yaml:"title,omitempty"`
	// The description keyword provides a more lengthy explanation about the purpose of the data described by the schema.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// The const keyword is used to restrict a value to a single value.
	//
	// https://json-schema.org/understanding-json-schema/reference/const
	Const string `json:"const,omitempty" yaml:"const,omitempty"`
	// The $comment keyword is strictly intended for adding comments to a schema.
	// Its value must always be a string.
	// Unlike the annotations title, description, and examples, JSON schema implementations aren’t allowed
	// to attach any meaning or behavior to it whatsoever, and may even strip them at any time.
	// Therefore, they are useful for leaving notes to future editors of a JSON schema,
	// but should not be used to communicate to users of the schema.
	//
	// https://json-schema.org/understanding-json-schema/reference/generic.html#comments
	Comment string `json:"$comment,omitempty" yaml:"$comment,omitempty"`
	// The enum keyword is used to restrict a value to a fixed set of values.
	// It must be an array with at least one element, where each element is unique.
	//
	// https://json-schema.org/understanding-json-schema/reference/generic.html#enumerated-values
	Enum []any `json:"enum,omitempty" yaml:"enum,omitempty"`
	// The examples keyword is a place to provide an array of examples that validate against the schema.
	// This isn't used for validation, but may help with explaining the effect and purpose of the schema to a reader.
	// Each entry should validate against the schema in which it resides, but that isn't strictly required.
	// There is no need to duplicate the default value in the examples array,
	// since default will be treated as another example.
	Examples []any `json:"examples,omitempty" yaml:"examples,omitempty"`
	// The readOnly indicates that a value should not be modified.
	// It could be used to indicate that a PUT request that changes a value would result in a 400 Bad Request response.
	ReadOnly bool `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`
	// The writeOnly indicates that a value may be set, but will remain hidden.
	// In could be used to indicate you can set a value with a PUT request,
	// but it would not be included when retrieving that record with a GET request.
	WriteOnly bool `json:"writeOnly,omitempty" yaml:"writeOnly,omitempty"`
	// The deprecated keyword is a boolean that indicates that the instance value the keyword applies to
	// should not be used and may be removed in the future.
	Deprecated bool `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`

	// *** Media Fields ***
	// String-encoding non-JSON data.
	//
	// https://json-schema.org/understanding-json-schema/reference/non_json_data#media:-string-encoding-non-json-data

	// https://json-schema.org/understanding-json-schema/reference/non_json_data#contentschema
	ContentSchema *RefOrSpec[Schema] `json:"contentSchema,omitempty" yaml:"contentSchema,omitempty"`
	// The contentMediaType keyword specifies the MIME type of the contents of a string, as described in RFC 2046.
	// There is a list of MIME types officially registered by the IANA, but the set of types supported will be
	// application and operating system dependent.
	//
	// https://json-schema.org/understanding-json-schema/reference/non_json_data#contentmediatype
	ContentMediaType string `json:"contentMediaType,omitempty" yaml:"contentMediaType,omitempty"`
	// The contentEncoding keyword specifies the encoding used to store the contents, as specified in RFC 2054, part 6.1 and RFC 4648.
	//
	// https://json-schema.org/understanding-json-schema/reference/non_json_data#contentencoding
	ContentEncoding string `json:"contentEncoding,omitempty" yaml:"contentEncoding,omitempty"`

	// *** Composition Fields ***
	//
	// https://json-schema.org/understanding-json-schema/reference/combining.html

	// The not keyword declares that an instance validates if it doesn’t validate against the given subschema.
	//
	// https://json-schema.org/understanding-json-schema/reference/combining.html#not
	Not *RefOrSpec[Schema] `json:"not,omitempty" yaml:"not,omitempty"`
	// To validate against allOf, the given data must be valid against all of the given subschemas.
	//
	// https://json-schema.org/understanding-json-schema/reference/combining.html#allof
	AllOf []*RefOrSpec[Schema] `json:"allOf,omitempty" yaml:"allOf,omitempty"`
	// To validate against anyOf, the given data must be valid against any (one or more) of the given subschemas.
	//
	// https://json-schema.org/understanding-json-schema/reference/combining.html#anyof
	AnyOf []*RefOrSpec[Schema] `json:"anyOf,omitempty" yaml:"anyOf,omitempty"`
	// To validate against oneOf, the given data must be valid against exactly one of the given subschemas.
	//
	// https://json-schema.org/understanding-json-schema/reference/combining.html#oneof
	OneOf []*RefOrSpec[Schema] `json:"oneOf,omitempty" yaml:"oneOf,omitempty"`

	// *** Conditional Fields ***
	// Applying Subschemas Conditionally.
	//
	// https://json-schema.org/understanding-json-schema/reference/conditionals.html

	// The dependentRequired keyword conditionally requires that certain properties must be present if
	// a given property is present in an object.
	// For example, suppose we have a schema representing a customer.
	// If you have their credit card number, you also want to ensure you have a billing address.
	// If you don’t have their credit card number, a billing address would not be required.
	// We represent this dependency of one property on another using the dependentRequired keyword.
	// The value of the dependentRequired keyword is an object.
	// Each entry in the object maps from the name of a property, p, to an array of strings listing properties that
	// are required if p is present.
	//
	// https://json-schema.org/understanding-json-schema/reference/conditionals.html#dependentrequired
	DependentRequired map[string][]string `json:"dependentRequired,omitempty" yaml:"dependentRequired,omitempty"`
	// The dependentSchemas keyword conditionally applies a subschema when a given property is present.
	// This schema is applied in the same way allOf applies schemas.
	// Nothing is merged or extended.
	// Both schemas apply independently.
	//
	// https://json-schema.org/understanding-json-schema/reference/conditionals.html#dependentschemas
	DependentSchemas map[string]*RefOrSpec[Schema] `json:"dependentSchemas,omitempty" yaml:"dependentSchemas,omitempty"`

	// https://json-schema.org/understanding-json-schema/reference/conditionals.html#if-then-else
	If   *RefOrSpec[Schema] `json:"if,omitempty" yaml:"if,omitempty"`
	Then *RefOrSpec[Schema] `json:"then,omitempty" yaml:"then,omitempty"`
	Else *RefOrSpec[Schema] `json:"else,omitempty" yaml:"else,omitempty"`

	// *** Number Type Fields ***
	//
	// https://json-schema.org/understanding-json-schema/reference/numeric.html#numeric-types

	// MultipleOf restricts the numbers to a multiple of a given number, using the multipleOf keyword.
	// It may be set to any positive number.
	//
	// https://json-schema.org/understanding-json-schema/reference/numeric.html#multiples
	MultipleOf *int `json:"multipleOf,omitempty" yaml:"multipleOf,omitempty"`
	// x ≥ minimum
	Minimum *int `json:"minimum,omitempty" yaml:"minimum,omitempty"`
	// x > exclusiveMinimum
	ExclusiveMinimum *int `json:"exclusiveMinimum,omitempty" yaml:"exclusiveMinimum,omitempty"`
	// x ≤ maximum
	Maximum *int `json:"maximum,omitempty" yaml:"maximum,omitempty"`
	// x < exclusiveMaximum
	ExclusiveMaximum *int `json:"exclusiveMaximum,omitempty" yaml:"exclusiveMaximum,omitempty"`

	// *** String Type Fields ***
	//
	// https://json-schema.org/understanding-json-schema/reference/string.html#string

	MinLength *int   `json:"minLength,omitempty" yaml:"minLength,omitempty"`
	MaxLength *int   `json:"maxLength,omitempty" yaml:"maxLength,omitempty"`
	Pattern   string `json:"pattern,omitempty" yaml:"pattern,omitempty"`
	Format    string `json:"format,omitempty" yaml:"format,omitempty"`

	// ** Array Type Fields ***
	//
	// https://json-schema.org/understanding-json-schema/reference/array.html#array

	// List validation is useful for arrays of arbitrary length where each item matches the same schema.
	// For this kind of array, set the items keyword to a single schema that will be used to validate all of the items in the array.
	//
	// https://json-schema.org/understanding-json-schema/reference/array#items
	Items *BoolOrSchema `json:"items,omitempty" yaml:"items,omitempty"`
	// https://json-schema.org/understanding-json-schema/reference/array#length
	MaxItems *int `json:"maxItems,omitempty" yaml:"maxItems,omitempty"`
	// The unevaluatedItems keyword is similar to unevaluatedProperties, but for items.
	//
	// https://json-schema.org/understanding-json-schema/reference/array#unevaluateditems
	UnevaluatedItems *BoolOrSchema `json:"unevaluatedItems,omitempty" yaml:"unevaluatedItems,omitempty"`
	// While the items schema must be valid for every item in the array, the contains schema only needs
	// to validate against one or more items in the array.
	//
	// https://json-schema.org/understanding-json-schema/reference/array.html#contains
	Contains    *RefOrSpec[Schema] `json:"contains,omitempty" yaml:"contains,omitempty"`
	MinContains *int               `json:"minContains,omitempty" yaml:"minContains,omitempty"`
	MaxContains *int               `json:"maxContains,omitempty" yaml:"maxContains,omitempty"`
	// https://json-schema.org/understanding-json-schema/reference/array.html#length
	MinItems *int `json:"minItems,omitempty" yaml:"minItems,omitempty"`
	// A schema can ensure that each of the items in an array is unique.
	// Simply set the uniqueItems keyword to true.
	//
	// https://json-schema.org/understanding-json-schema/reference/array.html#uniqueness
	UniqueItems *bool `json:"uniqueItems,omitempty" yaml:"uniqueItems,omitempty"`
	// The prefixItems is an array, where each item is a schema that corresponds to each index of the document’s array.
	// That is, an array where the first element validates the first element of the input array,
	// the second element validates the second element of the input array, etc.
	//
	// https://json-schema.org/understanding-json-schema/reference/array.html#tuple-validation
	PrefixItems []*RefOrSpec[Schema] `json:"prefixItems,omitempty" yaml:"prefixItems,omitempty"`

	// ** Object Type Fields ***
	//
	// https://json-schema.org/understanding-json-schema/reference/object.html#object

	// The properties (key-value pairs) on an object are defined using the properties keyword.
	// The value of properties is an object, where each key is the name of a property and each value is
	// a schema used to validate that property.
	// Any property that doesn't match any of the property names in the properties keyword is ignored by this keyword.
	//
	// https://json-schema.org/understanding-json-schema/reference/object.html#properties
	Properties map[string]*RefOrSpec[Schema] `json:"properties,omitempty" yaml:"properties,omitempty"`
	// Sometimes you want to say that, given a particular kind of property name, the value should match a particular schema.
	// That’s where patternProperties comes in: it maps regular expressions to schemas.
	// If a property name matches the given regular expression, the property value must validate against the corresponding schema.
	//
	// https://json-schema.org/understanding-json-schema/reference/object.html#pattern-properties
	PatternProperties map[string]*RefOrSpec[Schema] `json:"patternProperties,omitempty" yaml:"patternProperties,omitempty"`
	// The additionalProperties keyword is used to control the handling of extra stuff, that is,
	// properties whose names are not listed in the properties keyword or match any of the regular expressions
	// in the patternProperties keyword.
	// By default any additional properties are allowed.
	//
	// The value of the additionalProperties keyword is a schema that will be used to validate any properties in the instance
	// that are not matched by properties or patternProperties.
	// Setting the additionalProperties schema to false means no additional properties will be allowed.
	//
	// https://json-schema.org/understanding-json-schema/reference/object.html#additional-properties
	AdditionalProperties *BoolOrSchema `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"`
	// The unevaluatedProperties keyword is similar to additionalProperties except that it can recognize properties declared in subschemas.
	// So, the example from the previous section can be rewritten without the need to redeclare properties.
	//
	// https://json-schema.org/understanding-json-schema/reference/object.html#unevaluated-properties
	UnevaluatedProperties *BoolOrSchema `json:"unevaluatedProperties,omitempty" yaml:"unevaluatedProperties,omitempty"`
	// The names of properties can be validated against a schema, irrespective of their values.
	// This can be useful if you don’t want to enforce specific properties, but you want to make sure that
	// the names of those properties follow a specific convention.
	// You might, for example, want to enforce that all names are valid ASCII tokens so they can be used
	// as attributes in a particular programming language.
	//
	// https://json-schema.org/understanding-json-schema/reference/object.html#property-names
	PropertyNames *RefOrSpec[Schema] `json:"propertyNames,omitempty" yaml:"propertyNames,omitempty"`
	// The min number of properties on an object.
	//
	// https://json-schema.org/understanding-json-schema/reference/object.html#size
	MinProperties *int `json:"minProperties,omitempty" yaml:"minProperties,omitempty"`
	// The max number of properties on an object.
	//
	// https://json-schema.org/understanding-json-schema/reference/object.html#size
	MaxProperties *int `json:"maxProperties,omitempty" yaml:"maxProperties,omitempty"`
	// The required keyword takes an array of zero or more strings.
	// Each of these strings must be unique.
	//
	// https://json-schema.org/understanding-json-schema/reference/object.html#required-properties
	Required []string `json:"required,omitempty" yaml:"required,omitempty"`

	// *** OpenAPI Fields ***

	// Adds support for polymorphism.
	// The discriminator is an object name that is used to differentiate between other schemas which may satisfy the payload description.
	// See Composition and Inheritance for more details.
	Discriminator *Discriminator `json:"discriminator,omitempty" yaml:"discriminator,omitempty"`
	// Additional external documentation for this tag.
	// xml
	XML *Extendable[XML] `json:"xml,omitempty" yaml:"xml,omitempty"`
	// Additional external documentation for this schema.
	ExternalDocs *Extendable[ExternalDocs] `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	// A free-form property to include an example of an instance for this schema.
	// To represent examples that cannot be naturally represented in JSON or YAML, a string value can be used to
	// contain the example with escaping where necessary.
	//
	// Deprecated: The example property has been deprecated in favor of the JSON Schema examples keyword.
	// Use of example is discouraged, and later versions of this specification may remove it.
	Example any `json:"example,omitempty" yaml:"example,omitempty"`

	Extensions map[string]any `json:"-" yaml:"-"`
}

// AddExt sets the extension and returns the current object (self|this).
// Schema does not require special `x-` prefix.
// The extension will be ignored if the name overlaps with a struct field during marshalling to JSON or YAML.
func (o *Schema) AddExt(name string, value any) *Schema {
	if o.Extensions == nil {
		o.Extensions = make(map[string]any, 1)
	}
	o.Extensions[name] = value
	return o
}

func (o *Schema) GetExt(name string) any {
	if o.Extensions == nil {
		return nil
	}
	if !strings.HasPrefix(name, ExtensionPrefix) {
		name = ExtensionPrefix + name
	}
	return o.Extensions[name]
}

// returns the list of public fields for given tag and ignores `-` names
func getFields(t reflect.Type, tag string) map[string]struct{} {
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil
	}
	n := t.NumField()
	ret := make(map[string]struct{})
	for i := 0; i < n; i++ {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}
		if f.Anonymous {
			sub := getFields(f.Type, tag)
			for n, v := range sub {
				ret[n] = v
			}
			continue
		}
		name, _, _ := strings.Cut(f.Tag.Get(tag), ",")
		if name == "-" {
			continue
		}
		if name == "" {
			name = f.Name
		}
		ret[name] = struct{}{}
	}
	if len(ret) == 0 {
		return nil
	}
	return ret
}

type intSchema Schema // needed to avoid recursion in marshal/unmarshal

// MarshalJSON implements json.Marshaler interface.
func (o *Schema) MarshalJSON() ([]byte, error) {
	var raw map[string]json.RawMessage
	exts, err := json.Marshal(&o.Extensions)
	if err != nil {
		return nil, fmt.Errorf("%T.Extensions: %w", o, err)
	}
	if err := json.Unmarshal(exts, &raw); err != nil {
		return nil, fmt.Errorf("%T(raw extensions): %w", o, err)
	}
	s := intSchema(*o)
	fields, err := json.Marshal(&s)
	if err != nil {
		return nil, fmt.Errorf("%T: %w", o, err)
	}
	if err := json.Unmarshal(fields, &raw); err != nil {
		return nil, fmt.Errorf("%T(raw fields): %w", o, err)
	}
	data, err := json.Marshal(&raw)
	if err != nil {
		return nil, fmt.Errorf("%T(raw): %w", o, err)
	}
	return data, nil
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (o *Schema) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("%T: %w", o, err)
	}
	exts := make(map[string]any)
	keys := getFields(reflect.TypeOf(o), "json")
	for name, value := range raw {
		if _, ok := keys[name]; !ok {
			var v any
			if err := json.Unmarshal(value, &v); err != nil {
				return fmt.Errorf("%T.Extensions.%s: %w", o, name, err)
			}
			exts[name] = v
			delete(raw, name)
		}
	}
	fields, err := json.Marshal(&raw)
	if err != nil {
		return fmt.Errorf("%T(raw): %w", o, err)
	}
	var s intSchema
	if err := json.Unmarshal(fields, &s); err != nil {
		return fmt.Errorf("%T: %w", o, err)
	}
	s.Extensions = exts
	*o = Schema(s)
	return nil
}

// MarshalYAML implements yaml.Marshaler interface.
func (o *Schema) MarshalYAML() (any, error) {
	var raw map[string]any
	exts, err := yaml.Marshal(&o.Extensions)
	if err != nil {
		return nil, fmt.Errorf("%T.Extensions: %w", o, err)
	}
	if err := yaml.Unmarshal(exts, &raw); err != nil {
		return nil, fmt.Errorf("%T(raw extensions): %w", o, err)
	}
	s := intSchema(*o)
	fields, err := yaml.Marshal(&s)
	if err != nil {
		return nil, fmt.Errorf("%T: %w", o, err)
	}
	if err := yaml.Unmarshal(fields, &raw); err != nil {
		return nil, fmt.Errorf("%T(raw fields): %w", o, err)
	}
	return raw, nil
}

// UnmarshalYAML implements yaml.Unmarshaler interface.
func (o *Schema) UnmarshalYAML(node *yaml.Node) error {
	var raw map[string]any
	if err := node.Decode(&raw); err != nil {
		return fmt.Errorf("%T: %w", o, err)
	}
	exts := make(map[string]any)
	keys := getFields(reflect.TypeOf(o), "json")
	for name, value := range raw {
		if _, ok := keys[name]; !ok {
			exts[name] = value
			delete(raw, name)
		}
	}
	fields, err := yaml.Marshal(&raw)
	if err != nil {
		return fmt.Errorf("%T(raw): %w", o, err)
	}
	var s intSchema
	if err := yaml.Unmarshal(fields, &s); err != nil {
		return fmt.Errorf("%T: %w", o, err)
	}
	s.Extensions = exts
	*o = Schema(s)
	return nil
}

func (o *Schema) validateSpec(location string, opts *specValidationOptions) []*validationError {
	var errs []*validationError

	if o.Discriminator != nil {
		errs = append(errs, o.Discriminator.validateSpec(joinLoc(location, "discriminator"), opts)...)
	}
	if o.XML != nil {
		errs = append(errs, o.XML.validateSpec(joinLoc(location, "xml"), opts)...)
	}
	if o.ExternalDocs != nil {
		errs = append(errs, o.ExternalDocs.validateSpec(joinLoc(location, "externalDocs"), opts)...)
	}
	if o.Example != nil {
		if !opts.doNotValidateExamples {
			if e := opts.validator.ValidateDataAsJSON(location, o.Example); e != nil {
				errs = append(errs, newValidationError(joinLoc(location, "example"), e))
			}
		}
	}

	// JsonSchemaComposition
	if o.Not != nil {
		errs = append(errs, o.Not.validateSpec(joinLoc(location, "not"), opts)...)
	}
	if o.AllOf != nil {
		for i, v := range o.AllOf {
			errs = append(errs, v.validateSpec(joinLoc(location, "allOf", i), opts)...)
		}
	}
	if o.AnyOf != nil {
		for i, v := range o.AnyOf {
			errs = append(errs, v.validateSpec(joinLoc(location, "anyOf", i), opts)...)
		}
	}
	if o.OneOf != nil {
		for i, v := range o.OneOf {
			errs = append(errs, v.validateSpec(joinLoc(location, "oneOf", i), opts)...)
		}
	}

	// JsonSchemaCore
	if o.Schema != "" && o.Schema != Draft202012 {
		errs = append(errs, newValidationError(joinLoc(location, "schema"), "must be '%s', but got '%s'", Draft202012, o.Schema))
	}
	if len(o.Defs) > 0 {
		for k, v := range o.Defs {
			errs = append(errs, v.validateSpec(joinLoc(location, "defs", k), opts)...)
		}
	}
	if o.Type != nil {
		switch len(*o.Type) {
		case 0: // not type or any type
		case 1:
			switch v := (*o.Type)[0]; v {
			case StringType, NumberType, IntegerType, BooleanType, ObjectType, ArrayType, NullType:
			default:
				errs = append(errs, newValidationError(joinLoc(location, "type"), "invalid value, expected one of [%s, %s, %s, %s, %s, %s, %s], but got '%s'", StringType, NumberType, IntegerType, BooleanType, ObjectType, ArrayType, NullType, v))
			}
		default:
			for i, v := range *o.Type {
				switch v {
				case StringType, NumberType, IntegerType, BooleanType, ObjectType, ArrayType, NullType:
				default:
					errs = append(errs, newValidationError(joinLoc(location, "type", i), "invalid value, expected one of [%s, %s, %s, %s, %s, %s, %s], but got '%s'", StringType, NumberType, IntegerType, BooleanType, ObjectType, ArrayType, NullType, v))
				}
			}
		}
	}

	// JsonSchemaMedia
	if o.ContentSchema != nil {
		errs = append(errs, o.ContentSchema.validateSpec(joinLoc(location, "contentSchema"), opts)...)
	}
	if o.ContentEncoding != "" {
		switch o.ContentEncoding {
		case SevenBitEncoding, EightBitEncoding, BinaryEncoding, QuotedPrintableEncoding, Base16Encoding, Base32Encoding, Base64Encoding:
		default:
			errs = append(errs, newValidationError(joinLoc(location, "contentEncoding"), "invalid value, expected one of [%s, %s, %s, %s, %s, %s, %s], but got '%s'", SevenBitEncoding, EightBitEncoding, BinaryEncoding, QuotedPrintableEncoding, Base16Encoding, Base32Encoding, Base64Encoding, o.ContentEncoding))
		}
	}

	// JsonSchemaGeneric
	if o.Default != nil {
		if !opts.doNotValidateDefaultValues {
			if e := opts.validator.ValidateDataAsJSON(location, o.Default); e != nil {
				errs = append(errs, newValidationError(joinLoc(location, "default"), e))
			}
		}
		if len(o.Enum) > 0 {
			var found bool
			for _, v := range o.Enum {
				if reflect.DeepEqual(o.Default, v) {
					found = true
					break
				}
			}
			if !found {
				errs = append(errs, newValidationError(joinLoc(location, "default"), "invalid value, expected one of enum values: %v", o.Enum))
			}
		}
	}

	if len(o.Examples) > 0 && !opts.doNotValidateExamples {
		for k, v := range o.Examples {
			if e := opts.validator.ValidateDataAsJSON(location, v); e != nil {
				errs = append(errs, newValidationError(joinLoc(location, "examples", k), e))
			}
		}
	}

	if o.Type != nil {
		for _, t := range *o.Type {
			switch t {
			case ArrayType: // JsonSchemaTypeArray
				if o.Items != nil {
					errs = append(errs, o.Items.validateSpec(joinLoc(location, "items"), opts)...)
				}
				if o.MinItems != nil && *o.MinItems < 0 {
					errs = append(errs, newValidationError(joinLoc(location, "minItems"), "must be greater than or equal to 0"))
				}
				if o.MaxItems != nil && *o.MaxItems < 0 {
					errs = append(errs, newValidationError(joinLoc(location, "maxItems"), "must be greater than or equal to 0"))
					if o.MinItems != nil && *o.MaxItems < *o.MinItems {
						errs = append(errs, newValidationError(joinLoc(location, "maxItems"), "must be greater than or equal to minItems"))
					}
				}
				if o.UnevaluatedItems != nil {
					errs = append(errs, o.UnevaluatedItems.validateSpec(joinLoc(location, "unevaluatedItems"), opts)...)
				}
				if o.Contains != nil {
					errs = append(errs, o.Contains.validateSpec(joinLoc(location, "contains"), opts)...)
				}
				if o.MinContains != nil && *o.MinContains < 0 {
					errs = append(errs, newValidationError(joinLoc(location, "minContains"), "must be greater than or equal to 0"))
				}
				if o.MaxContains != nil && *o.MaxContains < 0 {
					errs = append(errs, newValidationError(joinLoc(location, "maxContains"), "must be greater than or equal to 0"))
					if o.MinContains != nil && *o.MaxContains < *o.MinContains {
						errs = append(errs, newValidationError(joinLoc(location, "maxContains"), "must be greater than or equal to minContains"))
					}
				}
				if len(o.PrefixItems) > 0 {
					for i, v := range o.PrefixItems {
						errs = append(errs, v.validateSpec(joinLoc(location, "prefixItems", i), opts)...)
					}
				}
			case ObjectType: // JsonSchemaTypeObject
				if o.Properties != nil {
					for k, v := range o.Properties {
						errs = append(errs, v.validateSpec(joinLoc(location, "properties", k), opts)...)
					}
				}
				if o.PatternProperties != nil {
					for k, v := range o.PatternProperties {
						errs = append(errs, v.validateSpec(joinLoc(location, "patternProperties", k), opts)...)
						if _, err := regexp.Compile(k); err != nil {
							errs = append(errs, newValidationError(joinLoc(location, "patternProperties", k), err))
						}
					}
				}
				if o.AdditionalProperties != nil {
					errs = append(errs, o.AdditionalProperties.validateSpec(joinLoc(location, "additionalProperties"), opts)...)
				}
				if o.UnevaluatedItems != nil {
					errs = append(errs, o.UnevaluatedItems.validateSpec(joinLoc(location, "unevaluatedItems"), opts)...)
				}
				if o.PropertyNames != nil {
					errs = append(errs, o.PropertyNames.validateSpec(joinLoc(location, "propertyNames"), opts)...)
				}
				if o.MinProperties != nil && *o.MinProperties < 0 {
					errs = append(errs, newValidationError(joinLoc(location, "minProperties"), "must be greater than or equal to 0"))
				}
				if o.MaxProperties != nil && *o.MaxProperties < 0 {
					errs = append(errs, newValidationError(joinLoc(location, "maxProperties"), "must be greater than or equal to 0"))
					if o.MinProperties != nil && *o.MaxProperties < *o.MinProperties {
						errs = append(errs, newValidationError(joinLoc(location, "maxProperties"), "must be greater than or equal to minProperties"))
					}
				}
				if len(o.Required) > 0 {
					for i, v := range o.Required {
						if _, ok := o.Properties[v]; !ok {
							errs = append(errs, newValidationError(joinLoc(location, "required", i), "must be a property in properties"))
						}
					}
				}
			case NumberType, IntegerType: // JsonSchemaTypeNumber
				if o.MultipleOf != nil && *o.MultipleOf <= 0 {
					errs = append(errs, newValidationError(joinLoc(location, "multipleOf"), "must be greater than 0"))
				}
				if o.Minimum != nil && *o.Minimum < 0 {
					errs = append(errs, newValidationError(joinLoc(location, "minimum"), "must be greater than or equal to 0"))
				}
				if o.Maximum != nil && *o.Maximum < 0 {
					errs = append(errs, newValidationError(joinLoc(location, "maximum"), "must be greater than or equal to 0"))
					if o.Minimum != nil && *o.Maximum < *o.Minimum {
						errs = append(errs, newValidationError(joinLoc(location, "maximum"), "must be greater than or equal to minimum"))
					}
				}
				if o.ExclusiveMinimum != nil && *o.ExclusiveMinimum < 0 {
					errs = append(errs, newValidationError(joinLoc(location, "exclusiveMinimum"), "must be greater than or equal to 0"))
				}
				if o.ExclusiveMaximum != nil && *o.ExclusiveMaximum < 0 {
					errs = append(errs, newValidationError(joinLoc(location, "exclusiveMaximum"), "must be greater than or equal to 0"))
					if o.ExclusiveMinimum != nil && *o.ExclusiveMaximum < *o.ExclusiveMinimum {
						errs = append(errs, newValidationError(joinLoc(location, "exclusiveMaximum"), "must be greater than or equal to exclusiveMinimum"))
					}
				}
				if o.Minimum != nil && o.ExclusiveMinimum != nil {
					errs = append(errs, newValidationError(joinLoc(location, "minimum&exclusiveMinimum"), ErrMutuallyExclusive))
				}
				if o.Maximum != nil && o.ExclusiveMaximum != nil {
					errs = append(errs, newValidationError(joinLoc(location, "maximum&exclusiveMaximum"), ErrMutuallyExclusive))
				}
			case StringType: // JsonSchemaTypeString
				if o.MinLength != nil && *o.MinLength < 0 {
					errs = append(errs, newValidationError(joinLoc(location, "minLength"), "must be greater than or equal to 0"))
				}
				if o.MaxLength != nil && *o.MaxLength < 0 {
					errs = append(errs, newValidationError(joinLoc(location, "maxLength"), "must be greater than or equal to 0"))
					if o.MinLength != nil && *o.MaxLength < *o.MinLength {
						errs = append(errs, newValidationError(joinLoc(location, "maxLength"), "must be greater than or equal to minLength"))
					}
				}
				if o.Pattern != "" {
					if _, err := regexp.Compile(o.Pattern); err != nil {
						errs = append(errs, newValidationError(joinLoc(location, "pattern"), err))
					}
				}
			}
		}
	}
	return errs
}

type SchemaBulder struct {
	spec *RefOrSpec[Schema]
}

func NewSchemaBuilder() *SchemaBulder {
	return &SchemaBulder{
		spec: NewRefOrSpec[Schema](&Schema{
			Schema: Draft202012,
		}),
	}
}

func (b *SchemaBulder) Build() *RefOrSpec[Schema] {
	return b.spec
}

func (b *SchemaBulder) Extensions(v map[string]any) *SchemaBulder {
	b.spec.Spec.Extensions = v
	return b
}

func (b *SchemaBulder) AddExt(name string, value any) *SchemaBulder {
	b.spec.Spec.AddExt(name, value)
	return b
}

func (b *SchemaBulder) Schema(v string) *SchemaBulder {
	b.spec.Spec.Schema = v
	return b
}

func (b *SchemaBulder) ID(v string) *SchemaBulder {
	b.spec.Spec.ID = v
	return b
}

func (b *SchemaBulder) Defs(v map[string]*RefOrSpec[Schema]) *SchemaBulder {
	b.spec.Spec.Defs = v
	return b
}

func (b *SchemaBulder) AddDef(name string, value *RefOrSpec[Schema]) *SchemaBulder {
	if b.spec.Spec.Defs == nil {
		b.spec.Spec.Defs = make(map[string]*RefOrSpec[Schema], 1)
	}
	b.spec.Spec.Defs[name] = value
	return b
}

func (b *SchemaBulder) DynamicRef(v string) *SchemaBulder {
	b.spec.Spec.DynamicRef = v
	return b
}

func (b *SchemaBulder) Vocabulary(v map[string]bool) *SchemaBulder {
	b.spec.Spec.Vocabulary = v
	return b
}

func (b *SchemaBulder) AddVocabulary(name string, value bool) *SchemaBulder {
	if b.spec.Spec.Vocabulary == nil {
		b.spec.Spec.Vocabulary = make(map[string]bool, 1)
	}
	b.spec.Spec.Vocabulary[name] = value
	return b
}

func (b *SchemaBulder) DynamicAnchor(v string) *SchemaBulder {
	b.spec.Spec.DynamicAnchor = v
	return b
}

func (b *SchemaBulder) Type(v ...string) *SchemaBulder {
	b.spec.Spec.Type = NewSingleOrArray[string](v...)
	return b
}

func (b *SchemaBulder) AddType(v ...string) *SchemaBulder {
	if b.spec.Spec.Type == nil {
		b.spec.Spec.Type = NewSingleOrArray[string](v...)
	} else {
		b.spec.Spec.Type.Add(v...)
	}
	return b
}

func (b *SchemaBulder) Default(v any) *SchemaBulder {
	b.spec.Spec.Default = v
	return b
}

func (b *SchemaBulder) Title(v string) *SchemaBulder {
	b.spec.Spec.Title = v
	return b
}

func (b *SchemaBulder) Description(v string) *SchemaBulder {
	b.spec.Spec.Description = v
	return b
}

func (b *SchemaBulder) Const(v string) *SchemaBulder {
	b.spec.Spec.Const = v
	return b
}

func (b *SchemaBulder) Comment(v string) *SchemaBulder {
	b.spec.Spec.Comment = v
	return b
}

func (b *SchemaBulder) Enum(v ...any) *SchemaBulder {
	b.spec.Spec.Enum = v
	return b
}

func (b *SchemaBulder) AddEnum(v ...any) *SchemaBulder {
	b.spec.Spec.Enum = append(b.spec.Spec.Enum, v...)
	return b
}

func (b *SchemaBulder) Examples(v ...any) *SchemaBulder {
	b.spec.Spec.Examples = v
	return b
}

func (b *SchemaBulder) AddExamples(v ...any) *SchemaBulder {
	b.spec.Spec.Examples = append(b.spec.Spec.Examples, v...)
	return b
}

func (b *SchemaBulder) ReadOnly(v bool) *SchemaBulder {
	b.spec.Spec.ReadOnly = v
	return b
}

func (b *SchemaBulder) WriteOnly(v bool) *SchemaBulder {
	b.spec.Spec.WriteOnly = v
	return b
}

func (b *SchemaBulder) Deprecated(v bool) *SchemaBulder {
	b.spec.Spec.Deprecated = v
	return b
}

func (b *SchemaBulder) ContentSchema(v *RefOrSpec[Schema]) *SchemaBulder {
	b.spec.Spec.ContentSchema = v
	return b
}

func (b *SchemaBulder) ContentMediaType(v string) *SchemaBulder {
	b.spec.Spec.ContentMediaType = v
	return b
}

func (b *SchemaBulder) ContentEncoding(v string) *SchemaBulder {
	b.spec.Spec.ContentEncoding = v
	return b
}

func (b *SchemaBulder) Not(v *RefOrSpec[Schema]) *SchemaBulder {
	b.spec.Spec.Not = v
	return b
}

func (b *SchemaBulder) AllOf(v ...*RefOrSpec[Schema]) *SchemaBulder {
	b.spec.Spec.AllOf = v
	return b
}

func (b *SchemaBulder) AddAllOf(v ...*RefOrSpec[Schema]) *SchemaBulder {
	b.spec.Spec.AllOf = append(b.spec.Spec.AllOf, v...)
	return b
}

func (b *SchemaBulder) AnyOf(v ...*RefOrSpec[Schema]) *SchemaBulder {
	b.spec.Spec.AnyOf = v
	return b
}

func (b *SchemaBulder) AddAnyOf(v ...*RefOrSpec[Schema]) *SchemaBulder {
	b.spec.Spec.AnyOf = append(b.spec.Spec.AnyOf, v...)
	return b
}

func (b *SchemaBulder) OneOf(v ...*RefOrSpec[Schema]) *SchemaBulder {
	b.spec.Spec.OneOf = v
	return b
}

func (b *SchemaBulder) AddOneOf(v ...*RefOrSpec[Schema]) *SchemaBulder {
	b.spec.Spec.OneOf = append(b.spec.Spec.OneOf, v...)
	return b
}

func (b *SchemaBulder) DependentRequired(v map[string][]string) *SchemaBulder {
	b.spec.Spec.DependentRequired = v
	return b
}

func (b *SchemaBulder) AddDependentRequired(name string, value ...string) *SchemaBulder {
	if b.spec.Spec.DependentRequired == nil {
		b.spec.Spec.DependentRequired = make(map[string][]string, 1)
	}
	b.spec.Spec.DependentRequired[name] = append(b.spec.Spec.DependentRequired[name], value...)
	return b
}

func (b *SchemaBulder) DependentSchemas(v map[string]*RefOrSpec[Schema]) *SchemaBulder {
	b.spec.Spec.DependentSchemas = v
	return b
}

func (b *SchemaBulder) AddDependentSchema(name string, value *RefOrSpec[Schema]) *SchemaBulder {
	if b.spec.Spec.DependentSchemas == nil {
		b.spec.Spec.DependentSchemas = make(map[string]*RefOrSpec[Schema], 1)
	}
	b.spec.Spec.DependentSchemas[name] = value
	return b
}

func (b *SchemaBulder) If(v *RefOrSpec[Schema]) *SchemaBulder {
	b.spec.Spec.If = v
	return b
}

func (b *SchemaBulder) Then(v *RefOrSpec[Schema]) *SchemaBulder {
	b.spec.Spec.Then = v
	return b
}

func (b *SchemaBulder) Else(v *RefOrSpec[Schema]) *SchemaBulder {
	b.spec.Spec.Else = v
	return b
}

func (b *SchemaBulder) MultipleOf(v int) *SchemaBulder {
	b.spec.Spec.MultipleOf = &v
	return b
}

func (b *SchemaBulder) Minimum(v int) *SchemaBulder {
	b.spec.Spec.Minimum = &v
	return b
}

func (b *SchemaBulder) ExclusiveMinimum(v int) *SchemaBulder {
	b.spec.Spec.ExclusiveMinimum = &v
	return b
}

func (b *SchemaBulder) Maximum(v int) *SchemaBulder {
	b.spec.Spec.Maximum = &v
	return b
}

func (b *SchemaBulder) ExclusiveMaximum(v int) *SchemaBulder {
	b.spec.Spec.ExclusiveMaximum = &v
	return b
}

func (b *SchemaBulder) MinLength(v int) *SchemaBulder {
	b.spec.Spec.MinLength = &v
	return b
}

func (b *SchemaBulder) MaxLength(v int) *SchemaBulder {
	b.spec.Spec.MaxLength = &v
	return b
}

func (b *SchemaBulder) Pattern(v string) *SchemaBulder {
	b.spec.Spec.Pattern = v
	return b
}

func (b *SchemaBulder) Format(v string) *SchemaBulder {
	b.spec.Spec.Format = v
	return b
}

func (b *SchemaBulder) Items(v *BoolOrSchema) *SchemaBulder {
	b.spec.Spec.Items = v
	return b
}

func (b *SchemaBulder) MaxItems(v int) *SchemaBulder {
	b.spec.Spec.MaxItems = &v
	return b
}

func (b *SchemaBulder) UnevaluatedItems(v *BoolOrSchema) *SchemaBulder {
	b.spec.Spec.UnevaluatedItems = v
	return b
}

func (b *SchemaBulder) Contains(v *RefOrSpec[Schema]) *SchemaBulder {
	b.spec.Spec.Contains = v
	return b
}

func (b *SchemaBulder) MinContains(v int) *SchemaBulder {
	b.spec.Spec.MinContains = &v
	return b
}

func (b *SchemaBulder) MaxContains(v int) *SchemaBulder {
	b.spec.Spec.MaxContains = &v
	return b
}

func (b *SchemaBulder) MinItems(v int) *SchemaBulder {
	b.spec.Spec.MinItems = &v
	return b
}

func (b *SchemaBulder) UniqueItems(v bool) *SchemaBulder {
	b.spec.Spec.UniqueItems = &v
	return b
}

func (b *SchemaBulder) PrefixItems(v ...*RefOrSpec[Schema]) *SchemaBulder {
	b.spec.Spec.PrefixItems = v
	return b
}

func (b *SchemaBulder) AddPrefixItems(v ...*RefOrSpec[Schema]) *SchemaBulder {
	b.spec.Spec.PrefixItems = append(b.spec.Spec.PrefixItems, v...)
	return b
}

func (b *SchemaBulder) Properties(v map[string]*RefOrSpec[Schema]) *SchemaBulder {
	b.spec.Spec.Properties = v
	return b
}

func (b *SchemaBulder) AddProperty(name string, value *RefOrSpec[Schema]) *SchemaBulder {
	if b.spec.Spec.Properties == nil {
		b.spec.Spec.Properties = make(map[string]*RefOrSpec[Schema], 1)
	}
	b.spec.Spec.Properties[name] = value
	return b
}

func (b *SchemaBulder) PatternProperties(v map[string]*RefOrSpec[Schema]) *SchemaBulder {
	b.spec.Spec.PatternProperties = v
	return b
}

func (b *SchemaBulder) AddPatternProperty(name string, value *RefOrSpec[Schema]) *SchemaBulder {
	if b.spec.Spec.PatternProperties == nil {
		b.spec.Spec.PatternProperties = make(map[string]*RefOrSpec[Schema], 1)
	}
	b.spec.Spec.PatternProperties[name] = value
	return b
}

func (b *SchemaBulder) AdditionalProperties(v *BoolOrSchema) *SchemaBulder {
	b.spec.Spec.AdditionalProperties = v
	return b
}

func (b *SchemaBulder) UnevaluatedProperties(v *BoolOrSchema) *SchemaBulder {
	b.spec.Spec.UnevaluatedProperties = v
	return b
}

func (b *SchemaBulder) PropertyNames(v *RefOrSpec[Schema]) *SchemaBulder {
	b.spec.Spec.PropertyNames = v
	return b
}

func (b *SchemaBulder) MinProperties(v int) *SchemaBulder {
	b.spec.Spec.MinProperties = &v
	return b
}

func (b *SchemaBulder) MaxProperties(v int) *SchemaBulder {
	b.spec.Spec.MaxProperties = &v
	return b
}

func (b *SchemaBulder) Required(v ...string) *SchemaBulder {
	b.spec.Spec.Required = v
	return b
}

func (b *SchemaBulder) AddRequired(v ...string) *SchemaBulder {
	b.spec.Spec.Required = append(b.spec.Spec.Required, v...)
	return b
}

func (b *SchemaBulder) Discriminator(v *Discriminator) *SchemaBulder {
	b.spec.Spec.Discriminator = v
	return b
}

func (b *SchemaBulder) XML(v *Extendable[XML]) *SchemaBulder {
	b.spec.Spec.XML = v
	return b
}

func (b *SchemaBulder) ExternalDocs(v *Extendable[ExternalDocs]) *SchemaBulder {
	b.spec.Spec.ExternalDocs = v
	return b
}

func (b *SchemaBulder) Example(v any) *SchemaBulder {
	b.spec.Spec.Example = v
	return b
}

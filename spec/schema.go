package spec

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
// https://spec.openapis.org/oas/v3.1.0#schema-object
type Schema struct {
	JsonSchema `yaml:",inline"`

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

// NewSchemaSpec creates Schema object.
func NewSchemaSpec() *RefOrSpec[Schema] {
	return NewRefOrSpec[Schema](nil, &Schema{})
}

// NewSchemaRef creates Ref object.
func NewSchemaRef(ref *Ref) *RefOrSpec[Schema] {
	return NewRefOrSpec[Schema](ref, nil)
}

// WithExt sets the extension and returns the current object (self|this).
// Schema does not require special `x-` prefix.
// The extension will be ignored if the name overlaps with a struct field during marshalling to JSON or YAML.
func (o *Schema) WithExt(name string, value any) *Schema {
	if o.Extensions == nil {
		o.Extensions = make(map[string]any, 1)
	}
	o.Extensions[name] = value
	return o
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

func (o *Schema) validateSpec(path string, opts *validationOptions) []*validationError {
	var errs []*validationError

	if o.Discriminator != nil {
		errs = append(errs, o.Discriminator.validateSpec(joinDot(path, "discriminator"), opts)...)
	}
	if o.XML != nil {
		errs = append(errs, o.XML.validateSpec(joinDot(path, "xml"), opts)...)
	}
	if o.ExternalDocs != nil {
		errs = append(errs, o.ExternalDocs.validateSpec(joinDot(path, "externalDocs"), opts)...)
	}
	if o.Example != nil {
		if !opts.doNotValidateExamples {
			if err := ValidateData(o.Example, o, opts.spec); err != nil {
				errs = append(errs, newValidationError(joinDot(path, "example"), err))
			}
		}
	}

	// JsonSchemaComposition
	if o.Not != nil {
		errs = append(errs, o.Not.validateSpec(joinDot(path, "not"), opts)...)
	}
	if o.AllOf != nil {
		for i, v := range o.AllOf {
			errs = append(errs, v.validateSpec(joinArrayItem(joinDot(path, "allOf"), i), opts)...)
		}
	}
	if o.AnyOf != nil {
		for i, v := range o.AnyOf {
			errs = append(errs, v.validateSpec(joinArrayItem(joinDot(path, "anyOf"), i), opts)...)
		}
	}
	if o.OneOf != nil {
		for i, v := range o.OneOf {
			errs = append(errs, v.validateSpec(joinArrayItem(joinDot(path, "oneOf"), i), opts)...)
		}
	}

	// JsonSchemaCore
	if o.Schema != "" && o.Schema != Draft202012 {
		errs = append(errs, newValidationError(joinDot(path, "schema"), "must be '%s', but got '%s'", Draft202012, o.Schema))
	}
	if len(o.Defs) > 0 {
		for k, v := range o.Defs {
			errs = append(errs, v.validateSpec(joinArrayItem(joinDot(path, "defs"), k), opts)...)
		}
	}
	if o.Type != nil {
		switch len(*o.Type) {
		case 0: // not type or any type
		case 1:
			switch v := (*o.Type)[0]; v {
			case StringType, NumberType, IntegerType, BooleanType, ObjectType, ArrayType, NullType:
			default:
				errs = append(errs, newValidationError(joinDot(path, "type"), "must be one of [%s, %s, %s, %s, %s, %s, %s], but got '%s'", StringType, NumberType, IntegerType, BooleanType, ObjectType, ArrayType, NullType, v))
			}
		default:
			for i, v := range *o.Type {
				switch v {
				case StringType, NumberType, IntegerType, BooleanType, ObjectType, ArrayType, NullType:
				default:
					errs = append(errs, newValidationError(joinArrayItem(joinDot(path, "type"), i), "must be one of [%s, %s, %s, %s, %s, %s, %s], but got '%s'", StringType, NumberType, IntegerType, BooleanType, ObjectType, ArrayType, NullType, v))
				}
			}
		}
	}

	// JsonSchemaMedia
	if o.ContentSchema != nil {
		errs = append(errs, o.ContentSchema.validateSpec(joinDot(path, "contentSchema"), opts)...)
	}
	if o.ContentEncoding != "" {
		switch o.ContentEncoding {
		case SevenBitEncoding, EightBitEncoding, BinaryEncoding, QuotedPrintableEncoding, Base16Encoding, Base32Encoding, Base64Encoding:
		default:
			errs = append(errs, newValidationError(joinDot(path, "contentEncoding"), "must be one of [%s, %s, %s, %s, %s, %s, %s], but got '%s'", SevenBitEncoding, EightBitEncoding, BinaryEncoding, QuotedPrintableEncoding, Base16Encoding, Base32Encoding, Base64Encoding, o.ContentEncoding))
		}
	}

	// JsonSchemaGeneric
	if o.Default != nil {
		if !opts.doNotValidateDefaultValues {
			if err := ValidateData(o.Default, o, opts.spec); err != nil {
				errs = append(errs, newValidationError(joinDot(path, "default"), err))
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
				errs = append(errs, newValidationError(joinDot(path, "default"), "must be one of enum values: %v", o.Enum))
			}
		}
	}

	if len(o.Examples) > 0 && !opts.doNotValidateExamples {
		for k, v := range o.Examples {
			if err := ValidateData(v, o, opts.spec); err != nil {
				errs = append(errs, newValidationError(joinArrayItem(joinDot(path, "examples"), k), err))
			}
		}
	}

	if o.Type != nil {
		for _, t := range *o.Type {
			switch t {
			case ArrayType: // JsonSchemaTypeArray
				if o.Items != nil {
					errs = append(errs, o.Items.validateSpec(joinDot(path, "items"), opts)...)
				}
				if o.MinItems != nil && *o.MinItems < 0 {
					errs = append(errs, newValidationError(joinDot(path, "minItems"), "must be greater than or equal to 0"))
				}
				if o.MaxItems != nil && *o.MaxItems < 0 {
					errs = append(errs, newValidationError(joinDot(path, "maxItems"), "must be greater than or equal to 0"))
					if o.MinItems != nil && *o.MaxItems < *o.MinItems {
						errs = append(errs, newValidationError(joinDot(path, "maxItems"), "must be greater than or equal to minItems"))
					}
				}
				if o.UnevaluatedItems != nil {
					errs = append(errs, o.UnevaluatedItems.validateSpec(joinDot(path, "unevaluatedItems"), opts)...)
				}
				if o.Contains != nil {
					errs = append(errs, o.Contains.validateSpec(joinDot(path, "contains"), opts)...)
				}
				if o.MinContains != nil && *o.MinContains < 0 {
					errs = append(errs, newValidationError(joinDot(path, "minContains"), "must be greater than or equal to 0"))
				}
				if o.MaxContains != nil && *o.MaxContains < 0 {
					errs = append(errs, newValidationError(joinDot(path, "maxContains"), "must be greater than or equal to 0"))
					if o.MinContains != nil && *o.MaxContains < *o.MinContains {
						errs = append(errs, newValidationError(joinDot(path, "maxContains"), "must be greater than or equal to minContains"))
					}
				}
				if len(o.PrefixItems) > 0 {
					for i, v := range o.PrefixItems {
						errs = append(errs, v.validateSpec(joinArrayItem(joinDot(path, "prefixItems"), i), opts)...)
					}
				}
			case ObjectType: // JsonSchemaTypeObject
				if o.Properties != nil {
					for k, v := range o.Properties {
						errs = append(errs, v.validateSpec(joinArrayItem(joinDot(path, "properties"), k), opts)...)
					}
				}
				if o.PatternProperties != nil {
					for k, v := range o.PatternProperties {
						errs = append(errs, v.validateSpec(joinArrayItem(joinDot(path, "patternProperties"), k), opts)...)
						if _, err := regexp.Compile(k); err != nil {
							errs = append(errs, newValidationError(joinArrayItem(joinDot(path, "patternProperties"), k), err))
						}
					}
				}
				if o.AdditionalProperties != nil {
					errs = append(errs, o.AdditionalProperties.validateSpec(joinDot(path, "additionalProperties"), opts)...)
				}
				if o.UnevaluatedItems != nil {
					errs = append(errs, o.UnevaluatedItems.validateSpec(joinDot(path, "unevaluatedItems"), opts)...)
				}
				if o.PropertyNames != nil {
					errs = append(errs, o.PropertyNames.validateSpec(joinDot(path, "propertyNames"), opts)...)
				}
				if o.MinProperties != nil && *o.MinProperties < 0 {
					errs = append(errs, newValidationError(joinDot(path, "minProperties"), "must be greater than or equal to 0"))
				}
				if o.MaxProperties != nil && *o.MaxProperties < 0 {
					errs = append(errs, newValidationError(joinDot(path, "maxProperties"), "must be greater than or equal to 0"))
					if o.MinProperties != nil && *o.MaxProperties < *o.MinProperties {
						errs = append(errs, newValidationError(joinDot(path, "maxProperties"), "must be greater than or equal to minProperties"))
					}
				}
				if len(o.Required) > 0 {
					for i, v := range o.Required {
						if _, ok := o.Properties[v]; !ok {
							errs = append(errs, newValidationError(joinArrayItem(joinDot(path, "required"), i), "must be a property in properties"))
						}
					}
				}
			case NumberType, IntegerType: // JsonSchemaTypeNumber
				if o.MultipleOf != nil && *o.MultipleOf <= 0 {
					errs = append(errs, newValidationError(joinDot(path, "multipleOf"), "must be greater than 0"))
				}
				if o.Minimum != nil && *o.Minimum < 0 {
					errs = append(errs, newValidationError(joinDot(path, "minimum"), "must be greater than or equal to 0"))
				}
				if o.Maximum != nil && *o.Maximum < 0 {
					errs = append(errs, newValidationError(joinDot(path, "maximum"), "must be greater than or equal to 0"))
					if o.Minimum != nil && *o.Maximum < *o.Minimum {
						errs = append(errs, newValidationError(joinDot(path, "maximum"), "must be greater than or equal to minimum"))
					}
				}
				if o.ExclusiveMinimum != nil && *o.ExclusiveMinimum < 0 {
					errs = append(errs, newValidationError(joinDot(path, "exclusiveMinimum"), "must be greater than or equal to 0"))
				}
				if o.ExclusiveMaximum != nil && *o.ExclusiveMaximum < 0 {
					errs = append(errs, newValidationError(joinDot(path, "exclusiveMaximum"), "must be greater than or equal to 0"))
					if o.ExclusiveMinimum != nil && *o.ExclusiveMaximum < *o.ExclusiveMinimum {
						errs = append(errs, newValidationError(joinDot(path, "exclusiveMaximum"), "must be greater than or equal to exclusiveMinimum"))
					}
				}
				if o.Minimum != nil && o.ExclusiveMinimum != nil {
					errs = append(errs, newValidationError(joinDot(path, "minimum&exclusiveMinimum"), ErrMutuallyExclusive))
				}
				if o.Maximum != nil && o.ExclusiveMaximum != nil {
					errs = append(errs, newValidationError(joinDot(path, "maximum&exclusiveMaximum"), ErrMutuallyExclusive))
				}
			case StringType: // JsonSchemaTypeString
				if o.MinLength != nil && *o.MinLength < 0 {
					errs = append(errs, newValidationError(joinDot(path, "minLength"), "must be greater than or equal to 0"))
				}
				if o.MaxLength != nil && *o.MaxLength < 0 {
					errs = append(errs, newValidationError(joinDot(path, "maxLength"), "must be greater than or equal to 0"))
					if o.MinLength != nil && *o.MaxLength < *o.MinLength {
						errs = append(errs, newValidationError(joinDot(path, "maxLength"), "must be greater than or equal to minLength"))
					}
				}
				if o.Pattern != "" {
					if _, err := regexp.Compile(o.Pattern); err != nil {
						errs = append(errs, newValidationError(joinDot(path, "pattern"), err))
					}
				}
			}
		}
	}
	return errs
}

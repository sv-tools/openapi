package openapi

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

type parseOptions struct {
	DefaultEncoding string // default encoding for the schema, used for []byte
}

type ParseOption func(*parseOptions)

// WithDefaultEncoding sets the default encoding for the schema for []byte fields.
// The default encoding is Base64Encoding.
func WithDefaultEncoding(v string) ParseOption {
	return func(opts *parseOptions) {
		opts.DefaultEncoding = v
	}
}

const is64Bit = uint64(^uintptr(0)) == ^uint64(0)

// ParseObject parses the object and returns the schema or the reference to the schema.
//
// The object can be a struct, pointer to struct, map, slice, pointer to map or slice, or any other type.
// The object can contain fields with `json`, `yaml` or `openapi` tags.
//
//	`openapi:"<name>[,ref:<ref> || any other tags]"` tag:
//	  - <name> is the name of the field in the schema, can be "-" to skip the field or empty to use the name from json, yaml tags or original field name.
//	json schema fields:
//	  - ref:<ref> is a reference to the schema, can not be used with jsonschema fields.
//	  - required, marks the field as required by adding it to the required list of the parent schema.
//	  - deprecated, marks the field as deprecated.
//	  - title:<title>, sets the title of the field or summary for the fereference.
//	  - summary:<summary>, sets the summary of the reference.
//	  - description:<description>, sets the description of the field.
//	  - type:<type> (boolean, integer, number, string, array, object), may be used multiple times.
//	    The first usage overrides the default type, all other types are added.
//	  - addtype:<type>, adds additional type, may be used multiple times.
//	  - format:<format>, sets the format of the type.
//
// The `components` parameter is needed to store the schemas of the structs, and to avoid the circular references.
// In case of the given object is struct, the function will return a reference to the schema stored in the components
// Otherwise, the function will return the schema itself.
func ParseObject(obj any, components *Extendable[Components], opts ...ParseOption) (*SchemaBuilder, error) {
	opt := &parseOptions{
		DefaultEncoding: Base64Encoding, // default encoding for []byte
	}
	for _, o := range opts {
		o(opt)
	}
	t := reflect.TypeOf(obj)
	if t == nil {
		return NewSchemaBuilder().Type(NullType).GoType("nil"), nil
	}
	return parseObject(joinLoc("", t.String()), t, components, opt)
}

// MapKeyMustBeStringError is an error that is returned when the map key is not a string.
type MapKeyMustBeStringError struct {
	Location string
	KeyType  reflect.Kind
}

func (e *MapKeyMustBeStringError) Error() string {
	return fmt.Sprintf("%s: unsupported map key type %s, expected string", e.Location, e.KeyType)
}

func (e *MapKeyMustBeStringError) Is(target error) bool {
	if target == nil {
		return false
	}
	_, ok := target.(*MapKeyMustBeStringError)
	return ok
}

// NewMapKeyMustBeStringError creates a new MapKeyMustBeStringError with the given location and key type.
func NewMapKeyMustBeStringError(location string, keyType reflect.Kind) *MapKeyMustBeStringError {
	return &MapKeyMustBeStringError{
		Location: location,
		KeyType:  keyType,
	}
}

func parseObject(location string, t reflect.Type, components *Extendable[Components], opt *parseOptions) (*SchemaBuilder, error) {
	if t == nil {
		return NewSchemaBuilder().Type(NullType).GoType("nil"), nil
	}
	kind := t.Kind()
	if kind == reflect.Ptr {
		builder, err := parseObject(location, t.Elem(), components, opt)
		if err != nil {
			return nil, err
		}
		if builder.IsRef() {
			builder = NewSchemaBuilder().OneOf(
				builder.Build(),
				NewSchemaBuilder().Type(NullType).Build(),
			)
		} else {
			builder.AddType(NullType)
		}
		return builder, nil
	}
	if kind == reflect.Interface {
		return NewSchemaBuilder().GoType("any"), nil
	}
	obj := reflect.New(t).Elem()
	builder := NewSchemaBuilder().GoType(fmt.Sprintf("%T", obj.Interface()))
	switch obj.Interface().(type) {
	case bool:
		builder.Type(BooleanType)
	case int, uint:
		if is64Bit {
			builder.Type(IntegerType).Format(Int64Format)
		} else {
			builder.Type(IntegerType).Format(Int32Format)
		}
	case int8, int16, int32, uint8, uint16, uint32:
		builder.Type(IntegerType).Format(Int32Format)
	case int64, uint64:
		builder.Type(IntegerType).Format(Int64Format)
	case float32:
		builder.Type(NumberType).Format(FloatFormat)
	case float64:
		builder.Type(NumberType).Format(DoubleFormat)
	case string:
		builder.Type(StringType)
	case []byte:
		builder.Type(StringType).ContentEncoding(opt.DefaultEncoding).GoType("[]byte")
	case json.Number:
		builder.Type(NumberType).GoPackage(t.PkgPath())
	case json.RawMessage:
		builder.Type(StringType).ContentMediaType("application/json").GoPackage(t.PkgPath())
	default:
		switch kind {
		case reflect.Array, reflect.Slice:
			var elemSchema any
			elem := t.Elem()
			if elem.Kind() == reflect.Interface {
				elemSchema = true
			} else {
				var err error
				elemSchema, err = parseObject(location, elem, components, opt)
				if err != nil {
					return nil, err
				}
			}
			builder.Type(ArrayType).Items(NewBoolOrSchema(elemSchema)).GoType("")
		case reflect.Map:
			if k := t.Key().Kind(); k != reflect.String {
				return nil, NewMapKeyMustBeStringError(location, k)
			}
			var elemSchema any
			elem := t.Elem()
			if elem.Kind() == reflect.Interface {
				elemSchema = true
			} else {
				var err error
				elemSchema, err = parseObject(location, elem, components, opt)
				if err != nil {
					return nil, err
				}
			}
			builder.Type(ObjectType).AdditionalProperties(NewBoolOrSchema(elemSchema)).GoType("")
		case reflect.Struct:
			objName := strings.ReplaceAll(t.PkgPath()+"."+t.Name(), "/", ".")
			if components.Spec.Schemas[objName] != nil {
				return NewSchemaBuilder().Ref("#/components/schemas/" + objName), nil
			}
			// add a temporary schema to avoid circular references
			if components.Spec.Schemas == nil {
				components.Spec.Schemas = make(map[string]*RefOrSpec[Schema], 1)
			}
			// reserve the name of the schema
			const toBeDeleted = "to be deleted"
			components.Spec.Schemas[objName] = NewSchemaBuilder().Ref(toBeDeleted).Build()
			var allOf []*RefOrSpec[Schema]
			for i := range t.NumField() {
				field := t.Field(i)
				// skip unexported fields
				if !field.IsExported() {
					continue
				}
				fieldSchema, err := parseObject(joinLoc(location, field.Name), obj.Field(i).Type(), components, opt)
				if err != nil {
					// remove the temporary schema
					delete(components.Spec.Schemas, objName)
					return nil, err
				}
				if field.Anonymous {
					allOf = append(allOf, fieldSchema.Build())
					continue
				}
				name := applyTag(&field, fieldSchema, builder)
				// skip the field if it's marked as "-"
				if name == "-" {
					continue
				}
				builder.AddProperty(name, fieldSchema.Build())
			}
			if len(allOf) > 0 {
				allOf = append(allOf, builder.Type(ObjectType).GoType("").Build())
				builder = NewSchemaBuilder().AllOf(allOf...).GoType(t.String())
			} else {
				builder.Type(ObjectType)
			}
			builder.GoPackage(t.PkgPath())
			components.Spec.Schemas[objName] = builder.Build()
			builder = NewSchemaBuilder().Ref("#/components/schemas/" + objName)
		default:
			// ignore unsupported types gracefully
		}
	}

	return builder, nil
}

func applyTag(field *reflect.StructField, schema, parent *SchemaBuilder) (name string) {
	name = field.Name

	for _, tagName := range []string{"json", "yaml"} {
		if tag, ok := field.Tag.Lookup(tagName); ok {
			parts := strings.SplitN(tag, ",", 2)
			if len(parts) > 0 {
				part := strings.TrimSpace(parts[0])
				if part != "" {
					name = part
					break
				}
			}
		}
	}

	tag, ok := field.Tag.Lookup("openapi")
	if !ok {
		return
	}
	parts := strings.Split(tag, ",")
	if len(parts) == 0 {
		return
	}

	if parts[0] != "" {
		name = parts[0]
	}
	if name == "-" {
		return parts[0]
	}
	parts = parts[1:]
	if len(parts) == 0 {
		return
	}

	if strings.HasPrefix(parts[0], "ref:") {
		schema.Ref(parts[0][4:])
	}

	var isTypeOverridden bool

	for _, part := range parts {
		prefixIndex := strings.Index(part, ":")
		var prefix string
		if prefixIndex == -1 {
			prefix = part
		} else {
			prefix = part[:prefixIndex]
			if prefixIndex == len(part)-1 {
				part = ""
			}
			part = part[prefixIndex+1:]
		}

		// the tags for the references only
		if schema.IsRef() {
			switch prefix {
			case "required":
				parent.AddRequired(name)
			case "description":
				schema.Description(part)
			case "title", "summary":
				schema.Title(part)
			default:
				// ignore unknown or unsupported tag prefixes gracefully
			}
			continue
		}

		switch prefix {
		case "required":
			parent.AddRequired(name)
		case "deprecated":
			schema.Deprecated(true)
		case "title":
			schema.Title(part)
		case "description":
			schema.Description(part)
		case "type":
			// first type overrides the default type, all other types are added
			if !isTypeOverridden {
				schema.Type(part)
				isTypeOverridden = true
			} else {
				schema.AddType(part)
			}
		case "addtype":
			schema.AddType(part)
		case "format":
			schema.Format(part)
		default:
			// handle unknown or unsupported tag prefixes gracefully
		}
	}

	return
}

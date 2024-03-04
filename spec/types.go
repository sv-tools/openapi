package spec

import (
	"fmt"
	"reflect"
)

const (
	// ******* Type-specific keywords *******
	//
	// https://json-schema.org/understanding-json-schema/reference/type.html

	// StringType is used for strings of text. It may contain Unicode characters.
	//
	// https://json-schema.org/understanding-json-schema/reference/string.html#string
	StringType = "string"
	// NumberType is used for any numeric type, either integers or floating point numbers.
	//
	// https://json-schema.org/understanding-json-schema/reference/numeric.html#number
	NumberType = "number"
	// IntegerType is used for integral numbers.
	// JSON does not have distinct types for integers and floating-point values.
	// Therefore, the presence or absence of a decimal point is not enough to distinguish between integers and non-integers.
	// For example, 1 and 1.0 are two ways to represent the same value in JSON.
	// JSON Schema considers that value an integer no matter which representation was used.
	//
	// https://json-schema.org/understanding-json-schema/reference/numeric.html#integer
	IntegerType = "integer"
	// ObjectType is the mapping type in JSON.
	// They map “keys” to “values”.
	// In JSON, the “keys” must always be strings.
	// Each of these pairs is conventionally referred to as a “property”.
	//
	// https://json-schema.org/understanding-json-schema/reference/object.html#object
	ObjectType = "object"
	// ArrayType is used for ordered elements.
	// In JSON, each element in an array may be of a different type.
	//
	// https://json-schema.org/understanding-json-schema/reference/array.html#array
	ArrayType = "array"
	// BooleanType matches only two special values: true and false.
	// Note that values that evaluate to true or false, such as 1 and 0, are not accepted by the schema.
	//
	// https://json-schema.org/understanding-json-schema/reference/boolean.html#boolean
	BooleanType = "boolean"
	// NullType has only one acceptable value: null.
	//
	// https://json-schema.org/understanding-json-schema/reference/null.html#null
	NullType = "null"

	// ******* Media: string-encoding non-JSON data *******
	//
	// https://json-schema.org/understanding-json-schema/reference/non_json_data.html

	SevenBitEncoding        = "7bit"
	EightBitEncoding        = "8bit"
	BinaryEncoding          = "binary"
	QuotedPrintableEncoding = "quoted-printable"
	Base16Encoding          = "base16"
	Base32Encoding          = "base32"
	Base64Encoding          = "base64"
)

// GetType returns the JSON Schema type of the given value.
func GetType(v any) (string, error) {
	if v == nil {
		return NullType, nil
	}
	return kindToType(getKind(v))
}

func getKind(v any) reflect.Kind {
	k := reflect.TypeOf(v).Kind()
	if k == reflect.Ptr {
		k = reflect.TypeOf(v).Elem().Kind()
	}
	return k
}

func kindToType(k reflect.Kind) (string, error) {
	switch k {
	case reflect.String:
		return StringType, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return IntegerType, nil
	case reflect.Float32, reflect.Float64:
		return NumberType, nil
	case reflect.Bool:
		return BooleanType, nil
	case reflect.Map, reflect.Struct:
		return ObjectType, nil
	case reflect.Slice, reflect.Array:
		return ArrayType, nil
	default:
		return "", fmt.Errorf("unsupported type: %s", k.String())
	}
}

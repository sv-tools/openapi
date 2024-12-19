package openapi

import "reflect"

const is64Bit = uint64(^uintptr(0)) == ^uint64(0)

func ParseObject(obj any) (*SchemaBulder, error) {
	t := reflect.TypeOf(obj)
	if t == nil && obj == nil {
		return NewSchemaBuilder().Type(NullType)
	}
	kind := t.Kind()
	if kind == reflect.Ptr {
		kind = t.Elem().Kind()
	}
	builder := NewSchemaBuilder().GoType(kind.String())
	switch obj.(type) {
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
	case complex64, complex128:
		builder.Type(ObjectType).OneOf(
			NewSchemaBuilder().Type(ObjectType).
				AddProperty("real", NewSchemaBuilder().Type(NumberType).Format(DoubleFormat).Build()).
				AddProperty("imaginary", NewSchemaBuilder().Type(NumberType).Format(DoubleFormat).Build()).
				Build(),
			NewSchemaBuilder().Type(ObjectType).
				AddProperty("real", NewSchemaBuilder().Type(NumberType).Format(DoubleFormat).Build()).
				AddProperty("imag", NewSchemaBuilder().Type(NumberType).Format(DoubleFormat).Build()).
				Build(),
		)
	case string, []rune:
		builder.Type(StringType)
	case rune:
		builder.Type(StringType).MinLength(1).MaxLength(1)
	case []byte:
		builder.Type(StringType).ContentEncoding(Base64Encoding).ContentSchema(NewSchemaBuilder().Type(StringType).Build())
	default:
		switch kind {
		case reflect.Array, reflect.Slice:

			builder.Type(ArrayType).Items()
		Map
		Struct
	}

	return builder
}

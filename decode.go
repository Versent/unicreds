package unicreds

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// decode decode the supplied struct from the dynamodb result map
func decode(name string, data map[string]*dynamodb.AttributeValue, val interface{}) (err error) {
	if data == nil {
		// If the data is nil, then we don't set anything.
		return nil
	}

	fields := make(map[*reflect.StructField]reflect.Value)

	v := reflect.ValueOf(val)
	if v.Kind() != reflect.Ptr {
		return errors.New("result must be a pointer")
	}

	v = v.Elem()
	if !v.CanAddr() {
		return errors.New("result must be addressable (a pointer)")
	}

	d := reflect.ValueOf(data)
	if !d.IsValid() {
		// If the data value is invalid, then we just set the value
		// to be the zero value.
		v.Set(reflect.Zero(v.Type()))
		return nil
	}

	if getKind(v) != reflect.Struct {
		return fmt.Errorf("%s: unsupported type: %s", name, getKind(v))
	}

	structType := v.Type()

	for i := 0; i < structType.NumField(); i++ {
		fieldType := structType.Field(i)

		// Normal struct field, store it away
		fields[&fieldType] = v.Field(i)
	}

	for fieldType, field := range fields {
		fieldName := fieldType.Name

		tagValue := fieldType.Tag.Get("ds")
		if tagValue != "" {
			fieldName = tagValue
		}

		if k := data[fieldName]; k == nil {
			continue
		}

		switch getKind(field) {
		case reflect.String:
			err = decodeString(fieldName, data[fieldName], field)
		case reflect.Int:
			err = decodeInt(fieldName, data[fieldName], field)
		default:
			return fmt.Errorf("%s: unsupported type: %s", fieldName, getKind(field))
		}
	}

	return err
}

func decodeInt(name string, data *dynamodb.AttributeValue, val reflect.Value) error {

	if data.N == nil {
		return nil
	}

	i, err := strconv.ParseInt(*data.N, 0, val.Type().Bits())

	if err == nil {
		val.SetInt(i)
	} else {
		return fmt.Errorf("cannot parse '%s' as int: %s", name, err)
	}

	return nil
}

func decodeString(name string, data *dynamodb.AttributeValue, val reflect.Value) error {

	if data.S == nil {
		return nil
	}
	val.SetString(*data.S)

	return nil
}

func getKind(val reflect.Value) reflect.Kind {
	kind := val.Kind()

	switch {
	case kind >= reflect.Int && kind <= reflect.Int64:
		return reflect.Int
	case kind >= reflect.Uint && kind <= reflect.Uint64:
		return reflect.Uint
	case kind >= reflect.Float32 && kind <= reflect.Float64:
		return reflect.Float32
	default:
		return kind
	}
}

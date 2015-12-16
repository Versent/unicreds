package unicreds

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// Decode decode the supplied struct from the dynamodb result map
//
// NOTE: this function needs a lot more validation and refinement.
func Decode(data map[string]*dynamodb.AttributeValue, rawVal interface{}) error {
	val := reflect.ValueOf(rawVal)
	if val.Kind() != reflect.Ptr {
		return errors.New("result must be a pointer")
	}

	val = val.Elem()
	if !val.CanAddr() {
		return errors.New("result must be addressable (a pointer)")
	}
	return decode("ds", data, val)
}

func decode(name string, data map[string]*dynamodb.AttributeValue, val reflect.Value) error {
	if data == nil {
		// If the data is nil, then we don't set anything.
		return nil
	}

	dataVal := reflect.ValueOf(data)
	if !dataVal.IsValid() {
		// If the data value is invalid, then we just set the value
		// to be the zero value.
		val.Set(reflect.Zero(val.Type()))
		return nil
	}

	var err error
	dataKind := getKind(val)
	switch dataKind {
	case reflect.Struct:
		err = decodeStruct(name, data, val)
	default:
		return fmt.Errorf("%s: unsupported type: %s", name, dataKind)
	}

	return err
}

func decodeStruct(name string, data map[string]*dynamodb.AttributeValue, val reflect.Value) (err error) {

	fields := make(map[*reflect.StructField]reflect.Value)

	structVal := val
	structType := structVal.Type()

	for i := 0; i < structType.NumField(); i++ {
		fieldType := structType.Field(i)

		// Normal struct field, store it away
		fields[&fieldType] = structVal.Field(i)
	}

	for fieldType, field := range fields {
		fieldName := fieldType.Name

		tagValue := fieldType.Tag.Get("ds")
		if tagValue != "" {
			fieldName = tagValue
		}

		keyVal := data[fieldName]
		if keyVal == nil {
			continue
		}

		switch getKind(field) {
		case reflect.String:
			err = decodeString(fieldName, keyVal, field)
		case reflect.Int:
			err = decodeInt(fieldName, keyVal, field)
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

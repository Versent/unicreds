package unicreds

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// Encode return the value encoded as a map of dynamo attributes.
//
// NOTE: this function needs a lot more validation and refinement.
func Encode(rawVal interface{}) (map[string]*dynamodb.AttributeValue, error) {
	val := reflect.ValueOf(rawVal)
	t := reflect.TypeOf(rawVal)

	var err error
	var data map[string]*dynamodb.AttributeValue

	switch t.Kind() {
	case reflect.Struct:
		data, err = encodeStruct("ds", val)
	case reflect.Ptr:
		val = val.Elem()
		data, err = encodeStruct("ds", val)

	default:
		return data, fmt.Errorf("%s: unsupported type: %s", "ds", t)
	}

	return data, err
}

func encodeStruct(name string, val reflect.Value) (map[string]*dynamodb.AttributeValue, error) {

	structType := val.Type()

	fields := make(map[*reflect.StructField]reflect.Value)
	data := map[string]*dynamodb.AttributeValue{}

	for i := 0; i < structType.NumField(); i++ {
		fieldType := structType.Field(i)

		// Normal struct field, store it away
		fields[&fieldType] = val.Field(i)
	}

	for fieldType, field := range fields {
		fieldName := fieldType.Name
		tagValue := fieldType.Tag.Get("ds")
		if tagValue != "" {
			fieldName = tagValue
		}

		// if this field is "empty" then don't save it's value
		if field.Interface() == reflect.Zero(fieldType.Type).Interface() {
			continue
		}

		switch getKind(field) {
		case reflect.String:
			data[fieldName] = &dynamodb.AttributeValue{
				S: aws.String(field.String()),
			}
		case reflect.Int:
			data[fieldName] = &dynamodb.AttributeValue{
				N: aws.String(strconv.FormatInt(field.Int(), 10)),
			}
		default:
			return data, fmt.Errorf("%s: unsupported type: %s", fieldName, getKind(field))
		}
	}

	return data, nil
}

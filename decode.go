package unicreds

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Decode decode the supplied struct from the dynamodb result map
func Decode(data map[string]*dynamodb.AttributeValue, rawVal interface{}) error {
	// Fix for issue https://github.com/fugue/credstash/issues/154 in credstash
	// This is needed until this issue is resolved, and also until we push new
	// values into credstash that have a string hmac value and not a binary hmac
	// value
	if h, ok := data["hmac"]; ok && len(h.B) > 0 {
		hmac := string(h.B)
		data["hmac"] = &dynamodb.AttributeValue{S: &hmac}
	}
	return dynamodbattribute.UnmarshalMap(data, rawVal)
}

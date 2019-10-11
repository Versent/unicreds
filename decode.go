package unicreds

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// adjustHmac will force the hmac to be a byte array if present as string
func adjustHmac(record map[string]*dynamodb.AttributeValue) {
	if val, ok := record["hmac"]; ok {
		if len(val.B) == 0 && val.S != nil {
			val.B = []byte(*val.S)
			val.S = nil
		}
	}
}

// Decode decode the supplied struct from the dynamodb result map
func Decode(data map[string]*dynamodb.AttributeValue, rawVal interface{}) error {
	// Fix for issue https://github.com/fugue/credstash/issues/154 in credstash
	// This is needed until this issue is resolved, and also until we push new
	// values into credstash that have a string hmac value and not a binary hmac
	// value
	//	if h, ok := data["hmac"]; ok && len(h.B) > 0 {
	//		hmac := string(h.B)
	//		data["hmac"] = &dynamodb.AttributeValue{S: &hmac}
	//	}

	adjustHmac(data)
	return dynamodbattribute.UnmarshalMap(data, rawVal)
}

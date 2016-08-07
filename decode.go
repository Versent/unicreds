package unicreds

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Decode decode the supplied struct from the dynamodb result map
func Decode(data map[string]*dynamodb.AttributeValue, rawVal interface{}) error {
	return dynamodbattribute.UnmarshalMap(data, rawVal)
}

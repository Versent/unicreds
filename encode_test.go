package unicreds

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
)

func TestEncode(t *testing.T) {

	cred := struct {
		Name      string `dynamodbav:"name"`
		Timestamp int64  `dynamodbav:"timestamp"`
	}{
		Name:      "data",
		Timestamp: 1449038525717338459,
	}

	expectedData := map[string]*dynamodb.AttributeValue{
		"name": &dynamodb.AttributeValue{
			S: aws.String("data"),
		},
		"timestamp": &dynamodb.AttributeValue{
			N: aws.String("1449038525717338459"),
		},
	}

	data, err := Encode(&cred)

	assert.Nil(t, err)
	assert.Equal(t, expectedData, data)
}

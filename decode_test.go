package unicreds

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
)

func TestDecode(t *testing.T) {

	cred := struct {
		Name      string `ds:"name"`
		Timestamp int64  `ds:"timestamp"`
	}{}

	data := map[string]*dynamodb.AttributeValue{
		"name": &dynamodb.AttributeValue{
			S: aws.String("data"),
		},
		"timestamp": &dynamodb.AttributeValue{
			N: aws.String("1449038525717338459"),
		},
	}

	err := Decode(data, &cred)
	if err != nil {
		fmt.Printf("%+v\n", err)
	}

	assert.Equal(t, "data", cred.Name)
	assert.Equal(t, int64(1449038525717338459), cred.Timestamp)
}

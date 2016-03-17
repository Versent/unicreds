package unicreds

import (
	"testing"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/atlassian/unicreds/mocks"
)

func init() {
	log.SetHandler(cli.Default)
}

func TestCredential(t *testing.T) {
	c := &Credential{}

	assert.Equal(t, c.CreatedAtDate(), CreatedAtNotAvailable)

	c.CreatedAt = 1458117788

	assert.NotEqual(t, c.CreatedAtDate(), CreatedAtNotAvailable)
}

func TestSetup(t *testing.T) {

	dsMock := configureMock()

	dsMock.On("CreateTable",
		mock.AnythingOfType("*dynamodb.CreateTableInput")).Return(nil, nil)

	dto := &dynamodb.DescribeTableOutput{
		Table: &dynamodb.TableDescription{TableStatus: aws.String("ACTIVE")},
	}

	dsMock.On("DescribeTable",
		mock.AnythingOfType("*dynamodb.DescribeTableInput")).Return(dto, nil)

	err := Setup()

	assert.Nil(t, err)
}

func TestGetSecretNotFound(t *testing.T) {

	dsMock := configureMock()

	qi := &dynamodb.QueryOutput{
		Items: []map[string]*dynamodb.AttributeValue{},
	}

	dsMock.On("Query", mock.AnythingOfType("*dynamodb.QueryInput")).Return(qi, nil)

	ds, err := GetSecret("test")

	assert.Error(t, err, "Secret Not Found")
	assert.Nil(t, ds)
}

func configureMock() *mocks.DynamoDBAPI {
	dsMock := &mocks.DynamoDBAPI{}

	dynamoSvc = dsMock

	return dsMock
}

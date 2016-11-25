package unicreds

import (
	"testing"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/versent/unicreds/mocks"
)

var (
	tableName     = "credential-store"
	readCapacity  = int64(4)
	writeCapacity = int64(4)
	dsPlainText   = []byte{
		0x6a, 0xcf, 0xeb, 0xd6, 0xe9, 0xa6, 0x19, 0xc1,
		0x38, 0xb9, 0xfc, 0x2d, 0x53, 0x23, 0x4d, 0x78,
		0x85, 0x48, 0x96, 0xd6, 0xd2, 0xf6, 0xf4, 0x42,
		0x99, 0x9d, 0x8e, 0xa9, 0xed, 0xf0, 0xb3, 0xf2,
	}

	itemsFixture = []map[string]*dynamodb.AttributeValue{
		{
			"name":     &dynamodb.AttributeValue{S: aws.String("test")},
			"version":  &dynamodb.AttributeValue{S: aws.String("1")},
			"contents": &dynamodb.AttributeValue{S: aws.String("o8we1zr9GD+KstVv3x2YTeT2")},
			"hmac":     &dynamodb.AttributeValue{S: aws.String("1e2d485cf52ec57d9db5c05eda678b45eee8d3dabcc6c1ee7c0999712026f6aa")},
		},
	}
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

	dsMock, _ := configureMock()

	dsMock.On("CreateTable",
		mock.AnythingOfType("*dynamodb.CreateTableInput")).Return(nil, nil)

	dto := &dynamodb.DescribeTableOutput{
		Table: &dynamodb.TableDescription{TableStatus: aws.String("ACTIVE")},
	}

	dsMock.On("DescribeTable",
		mock.AnythingOfType("*dynamodb.DescribeTableInput")).Return(dto, nil)

	err := Setup(&tableName, &readCapacity, &writeCapacity)

	assert.Nil(t, err)
}

func TestGetHighestVersionSecretNotFound(t *testing.T) {

	dsMock, _ := configureMock()

	qi := &dynamodb.QueryOutput{
		Items: []map[string]*dynamodb.AttributeValue{},
	}

	dsMock.On("Query", mock.AnythingOfType("*dynamodb.QueryInput")).Return(qi, nil)

	ds, err := GetHighestVersionSecret(&tableName, "test", NewEncryptionContextValue())

	assert.Error(t, err, "Secret Not Found")
	assert.Nil(t, ds)
}

func TestGetHighestVersionSecret(t *testing.T) {

	dsMock, kmsMock := configureMock()

	qi := &dynamodb.QueryOutput{
		Items: itemsFixture,
	}

	ki := &kms.DecryptOutput{Plaintext: dsPlainText}

	dsMock.On("Query", mock.AnythingOfType("*dynamodb.QueryInput")).Return(qi, nil)
	kmsMock.On("Decrypt", mock.AnythingOfType("*kms.DecryptInput")).Return(ki, nil)

	ds, err := GetHighestVersionSecret(&tableName, "test", NewEncryptionContextValue())

	assert.Nil(t, err)
	assert.Equal(t, ds.Secret, "something test 123")
}

func TestGetSecretNotFound(t *testing.T) {

	dsMock, _ := configureMock()

	gi := &dynamodb.GetItemOutput{
		Item: map[string]*dynamodb.AttributeValue{},
	}

	dsMock.On("GetItem", mock.AnythingOfType("*dynamodb.GetItemInput")).Return(gi, nil)

	ds, err := GetSecret(&tableName, "test", "1", NewEncryptionContextValue())

	assert.Error(t, err, "Secret Not Found")
	assert.Nil(t, ds)
}

func TestGetSecret(t *testing.T) {

	dsMock, kmsMock := configureMock()

	gi := &dynamodb.GetItemOutput{
		Item: itemsFixture[0],
	}

	ki := &kms.DecryptOutput{Plaintext: dsPlainText}

	dsMock.On("GetItem", mock.AnythingOfType("*dynamodb.GetItemInput")).Return(gi, nil)
	kmsMock.On("Decrypt", mock.AnythingOfType("*kms.DecryptInput")).Return(ki, nil)

	ds, err := GetSecret(&tableName, "test", "1", NewEncryptionContextValue())

	assert.Nil(t, err)
	assert.Equal(t, ds.Secret, "something test 123")
}

func TestGetAllSecrets(t *testing.T) {

	dsMock, kmsMock := configureMock()

	qs := &dynamodb.ScanOutput{
		Count: aws.Int64(0),
		Items: itemsFixture,
	}

	ki := &kms.DecryptOutput{Plaintext: dsPlainText}

	dsMock.On("Scan", mock.AnythingOfType("*dynamodb.ScanInput")).Return(qs, nil)
	kmsMock.On("Decrypt", mock.AnythingOfType("*kms.DecryptInput")).Return(ki, nil)

	ds, err := GetAllSecrets(&tableName, false, NewEncryptionContextValue())

	assert.Nil(t, err)
	assert.Len(t, ds, 1)
}

func TestGetAllSecretsDecryptFailed(t *testing.T) {

	dsMock, kmsMock := configureMock()

	qs := &dynamodb.ScanOutput{
		Count: aws.Int64(0),
		Items: itemsFixture,
	}

	awsErr := awserr.New("AccessDeniedException", "KMS access denied", nil)

	dsMock.On("Scan", mock.AnythingOfType("*dynamodb.ScanInput")).Return(qs, nil)
	kmsMock.On("Decrypt", mock.AnythingOfType("*kms.DecryptInput")).Return(nil, awsErr)

	ds, err := GetAllSecrets(&tableName, true, NewEncryptionContextValue())

	assert.Nil(t, err)
	assert.Len(t, ds, 0)
}

func TestGetAllSecretsEncryptionContextFailed(t *testing.T) {

	dsMock, kmsMock := configureMock()

	qs := &dynamodb.ScanOutput{
		Count: aws.Int64(0),
		Items: itemsFixture,
	}

	awsErr := awserr.New("InvalidCiphertextException", "", nil)

	dsMock.On("Scan", mock.AnythingOfType("*dynamodb.ScanInput")).Return(qs, nil)
	kmsMock.On("Decrypt", mock.AnythingOfType("*kms.DecryptInput")).Return(nil, awsErr)

	ec := NewEncryptionContextValue()
	ec.Set("Unknown:Context")

	ds, err := GetAllSecrets(&tableName, true, ec)

	assert.Nil(t, err)
	assert.Len(t, ds, 0)
}

func TestListSecrets(t *testing.T) {

	dsMock, _ := configureMock()

	qs := &dynamodb.ScanOutput{
		Count: aws.Int64(0),
		Items: itemsFixture,
	}

	dsMock.On("Scan", mock.AnythingOfType("*dynamodb.ScanInput")).Return(qs, nil)

	ds, err := ListSecrets(&tableName, true)

	assert.Nil(t, err)
	assert.Len(t, ds, 1)
}

func configureMock() (*mocks.DynamoDBAPI, *mocks.KMSAPI) {
	dsMock := &mocks.DynamoDBAPI{}
	kmsMock := &mocks.KMSAPI{}

	dynamoSvc = dsMock
	kmsSvc = kmsMock

	return dsMock, kmsMock
}

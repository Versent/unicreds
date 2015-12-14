package unicreds

import (
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

const (
	// Table the name of the dynamodb table
	Table = "credential-store"

	// KmsKey default KMS key alias name
	KmsKey = "alias/credstash"

	tableCreateTimeout = 30 * time.Second
)

var (
	dynamoSvc dynamodbiface.DynamoDBAPI

	// ErrSecretNotFound returned when unable to find the specified secret in dynamodb
	ErrSecretNotFound = errors.New("Secret Not Found")

	// ErrHmacValidationFailed returned when the hmac signature validation fails
	ErrHmacValidationFailed = errors.New("Secret HMAC validation failed")

	// ErrTimeout timeout occured waiting for dynamodb table to create
	ErrTimeout = errors.New("Timed out waiting for dynamodb table to become active")
)

func init() {
	dynamoSvc = dynamodb.New(session.New(), aws.NewConfig())
}

// Credential managed credential information
type Credential struct {
	Name     string `ds:"name"`
	Version  string `ds:"version"`
	Key      string `ds:"key"`
	Contents string `ds:"contents"`
	Hmac     string `ds:"hmac"`
}

// DecryptedCredential managed credential information
type DecryptedCredential struct {
	*Credential
	Secret string
}

// Setup create the table which stores credentials
func Setup() (err error) {

	res, err := dynamoSvc.CreateTable(&dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("name"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("version"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("name"),
				KeyType:       aws.String(dynamodb.KeyTypeHash),
			},
			{
				AttributeName: aws.String("version"),
				KeyType:       aws.String(dynamodb.KeyTypeRange),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
		TableName: aws.String(Table),
	})

	if err != nil {
		return
	}

	fmt.Printf("res = %+v\n", res)

	err = waitForTable()

	return
}

// GetSecret retrieve the secret from dynamodb using the name
func GetSecret(name string) (*DecryptedCredential, error) {

	res, err := dynamoSvc.Query(&dynamodb.QueryInput{
		TableName: aws.String(Table),
		ExpressionAttributeNames: map[string]*string{
			"#N": aws.String("name"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":name": &dynamodb.AttributeValue{
				S: aws.String(name),
			},
		},
		KeyConditionExpression: aws.String("#N = :name"),
		Limit:            aws.Int64(1),
		ConsistentRead:   aws.Bool(true),
		ScanIndexForward: aws.Bool(false), // descending order
	})

	if err != nil {
		return nil, err
	}

	cred := new(Credential)

	if len(res.Items) == 0 {
		return nil, ErrSecretNotFound
	}

	err = Decode(res.Items[0], cred)

	if err != nil {
		return nil, err
	}

	return decryptCredential(cred)
}

// ListSecrets return a list of secrets
func ListSecrets() ([]*DecryptedCredential, error) {

	res, err := dynamoSvc.Scan(&dynamodb.ScanInput{
		TableName: aws.String(Table),
		AttributesToGet: []*string{
			aws.String("name"),
			aws.String("version"),
			aws.String("key"),
			aws.String("contents"),
			aws.String("hmac"),
		},
		ConsistentRead: aws.Bool(true),
	})

	if err != nil {
		return nil, err
	}

	var results []*DecryptedCredential

	for _, item := range res.Items {
		cred := new(Credential)

		err = Decode(item, cred)
		if err != nil {
			return nil, err
		}

		dcred, err := decryptCredential(cred)
		if err != nil {
			return nil, err
		}

		results = append(results, dcred)
	}

	return results, nil
}

// PutSecret retrieve the secret from dynamodb
func PutSecret(name, secret, version string) error {

	if version == "" {
		version = "1"
	}

	dk, err := GenerateDataKey(KmsKey, 64)

	dataKey := dk.Plaintext[:32]
	hmacKey := dk.Plaintext[32:]
	wrappedKey := dk.CiphertextBlob

	ctext, err := Encrypt(dataKey, []byte(secret))

	b64hmac := ComputeHmac256(ctext, hmacKey)

	b64ctext := base64.StdEncoding.EncodeToString(ctext)

	cred := &Credential{
		Name:     name,
		Version:  version,
		Key:      base64.StdEncoding.EncodeToString(wrappedKey),
		Contents: b64ctext,
		Hmac:     b64hmac,
	}

	data, err := Encode(cred)

	if err != nil {
		return err
	}

	_, err = dynamoSvc.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(Table),
		Item:      data,
	})

	return err
}

// DeleteSecret delete a secret
func DeleteSecret(name string) error {

	res, err := dynamoSvc.Query(&dynamodb.QueryInput{
		TableName: aws.String(Table),
		ExpressionAttributeNames: map[string]*string{
			"#N": aws.String("name"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":name": &dynamodb.AttributeValue{
				S: aws.String(name),
			},
		},
		KeyConditionExpression: aws.String("#N = :name"),
		ConsistentRead:         aws.Bool(true),
		ScanIndexForward:       aws.Bool(false), // descending order
	})

	if err != nil {
		return err
	}

	for _, item := range res.Items {
		cred := new(Credential)

		err = Decode(item, cred)
		if err != nil {
			return err
		}

		fmt.Printf("deleting name=%s version=%s\n", cred.Name, cred.Version)

		_, err = dynamoSvc.DeleteItem(&dynamodb.DeleteItemInput{
			TableName: aws.String(Table),
			Key: map[string]*dynamodb.AttributeValue{
				"name": &dynamodb.AttributeValue{
					S: aws.String(cred.Name),
				},
				"version": &dynamodb.AttributeValue{
					S: aws.String(cred.Version),
				},
			},
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func decryptCredential(cred *Credential) (*DecryptedCredential, error) {

	wrappedKey, err := base64.StdEncoding.DecodeString(cred.Key)

	if err != nil {
		return nil, err
	}

	dk, err := DecryptDataKey(wrappedKey)

	if err != nil {
		return nil, err
	}

	dataKey := dk.Plaintext[:32]
	hmacKey := dk.Plaintext[32:]

	contents, err := base64.StdEncoding.DecodeString(cred.Contents)
	if err != nil {
		return nil, err
	}

	hexhmac := ComputeHmac256(contents, hmacKey)

	if hexhmac != cred.Hmac {
		return nil, ErrHmacValidationFailed
	}

	secret, err := Decrypt(dataKey, contents)

	if err != nil {
		return nil, err
	}

	plainText := string(secret)

	return &DecryptedCredential{Credential: cred, Secret: plainText}, nil
}

func waitForTable() error {

	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(tableCreateTimeout)
		timeout <- true
	}()

	ticker := time.NewTicker(1 * time.Second)

	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// a read from ch has occurred
			res, err := dynamoSvc.DescribeTable(&dynamodb.DescribeTableInput{
				TableName: aws.String(Table),
			})

			if err != nil {
				return err
			}

			if *res.Table.TableStatus == "ACTIVE" {
				return nil
			}

		case <-timeout:
			// polling for table status has taken more than the timeout
			return ErrTimeout
		}
	}

}

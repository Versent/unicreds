package unicreds

import (
	"encoding/base64"
	"errors"
	"sort"
	"strconv"
	"time"

	"github.com/apex/log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

const (
	// Table the name of the dynamodb table
	Table = "credential-store"

	// DefaultKmsKey default KMS key alias name
	DefaultKmsKey = "alias/credstash"

	// CreatedAtNotAvailable returned to indicate the created at field is missing
	// from the secret
	CreatedAtNotAvailable = "Not Available"

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

// SetDynamoDBConfig override the default aws configuration
func SetDynamoDBConfig(config *aws.Config) {
	dynamoSvc = dynamodb.New(session.New(), config)
}

// Credential managed credential information
type Credential struct {
	Name      string `ds:"name"`
	Version   string `ds:"version"`
	Key       string `ds:"key"`
	Contents  string `ds:"contents"`
	Hmac      string `ds:"hmac"`
	CreatedAt int64  `ds:"created_at"`
}

// CreatedAtDate convert the timestamp field to a date string
func (c *Credential) CreatedAtDate() string {
	if c.CreatedAt == 0 {
		return CreatedAtNotAvailable
	}
	tm := time.Unix(c.CreatedAt, 0)
	return tm.String()
}

// DecryptedCredential managed credential information
type DecryptedCredential struct {
	*Credential
	Secret string
}

// ByVersion sort helper for credentials
type ByVersion []*Credential

func (a ByVersion) Len() int      { return len(a) }
func (a ByVersion) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func (a ByVersion) Less(i, j int) bool {
	aiv, _ := strconv.Atoi(a[i].Version)
	ajv, _ := strconv.Atoi(a[j].Version)

	return aiv < ajv
}

// Setup create the table which stores credentials
func Setup() (err error) {
	log.Debug("Running Setup")

	_, err = dynamoSvc.CreateTable(&dynamodb.CreateTableInput{
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

	err = waitForTable()

	return
}

// GetSecret retrieve the secret from dynamodb using the name
func GetSecret(name string) (*DecryptedCredential, error) {
	log.Debug("Getting secret")

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

// GetHighestVersion look up the highest version for a given name
func GetHighestVersion(name string) (string, error) {
	log.WithField("name", name).Debug("Looking up highest version")

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
		Limit:                aws.Int64(1),
		ConsistentRead:       aws.Bool(true),
		ScanIndexForward:     aws.Bool(false), // descending order
		ProjectionExpression: aws.String("version"),
	})

	if err != nil {
		return "", err
	}

	if len(res.Items) == 0 {
		return "", ErrSecretNotFound
	}

	v := res.Items[0]["version"]

	if v == nil {
		return "", ErrSecretNotFound
	}

	return aws.StringValue(v.S), nil
}

// ListSecrets returns a list of all secrets
func ListSecrets(all bool) ([]*Credential, error) {
	log.Debug("Listing secrets")

	res, err := dynamoSvc.Scan(&dynamodb.ScanInput{
		TableName: aws.String(Table),
		ExpressionAttributeNames: map[string]*string{
			"#N": aws.String("name"),
		},
		ProjectionExpression: aws.String("#N, version, created_at"),
		ConsistentRead:       aws.Bool(true),
	})
	if err != nil {
		return nil, err
	}

	if all {
		return decodeCredential(res.Items)
	}

	creds, err := decodeCredential(res.Items)
	if err != nil {
		return nil, err
	}

	return filterLatest(creds)
}

// GetAllSecrets returns a list of all secrets
func GetAllSecrets(all bool) ([]*DecryptedCredential, error) {
	log.Debug("Getting all secrets")

	res, err := dynamoSvc.Scan(&dynamodb.ScanInput{
		TableName: aws.String(Table),
		AttributesToGet: []*string{
			aws.String("name"),
			aws.String("version"),
			aws.String("key"),
			aws.String("contents"),
			aws.String("hmac"),
			aws.String("created_at"),
		},
		ConsistentRead: aws.Bool(true),
	})
	if err != nil {
		return nil, err
	}

	creds, err := decodeCredential(res.Items)
	if err != nil {
		return nil, err
	}

	var results []*DecryptedCredential

	for _, cred := range creds {

		dcred, err := decryptCredential(cred)
		if err != nil {
			return nil, err
		}

		results = append(results, dcred)
	}

	return results, nil
}

// PutSecret retrieve the secret from dynamodb
func PutSecret(alias, name, secret, version string) error {
	log.Debug("Putting secret")

	kmsKey := DefaultKmsKey

	if alias != "" {
		kmsKey = alias
	}

	if version == "" {
		version = "1"
	}

	dk, err := GenerateDataKey(kmsKey, 64)
	if err != nil {
		return err
	}

	dataKey := dk.Plaintext[:32]
	hmacKey := dk.Plaintext[32:]
	wrappedKey := dk.CiphertextBlob

	ctext, err := Encrypt(dataKey, []byte(secret))
	if err != nil {
		return err
	}

	b64hmac := ComputeHmac256(ctext, hmacKey)

	b64ctext := base64.StdEncoding.EncodeToString(ctext)

	cred := &Credential{
		Name:      name,
		Version:   version,
		Key:       base64.StdEncoding.EncodeToString(wrappedKey),
		Contents:  b64ctext,
		Hmac:      b64hmac,
		CreatedAt: time.Now().Unix(),
	}

	data, err := Encode(cred)

	if err != nil {
		return err
	}

	_, err = dynamoSvc.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(Table),
		Item:      data,
		ExpressionAttributeNames: map[string]*string{
			"#N": aws.String("name"),
		},
		ConditionExpression: aws.String("attribute_not_exists(#N)"),
	})

	return err
}

// DeleteSecret delete a secret
func DeleteSecret(name string) error {
	log.Debug("Deleting secret")

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

		log.WithFields(log.Fields{"name": cred.Name, "version": cred.Version}).Info("deleting")

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

// ResolveVersion calculate the version given a name and version
func ResolveVersion(name string, version int) (string, error) {
	log.Debug("Resolving version")

	if version != 0 {
		return strconv.Itoa(version), nil
	}

	ver, err := GetHighestVersion(name)
	if err != nil {
		if err == ErrSecretNotFound {
			return "1", nil
		}
		return "", err
	}

	if version, err = strconv.Atoi(ver); err != nil {
		return "", err
	}

	version++

	return strconv.Itoa(version), nil
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

func decodeCredential(items []map[string]*dynamodb.AttributeValue) ([]*Credential, error) {

	results := make([]*Credential, 0, len(items))

	for _, item := range items {
		cred := new(Credential)

		err := Decode(item, cred)
		if err != nil {
			return nil, err
		}

		results = append(results, cred)
	}
	return results, nil
}

func filterLatest(creds []*Credential) ([]*Credential, error) {

	sort.Sort(ByVersion(creds))

	names := map[string]*Credential{}

	for _, cred := range creds {
		names[cred.Name] = cred
	}

	results := make([]*Credential, 0, len(names))

	for _, val := range names {
		results = append(results, val)
	}

	// because maps key order is randomised in golang
	sort.Sort(ByVersion(results))

	return results, nil
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

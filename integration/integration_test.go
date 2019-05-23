// +build integration

package integration

import (
	"fmt"
	"os"
	"testing"

	"github.com/Versent/unicreds"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetHandler(cli.Default)
	log.SetLevel(log.DebugLevel)
}

func TestIntegrationGetSecret(t *testing.T) {

	var err error
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-west-2"
	}
	alias := os.Getenv("UNICREDS_KEY_ALIAS")
	if alias == "" {
		alias = "alias/unicreds"
	}
	tableName := os.Getenv("UNICREDS_TABLE_NAME")
	if tableName == "" {
		tableName = "credential-store"
	}

	unicreds.SetAwsConfig(aws.String(region), nil)

	encContext := unicreds.NewEncryptionContextValue()

	(*encContext)["test"] = aws.String("123")

	for i := 0; i < 15; i++ {
		err = unicreds.PutSecret(aws.String(tableName), alias, "Integration1", fmt.Sprintf("secret%d", i), unicreds.PaddedInt(i), encContext)
		if err != nil {
			log.Errorf("put err: %v", err)
		}

		assert.Nil(t, err)
	}

	cred, err := unicreds.GetHighestVersionSecret(aws.String(tableName), "Integration1", encContext)
	assert.Nil(t, err)
	assert.Equal(t, cred.Name, "Integration1")
	assert.Equal(t, cred.Secret, "secret14")
	assert.NotZero(t, cred.CreatedAt)
	assert.NotZero(t, cred.Version)

	for i := 0; i < 15; i++ {
		cred, err := unicreds.GetSecret(aws.String(tableName), "Integration1", unicreds.PaddedInt(i), encContext)
		assert.Nil(t, err)
		assert.Equal(t, cred.Name, "Integration1")
		assert.Equal(t, cred.Secret, fmt.Sprintf("secret%d", i))
		assert.NotZero(t, cred.CreatedAt)
		assert.Equal(t, cred.Version, unicreds.PaddedInt(i))
	}

	creds, err := unicreds.GetAllSecrets(aws.String(tableName), true)
	assert.Nil(t, err)
	assert.Len(t, creds, 15)

	// change the context and ensure this triggers an error
	(*encContext)["test"] = aws.String("345")

	cred, err = unicreds.GetHighestVersionSecret(aws.String(tableName), "Integration1", encContext)
	assert.Error(t, err)
	assert.Nil(t, cred)

	err = unicreds.DeleteSecret(aws.String(tableName), "Integration1")
	assert.Nil(t, err)
}

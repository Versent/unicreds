// +build integration

package integration

import (
	"fmt"
	"testing"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"
	"github.com/versent/unicreds"
)

func init() {
	log.SetHandler(cli.Default)
	log.SetLevel(log.DebugLevel)
}

func TestIntegrationGetSecret(t *testing.T) {

	var err error

	unicreds.SetRegion(aws.String("us-west-2"))

	encContext := unicreds.NewEncryptionContextValue()

	(*encContext)["test"] = aws.String("123")

	for i := 0; i < 15; i++ {
		err = unicreds.PutSecret(aws.String("credential-store"), "alias/accounting", "Integration1", "secret1", fmt.Sprintf("%d", i), encContext)
		if err != nil {
			log.Errorf("put err: %v", err)
		}

		assert.Nil(t, err)
	}

	cred, err := unicreds.GetSecret(aws.String("credential-store"), "Integration1", encContext)
	assert.Nil(t, err)
	assert.Equal(t, cred.Name, "Integration1")
	assert.Equal(t, cred.Secret, "secret1")
	assert.NotZero(t, cred.CreatedAt)
	assert.NotZero(t, cred.Version)

	creds, err := unicreds.GetAllSecrets(aws.String("credential-store"), true)
	assert.Nil(t, err)
	assert.Len(t, creds, 24)

	// change the context and ensure this triggers an error
	(*encContext)["test"] = aws.String("345")

	cred, err = unicreds.GetSecret(aws.String("credential-store"), "Integration1", encContext)
	assert.Error(t, err)
	assert.Nil(t, cred)

	err = unicreds.DeleteSecret(aws.String("credential-store"), "Integration1")
	assert.Nil(t, err)

}

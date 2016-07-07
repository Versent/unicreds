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

	for i := 0; i < 15; i++ {
		err = unicreds.PutSecret(aws.String("credential-store"), "alias/accounting", "Integration1", "secret1", fmt.Sprintf("%d", i))
		if err != nil {
			log.Errorf("put err: %v", err)
		}

		assert.Nil(t, err)
	}

	cred, err := unicreds.GetSecret(aws.String("credential-store"), "Integration1")
	assert.Nil(t, err)
	assert.Equal(t, cred.Name, "Integration1")
	assert.Equal(t, cred.Secret, "secret1")

	creds, err := unicreds.GetAllSecrets(aws.String("credential-store"), true)
	assert.Nil(t, err)
	assert.Len(t, creds, 24)

	err = unicreds.DeleteSecret(aws.String("credential-store"), "Integration1")
	assert.Nil(t, err)

}

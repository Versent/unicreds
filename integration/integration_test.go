// +build integration

package integration

import (
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

	unicreds.SetRegion(aws.String("us-west-2"))

	err := unicreds.PutSecret("alias/accounting", "Integration1", "secret1", "")
	assert.Nil(t, err)

	cred, err := unicreds.GetSecret("Integration1")
	assert.Nil(t, err)
	assert.Equal(t, cred.Name, "Integration1")
	assert.Equal(t, cred.Secret, "secret1")

	creds, err := unicreds.GetAllSecrets(true)
	assert.Nil(t, err)
	assert.Len(t, creds, 1)

	err = unicreds.DeleteSecret("Integration1")
	assert.Nil(t, err)

}

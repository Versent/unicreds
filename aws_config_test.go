package unicreds

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {

	err := SetAwsConfig(nil, nil, nil, nil)
	assert.Nil(t, err)

	err = SetAwsConfig(aws.String(""), aws.String(""), aws.String(""), aws.String(""))
	assert.Nil(t, err)

	err = SetAwsConfig(aws.String(""), aws.String("wolfeidau"), aws.String(""), aws.String(""))
	assert.Error(t, err)

	err = SetAwsConfig(aws.String("us-west-2"), aws.String("wolfeidau"), aws.String(""), aws.String(""))
	assert.Nil(t, err)

	err = SetAwsConfig(aws.String("us-west-2"), aws.String("wolfeidau"), aws.String("role"), aws.String("localstack"))
	assert.Nil(t, err)
}

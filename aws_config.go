package unicreds

import (
	"github.com/apex/log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

const (
	zoneURL = "http://169.254.169.254/latest/meta-data/placement/availability-zone"
)

// SetAwsConfig configure the AWS region with a fallback for discovery
// on EC2 hosts.
func SetAwsConfig(region, profile *string) error {
	if region == nil {
		// Try to get our region based on instance metadata
		region, err := getRegion()
		if err != nil {
			return err
		}
		// Update the aws config overrides if present
		setAwsConfig(region, profile)
		return nil
	}

	setAwsConfig(region, profile)
	return nil
}

func setAwsConfig(region, profile *string) {
	log.WithFields(log.Fields{"region": aws.StringValue(region), "profile": aws.StringValue(profile)}).Debug("Configure AWS")
	config := &aws.Config{Region: region}
	// if a profile is supplied then just use the shared credentials provider
	// as per docs this will look in $HOME/.aws/credentials if the filename is ""
	if aws.StringValue(profile) != "" {
		config.Credentials = credentials.NewSharedCredentials("", *profile)
	}
	SetDynamoDBConfig(config)
	SetKMSConfig(config)
}

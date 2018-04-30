package unicreds

import (
	"fmt"

	"github.com/apex/log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
)

const (
	zoneURL = "http://169.254.169.254/latest/meta-data/placement/availability-zone"
)

// SetAwsConfig configure the AWS region with a fallback for discovery
// on EC2 hosts.
func SetAwsConfig(region, profile *string, role *string) (err error) {
	if region == nil {
		// Try to get our region based on instance metadata
		region, err = getRegion()
		if err != nil {
			return err
		}
	}

	// This is to work around a limitation of the credentials
	// chain when providing an AWS profile as a flag
	if aws.StringValue(region) == "" && aws.StringValue(profile) != "" {
		return fmt.Errorf("Must provide a region flag when specifying a profile")
	}

	setAwsConfig(region, profile, role)
	return nil
}

func setAwsConfig(region, profile *string, role *string) {
	log.WithFields(log.Fields{"region": aws.StringValue(region), "profile": aws.StringValue(profile)}).Debug("Configure AWS")

	sess := getAwsSession(region, profile, role)

	SetDynamoDBSession(sess)
	SetKMSSession(sess)
}

func getAwsSession(region, profile, role *string) *session.Session {
	config := aws.Config{Region: region}

	// If a role is supplied, get credentials from STS Session
	if aws.StringValue(role) != "" {
		// If the role is being assumed through a non-default profile it must be added to the config.
		// When an empty string is provided, it uses the default credentials file.
		if aws.StringValue(profile) != "" {
			config.Credentials = credentials.NewSharedCredentials("", *profile)
		}

		sts_session := session.Must(session.NewSession(&config))
		log.WithFields(log.Fields{"role": aws.StringValue(role), "profile": aws.StringValue(profile)}).Debug("AssumeRole")
		config.Credentials = stscreds.NewCredentials(sts_session, *role)

		return session.Must(session.NewSession(&config))
	}

	// If no role is supplied, use the shared AWS config
	return session.Must(session.NewSessionWithOptions(session.Options{
		Config:            config,
		SharedConfigState: session.SharedConfigEnable,
		Profile:           *profile,
	}))
}

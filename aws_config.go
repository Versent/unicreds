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

	if aws.StringValue(region) == "" && aws.StringValue(profile) == "" {
		return nil
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
	config := aws.Config{Region: region}

	// if a profile is supplied then just use the shared credentials provider
	// as per docs this will look in $HOME/.aws/credentials if the filename is ""
	if aws.StringValue(profile) != "" {
		config.Credentials = credentials.NewSharedCredentials("", *profile)
	}

	// Are we assuming a role?
	if aws.StringValue(role) != "" {
		// Must request credentials from STS service and replace before passing on
		sts_sess := session.Must(session.NewSession(&config))
		log.WithFields(log.Fields{"role": aws.StringValue(role)}).Debug("AssumeRole")
		config.Credentials = stscreds.NewCredentials(sts_sess, *role)
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config:            config,
		SharedConfigState: session.SharedConfigEnable,
	}))

	SetDynamoDBSession(sess)
	SetKMSSession(sess)
}

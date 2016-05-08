package unicreds

import (
	"github.com/apex/log"
	"github.com/aws/aws-sdk-go/aws"
)

const (
	zoneURL = "http://169.254.169.254/latest/meta-data/placement/availability-zone"
)

// SetRegion configure the AWS region with a fallback for discovery
// on EC2 hosts.
func SetRegion(region *string) error {
	if region == nil {
		// Try to get our region based on instance metadata
		region, err := getRegion()
		if err != nil {
			return err
		}
		// Update the aws config overrides if present
		setRegion(region)
		return nil
	}

	setRegion(region)
	return nil
}

func setRegion(region *string) {
	log.WithField("region", *region).Debug("Setting region")
	SetDynamoDBConfig(&aws.Config{Region: region})
	SetKMSConfig(&aws.Config{Region: region})
}

package unicreds

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/kms/kmsiface"
)

var kmsSvc kmsiface.KMSAPI

func init() {
	kmsSvc = kms.New(session.New(), aws.NewConfig())
}

// SetKMSConfig override the default aws configuration
func setKMSConfig(config *aws.Config) {
	kmsSvc = kms.New(session.New(), config)
}

// DataKey which contains the details of the KMS key
type dataKey struct {
	CiphertextBlob []byte
	Plaintext      []byte
}

// GenerateDataKey simplified method for generating a datakey with kms
func generateDataKey(alias string, size int) (*dataKey, error) {

	numberOfBytes := int64(size)

	params := &kms.GenerateDataKeyInput{
		KeyId:             aws.String(alias),
		EncryptionContext: map[string]*string{},
		GrantTokens:       []*string{},
		NumberOfBytes:     aws.Int64(numberOfBytes),
	}

	resp, err := kmsSvc.GenerateDataKey(params)

	if err != nil {
		return nil, err
	}

	return &dataKey{
		CiphertextBlob: resp.CiphertextBlob,
		Plaintext:      resp.Plaintext, // return the plain text key after generation
	}, nil
}

// DecryptDataKey ask kms to decrypt the supplied data key
func decryptDataKey(ciphertext []byte) (*dataKey, error) {

	params := &kms.DecryptInput{
		CiphertextBlob:    ciphertext,
		EncryptionContext: map[string]*string{},
		GrantTokens:       []*string{},
	}
	resp, err := kmsSvc.Decrypt(params)

	if err != nil {
		return nil, err
	}

	return &dataKey{
		CiphertextBlob: ciphertext,
		Plaintext:      resp.Plaintext, // transfer the plain text key after decryption
	}, nil
}

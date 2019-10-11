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
func SetKMSConfig(config *aws.Config) {
	kmsSvc = kms.New(session.New(), config)
}

func SetKMSSession(sess *session.Session) {
	kmsSvc = kms.New(sess)
}

// DataKey which contains the details of the KMS key
type DataKey struct {
	CiphertextBlob []byte
	Plaintext      []byte
}

// GenerateDataKey simplified method for generating a datakey with kms
func GenerateDataKey(alias string, encContext *EncryptionContextValue, size int) (*DataKey, error) {

	numberOfBytes := int64(size)

	params := &kms.GenerateDataKeyInput{
		KeyId:             aws.String(alias),
		EncryptionContext: *encContext,
		GrantTokens:       []*string{},
		NumberOfBytes:     aws.Int64(numberOfBytes),
	}

	resp, err := kmsSvc.GenerateDataKey(params)

	if err != nil {
		return nil, err
	}

	return &DataKey{
		CiphertextBlob: resp.CiphertextBlob,
		Plaintext:      resp.Plaintext, // return the plain text key after generation
	}, nil
}

// DecryptDataKey ask kms to decrypt the supplied data key
func DecryptDataKey(ciphertext []byte, encContext *EncryptionContextValue) (*DataKey, error) {

	params := &kms.DecryptInput{
		CiphertextBlob:    ciphertext,
		EncryptionContext: *encContext,
		GrantTokens:       []*string{},
	}
	resp, err := kmsSvc.Decrypt(params)

	if err != nil {
		return nil, err
	}

	return &DataKey{
		CiphertextBlob: ciphertext,
		Plaintext:      resp.Plaintext, // transfer the plain text key after decryption
	}, nil
}

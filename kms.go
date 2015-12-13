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

// DataKey which contains the details of the KMS key
type DataKey struct {
	CiphertextBlob []byte
	Plaintext      []byte
}

// GenerateDataKey simplified method for generating a datakey with kms
func GenerateDataKey(alias string, size int) (*DataKey, error) {

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

	return &DataKey{
		CiphertextBlob: resp.CiphertextBlob,
		Plaintext:      resp.Plaintext, // return the plain text key after generation
	}, nil
}

// DecryptDataKey ask kms to decrypt the supplied data key
func DecryptDataKey(ciphertext []byte) (*DataKey, error) {

	params := &kms.DecryptInput{
		CiphertextBlob:    ciphertext,
		EncryptionContext: map[string]*string{},
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

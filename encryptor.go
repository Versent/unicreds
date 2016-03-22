package unicreds

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// Encrypt AES encryption method which matches the pycrypto package
// using CTR and AES256. Note this routine seeds the counter/iv with a value of 1
// then throws it away?!
func Encrypt(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, len(plaintext))

	initialCounter := newCounter()
	stream := cipher.NewCTR(block, initialCounter)
	stream.XORKeyStream(ciphertext, plaintext)

	return ciphertext, nil
}

// ComputeHmac256 compute a hmac256 signature of the supplied message and return
// the value hex encoded
func ComputeHmac256(message, secret []byte) string {
	h := hmac.New(sha256.New, secret)
	h.Write(message)
	return hex.EncodeToString(h.Sum(nil))
}

// Decrypt AES encryption method which matches the pycrypto package
// using CTR and AES256. Note this routine seeds the counter/iv with a value of 1
// then throws it away?!
func Decrypt(key, ciphertext []byte) ([]byte, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	initialCounter := newCounter()

	plaintext := ciphertext

	stream := cipher.NewCTR(block, initialCounter)

	stream.XORKeyStream(plaintext, ciphertext)

	return plaintext, nil
}

// start with a counter block with a default of 1 to be compatible with the python encryptor
// see https://pythonhosted.org/pycrypto/Crypto.Util.Counter-module.html for more info
func newCounter() []byte {
	return []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1}
}

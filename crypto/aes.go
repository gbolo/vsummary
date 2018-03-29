package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/spf13/viper"
)

// key to be used for all crypto operations
var key []byte

// ensure that key is properly initialized
func initKey() error {

	// if key is not a valid length, try to load it from config
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {

		key = []byte(viper.GetString("aes_key"))

		// check if key from config is proper size
		if len(key) != 16 && len(key) != 24 && len(key) != 32 {
			return fmt.Errorf("aes_key specified in config is not of correct length: %d", len(key))
		}

	}

	return nil

}

// encrypt string to base64 crypto using AES
func Encrypt(text string) (encryptedText string, err error) {

	// validate key
	if err = initKey(); err != nil {
		return
	}

	// create cipher block from key
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	plaintext := []byte(text)

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, rErr := io.ReadFull(rand.Reader, iv); rErr != nil {
		err = fmt.Errorf("iv ciphertext err: %s", rErr)
		return
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// convert to base64
	encryptedText = base64.URLEncoding.EncodeToString(ciphertext)

	return
}

// decrypt from base64 to decrypted string
func Decrypt(cryptoText string) (decryptedText string, err error) {

	// validate key
	if err = initKey(); err != nil {
		return
	}

	// create cipher block from key
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		err = fmt.Errorf("ciphertext is too small")
		return
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	decryptedText = fmt.Sprintf("%s", ciphertext)

	return
}

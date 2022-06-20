package lib

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"log"
	"os"
	"strings"
)

type Crypto interface {
	Encrypt(plainText string) (string, error)
	Decrypt(cipherIvKey string) (string, error)
}

type NiceCrypto struct {
	CipherKey   string
	CipherIvKey string
}

var SystemCipher Crypto

func CreateCipher(secretKey string) Crypto {
	var err error
	Cipher, err := NewNiceCrypto(secretKey, os.Getenv("cipher_iv_key"))
	if err != nil {
		log.Panic(err)
	}

	return Cipher
}

func (c NiceCrypto) Encrypt(plainText string) (string, error) {
	if strings.TrimSpace(plainText) == "" {
		return plainText, nil
	}

	block, err := aes.NewCipher([]byte(c.CipherKey))
	if err != nil {
		return "", err
	}

	encrypter := cipher.NewCBCEncrypter(block, []byte(c.CipherIvKey))
	paddedPlainText := padPKCS7([]byte(plainText), encrypter.BlockSize())

	cipherText := make([]byte, len(paddedPlainText))
	encrypter.CryptBlocks(cipherText, paddedPlainText)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func (c NiceCrypto) Decrypt(cipherText string) (string, error) {
	if strings.TrimSpace(cipherText) == "" {
		return cipherText, nil
	}

	decodedCipherText, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(c.CipherKey))
	if err != nil {
		return "", err
	}

	decrypt := cipher.NewCBCDecrypter(block, []byte(c.CipherIvKey))
	plainText := make([]byte, len(decodedCipherText))

	decrypt.CryptBlocks(plainText, decodedCipherText)
	trimmedPlainText := trimPKCS5(plainText)

	return string(trimmedPlainText), nil
}

func NewNiceCrypto(cipherKey, cipherIvKey string) (Crypto, error) {
	if ck := len(cipherKey); ck != 32 {
		return nil, aes.KeySizeError(ck)
	}

	if cik := len(cipherIvKey); cik != 16 {
		return nil, aes.KeySizeError(cik)
	}

	return &NiceCrypto{cipherKey, cipherIvKey}, nil
}

func padPKCS7(plainText []byte, blockSize int) []byte {
	padding := blockSize - len(plainText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(plainText, padText...)
}

func trimPKCS5(text []byte) []byte {
	padding := text[len(text)-1]
	return text[:len(text)-int(padding)]
}

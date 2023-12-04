package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

type AES struct {
}

func (a *AES) Encrypt(plainText, key []byte, uppercase bool) (cipherText, nonce string, err error) {
	if !(len(key) == 16 || len(key) == 24 || len(key) == 32) {
		return "", "", fmt.Errorf("invalid aes key size %d", len(key))
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", "", err
	}
	nonceData := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonceData); err != nil {
		return "", "", err
	}
	nonce = string(nonceData)
	cipherData := gcm.Seal(nil, nonceData, plainText, nil)
	cipherText = fmt.Sprintf("%02x", cipherData)
	if uppercase {
		cipherText = fmt.Sprintf("%02X", cipherData)
	}
	return
}

func (a *AES) Decrypt(cipherText, nonce, key []byte) (plainText string, err error) {
	if !(len(key) == 16 || len(key) == 24 || len(key) == 32) {
		return "", fmt.Errorf("invalid aes key size %d", len(key))
	}
	cipherData, err := hex.DecodeString(string(cipherText))
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	plainData, err := gcm.Open(nil, nonce, cipherData, nil)
	if err != nil {
		return "", err
	}
	plainText = string(plainData)
	return
}

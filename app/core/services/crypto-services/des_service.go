package crypto

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/hex"
	"fmt"
)

type DES struct {
}

func (d *DES) CBCEncrypt(plainText, key, iv []byte, uppercase bool) (cipherText string, err error) {
	if len(key) != 8 {
		return "", fmt.Errorf("invalid des key size %d", len(key))
	}
	if len(iv) != 8 {
		return "", fmt.Errorf("invalid des iv size %d", len(iv))
	}
	block, err := des.NewCipher(key)
	if err != nil {
		return "", err
	}
	moder := cipher.NewCBCEncrypter(block, iv)
	padding := block.BlockSize() - len(plainText)%block.BlockSize()
	paddedPlainText := append(plainText, bytes.Repeat([]byte{byte(padding)}, padding)...) //将要填充的大小添加在明文尾部,形成加密块大小的整数倍
	cipherData := make([]byte, len(paddedPlainText))
	moder.CryptBlocks(cipherData, paddedPlainText)
	cipherText = fmt.Sprintf("%02x", cipherData)
	if uppercase {
		cipherText = fmt.Sprintf("%02X", cipherData)
	}
	return
}

func (d *DES) CBCDecrypt(cipherText, key, iv []byte) (plainText string, err error) {
	if len(key) != 8 {
		return "", fmt.Errorf("invalid des key size %d", len(key))
	}
	if len(iv) != 8 {
		return "", fmt.Errorf("invalid des iv size %d", len(iv))
	}
	cipherData, err := hex.DecodeString(string(cipherText))
	if err != nil {
		return "", err
	}
	block, err := des.NewCipher(key)
	if err != nil {
		return "", err
	}
	plainData := make([]byte, len(cipherData))
	moder := cipher.NewCBCDecrypter(block, iv)
	moder.CryptBlocks(plainData, cipherData)
	unPadding := int(plainData[len(plainData)-1])    //从尾部获取添加的填充大小值
	plainData = plainData[:len(plainData)-unPadding] //去除填充值,从而获取到真正的明文
	plainText = string(plainData)
	return
}

func (d *DES) ECBEncrypt(plainText, key []byte, uppercase bool) (cipherText string, err error) {
	if len(key) != 8 {
		return "", fmt.Errorf("invalid des key size %d", len(key))
	}
	block, err := des.NewCipher(key)
	if err != nil {
		return "", err
	}
	padding := block.BlockSize() - len(plainText)%block.BlockSize()
	paddedPlainText := append(plainText, bytes.Repeat([]byte{byte(padding)}, padding)...) //将要填充的大小添加在明文尾部,形成加密块大小的整数倍
	cipherData := make([]byte, len(paddedPlainText))
	blockSize := block.BlockSize()
	blockCipherData := cipherData
	for len(paddedPlainText) > 0 {
		block.Encrypt(blockCipherData, paddedPlainText[:blockSize])
		paddedPlainText = paddedPlainText[blockSize:]
		blockCipherData = blockCipherData[blockSize:]
	}
	cipherText = fmt.Sprintf("%02x", cipherData)
	if uppercase {
		cipherText = fmt.Sprintf("%02X", cipherData)
	}
	return
}

func (d *DES) ECBDecrypt(cipherText, key []byte) (plainText string, err error) {
	if len(key) != 8 {
		return "", fmt.Errorf("invalid des key size %d", len(key))
	}
	cipherData, err := hex.DecodeString(string(cipherText))
	if err != nil {
		return "", err
	}
	block, err := des.NewCipher(key)
	if err != nil {
		return "", err
	}
	plainData := make([]byte, len(cipherData))
	blockSize := block.BlockSize()
	blockPlainData := plainData
	for len(cipherData) > 0 {
		block.Decrypt(blockPlainData, cipherData[:blockSize])
		cipherData = cipherData[blockSize:]
		blockPlainData = blockPlainData[blockSize:]
	}
	unPadding := int(plainData[len(plainData)-1])    //从尾部获取添加的填充大小值
	plainData = plainData[:len(plainData)-unPadding] //去除填充值,从而获取到真正的明文
	plainText = string(plainData)
	return
}

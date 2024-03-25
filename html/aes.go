package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

// PKCS7Padding 填充模式（加密数据块大小必须为aes.BlockSize的整数倍）
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS7UnPadding 去除填充模式
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// AesEncrypt AES加密
func AesEncrypt(origData, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	origData = PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	encodingData := base64.StdEncoding.EncodeToString(crypted)
	return encodingData, nil
}

// AesDecrypt AES解密
func AesDecrypt(crypted string, key []byte) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(crypted)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(decoded))
	blockMode.CryptBlocks(origData, decoded)
	origData = PKCS7UnPadding(origData)
	return origData, nil
}

func test() {
	// 需要加密的数据
	data := "Hello, World!"
	// AES的密钥，长度必须是16, 24或32字节，分别对应AES-128, AES-192, AES-256
	key := "1234567890123456"

	// 加密
	ciphertext, err := AesEncrypt([]byte(data), []byte(key))
	if err != nil {
		fmt.Println("Error encrypting:", err)
		return
	}
	fmt.Printf("Ciphertext: %s\n", ciphertext)

	// 解密
	plaintext, err := AesDecrypt(ciphertext, []byte(key))
	if err != nil {
		fmt.Println("Error decrypting:", err)
		return
	}
	fmt.Printf("Plaintext: %s\n", plaintext)
}

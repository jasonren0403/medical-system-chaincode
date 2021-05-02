package CryptoUtils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
)

type CryptMode int8

const (
	OFB        CryptMode = 0x1
	CBC        CryptMode = 0x2
	CTR        CryptMode = 0x3
	USE_BASE64 CryptMode = 0x10
	USE_HEX    CryptMode = 0x11
)

// DoAESEncrypt does an OFB style encryption to given text
// key[:blockSize] as IV
func DoAESEncrypt(key []byte, plaintext string, c CryptMode) (error, string) {
	plaindata := []byte(plaintext)

	block, err := aes.NewCipher(key)
	if err != nil {
		return errors.New(fmt.Sprintf("key 长度必须 16/24/32长度: %s", err.Error())), ""
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewOFB(block, key[:blockSize])
	crypted := make([]byte, len(plaindata))
	blockMode.XORKeyStream(crypted, plaindata)
	switch c ^ 0x10 {
	case 0:
		return nil, base64.RawURLEncoding.EncodeToString(crypted)
	case 1:
		return nil, hex.EncodeToString(crypted)
	default:
		return errors.New("cannot reach here, bad c value"), ""
	}
}

// DoAESDecrypt does an OFB style decryption to given text
// key[:blockSize] as IV
func DoAESDecrypt(key []byte, ciphertext string, c CryptMode) (error, string) {
	var cryptedByte []byte
	switch c ^ 0x10 {
	case 0:
		cryptedByte, _ = base64.RawURLEncoding.DecodeString(ciphertext)
	case 1:
		cryptedByte, _ = hex.DecodeString(ciphertext)
	default:
		return errors.New("cannot reach here, bad c value"), ""
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return errors.New(fmt.Sprintf("key 长度必须 16/24/32长度: %s", err.Error())), ""
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewOFB(block, key[:blockSize])
	orig := make([]byte, len(cryptedByte))
	blockMode.XORKeyStream(orig, cryptedByte)
	return nil, string(orig)
}

func PKCS7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

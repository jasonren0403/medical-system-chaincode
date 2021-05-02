package CryptoUtils

import (
	"encoding/hex"
)

func DoAESEncrypt(key []byte, plaintext string) []byte {
	return nil //todo
}

func DoAESDecrypt(key []byte, ciphertext []byte) string {
	return "" //todo
}

func toHexString(src []byte) string {
	return hex.EncodeToString(src)
}

func hexToByte(src string) []byte {
	decoded, err := hex.DecodeString(src)
	if err != nil {
		panic(err)
	}
	return decoded
}

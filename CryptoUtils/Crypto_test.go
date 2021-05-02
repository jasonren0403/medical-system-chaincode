package CryptoUtils

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

/* Crypto_test.go -- Test if crypto function work correctly */

func TestAESByB64(t *testing.T) {
	key := []byte("testKey1testKey1")
	plaintxt := "aesEncTest"
	err, val := DoAESEncrypt(key, plaintxt, USE_BASE64)
	if assert.NoError(t, err, "Encrypt should not cause error") {
		assert.NotEmpty(t, val, "cipher should not be empty")
		log.Printf("Enc value = %s", val)
	}
	log.Printf("Decrypt '%s' --> (key=%s)", val, key)
	err, val2 := DoAESDecrypt(key, "thbBaJFTT1LmGg", USE_BASE64)
	if assert.NoError(t, err, "Decrypt should not cause error") {
		assert.NotEmpty(t, val2, "cipher should not be empty")
		assert.EqualValues(t, plaintxt, val2, "Decrypt should be equal to original plaintext")
	}
}

func TestAESByHex(t *testing.T) {
	key := []byte("testKey1testKey1")
	plaintxt := "aesEncTest"
	err, val := DoAESEncrypt(key, plaintxt, USE_HEX)
	if assert.NoError(t, err, "Encrypt should not cause error") {
		assert.NotEmpty(t, val, "cipher should not be empty")
		log.Printf("Enc value = %s", val)
	}
	log.Printf("Decrypt '%s' --> (key=%s)", val, key)
	err, val2 := DoAESDecrypt(key, "b616c16891534f52e61a", USE_HEX)
	if assert.NoError(t, err, "Decrypt should not cause error") {
		assert.NotEmpty(t, val2, "cipher should not be empty")
		assert.EqualValues(t, plaintxt, val2, "Decrypt should be equal to original plaintext")
	}
}

func TestAESEdgeCase(t *testing.T) {
	key1 := []byte("Not Valid key length")
	err, _ := DoAESEncrypt(key1, "lllllllllllllllllllllllllllllllllllllll", USE_BASE64)
	assert.Error(t, err, "AES should only use key length 16/24/32, this key's length is", len(key1),
		"so it should not be received")
}

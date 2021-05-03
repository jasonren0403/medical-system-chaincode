package signUtils

import (
	"crypto/ed25519"
	"crypto/rand"
)

func Setup() (ed25519.PublicKey, ed25519.PrivateKey) {
	// returns (gpk,gmsk)
	pub, pri, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	return pub, pri
}

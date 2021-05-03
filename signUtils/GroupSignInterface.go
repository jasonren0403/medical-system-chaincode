package signUtils

import "crypto"

type IGroupSign interface {
	// Setup creates a pair of *group* key
	Setup() (crypto.PublicKey, crypto.PrivateKey)
	// Join returns a given member with private key and certificate, and the group manager get the manage key
	Join(Member string) (crypto.PrivateKey, string)
	// Sign outputs a signature for the given message and the member's private key
	Sign(key crypto.PrivateKey, message []byte) crypto.Hash
	// Verify uses the *group*'s public key to verify the validity of the given message
	Verify(key crypto.PublicKey, message []byte) bool
	// Open tracks the user certificate from the message with manage key
	Open(manageKey crypto.PrivateKey, message []byte) string
}

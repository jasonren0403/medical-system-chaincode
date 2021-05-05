package signUtils

import (
	"crypto/ed25519"
	"github.com/Nik-U/pbc"
)

// bbs04

type Group struct {
	g1, h, u, v, g2, w, ehw, ehg2, minusEg1g2 *pbc.Element
	pairing                                   *pbc.Pairing
}

type PrivateKey struct {
	*Group
	xi1, xi2, gamma *pbc.Element
}
type MemberKey struct {
	*Group
	x_, h_, u_ *pbc.Element
}
type Cert struct {
	*Group
	A, a *pbc.Element
}
type Sig struct {
	t1, t2, t3, c1, c2, c3, c, salpha, sbeta, sa, sx, sdelta1, sdelta2 *pbc.Element
}

type GroupInfo struct {
	// Gpk public key of the group
	Gpk ed25519.PublicKey `json:"pk_group"`
	// gamma private key of the group
	gamma ed25519.PrivateKey
	// gmsk "private" keys of the group manager (cert -> member's public key)
	gmsk map[string]ed25519.PublicKey
	// Members members of the group
	Members []GroupMember `json:"members"`
}

type GroupMember struct {
	// Cert identify the member
	Cert string `json:"cert"`
	// sk private key of the member
	sk ed25519.PrivateKey
	//sessionKey string
	// Role enum(manager|member)
	Role string `validate:"oneof=manager member" json:"role"`
}

type IGroupSign interface {
	// Setup creates a pair of *group* key
	Setup() *GroupInfo
	// Join returns a given member with private key and certificate, and the group manager get the manage key
	Join(initParam map[string]string, g GroupInfo) *GroupMember
	// Sign outputs a signature for the given message and the member's private key
	// returns string because crypto.Hash does not support crypto/ed25519 algorithm
	Sign(gm GroupMember, message string) string
	// Verify uses the *group*'s public key to verify the validity of the given message
	Verify(g GroupInfo, sig string, message string) bool
	// Open tracks the user certificate from the message with manage key
	Open(g GroupInfo, manageKey ed25519.PrivateKey, message string) string
	// Revoke updates gpk when group member revoked, changes the gpk as result
	Revoke(gm GroupMember) string
	// UpdateParams updates the group params after a member is revoked
	UpdateParams() bool
}

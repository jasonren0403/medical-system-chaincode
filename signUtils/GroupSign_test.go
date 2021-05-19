package signUtils

import (
	"github.com/Nik-U/pbc"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

/* GroupSign_test.go -- Test if signing algorithm work correctly */

func TestGroupSign(t *testing.T) {
	params := pbc.GenerateA(160, 512)
	pairing := params.NewPairing()
	g1 := pairing.NewG1().Rand()
	g2 := pairing.NewG2().Rand()
	priv := Setup(g1, g2, pairing)
	log.Println("Private key =", priv)
	//genarate  new  member
	member := priv.NewMember()
	//generate  new  cert
	cert := priv.Cert(member.u_)
	log.Println("cert.A=", cert.A.String())
	//verify  cert
	assert.True(t, member.VerifyCert(cert), "member verify should be success")
	//generate cipher
	c1 := pairing.NewG1().Rand()
	c2 := pairing.NewG1().Rand()
	c3 := pairing.NewG1().Mul(c1, c2)
	//generate signature
	sig := member.Sign(cert, c1, c2, c3)
	//verify  signature
	assert.True(t, priv.Group.Verify(sig, member.h_), "sig===member.h_")
	log.Println("sig-->", priv.Open(sig))
	assert.True(t, cert.A.Equals(priv.Open(sig)), "cert.A===open(sig)")
}

func BenchmarkGroupSignWithAddingMembers(b *testing.B) {
	params := pbc.GenerateA(160, 512)
	pairing := params.NewPairing()
	g1 := pairing.NewG1().Rand()
	g2 := pairing.NewG2().Rand()
	priv := Setup(g1, g2, pairing)
	b.Log("Init Group over")
	b.ReportAllocs()
	b.ResetTimer()
	b.Log("test times:", b.N)
	for i := 0; i < b.N; i++ {
		member := priv.NewMember()
		//log.Printf("Now there are %d member(s)\n", i+1)
		cert := priv.Cert(member.u_)
		//log.Println("member's cert:", cert.A)
		c1 := pairing.NewG1().Rand()
		c2 := pairing.NewG1().Rand()
		c3 := pairing.NewG1().Mul(c1, c2)
		//log.Printf("c1 len %d, c2 len %d, c3 len %d\n", c1.Len(), c2.Len(), c3.Len())
		sig := member.Sign(cert, c1, c2, c3)
		//log.Println("sig=", sig)
		priv.Group.Verify(sig, member.h_)
		priv.Open(sig)
	}
}

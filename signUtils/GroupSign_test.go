package signUtils

import (
	"fmt"
	"github.com/Nik-U/pbc"
	"testing"
)

/* GroupSign_test.go -- Test if signing algorithm work correctly */
// todo: build pbc library

func TestGroupSign(t *testing.T) {
	params := pbc.GenerateA(160, 512)
	pairing := params.NewPairing()
	g1 := pairing.NewG1().Rand()
	g2 := pairing.NewG2().Rand()
	priv := Setup(g1, g2, pairing)
	//genarate  new  member
	member := priv.NewMember()
	//generate  new  cert
	cert := priv.Cert(member.u_)
	fmt.Println("cert.A=", cert.A.String())

	member1 := priv.NewMember()
	cert1 := priv.Cert(member1.u_)
	fmt.Println("cert1.A=", cert1.A.String())

	//verify  cert
	member.VerifyCert(cert)
	member1.VerifyCert(cert1)
	//generate   mima
	c1 := pairing.NewG1().Rand()
	c2 := pairing.NewG1().Rand()
	c3 := pairing.NewG1().Mul(c1, c2)
	//generate  signature
	sig := member.Sign(cert, c1, c2, c3)
	sig1 := member1.Sign(cert1, c1, c2, c3)
	//verify    signature
	fmt.Println("sig===member.h_?", priv.Group.Verify(sig, member.h_))
	fmt.Println("sig1===member.h_?", priv.Group.Verify(sig1, member1.h_))
	fmt.Println("sig is opened to ==>", priv.Open(sig))
	fmt.Println("sig1 is opened to ==>", priv.Open(sig1))
}

func BenchmarkGroupSign(b *testing.B) {
	for i := 0; i < b.N; i++ {

	}
}

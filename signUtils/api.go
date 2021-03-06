package signUtils

import (
	"crypto/sha256"
	"github.com/Nik-U/pbc"
)

type SignAPI struct {
}

//BBS04: http://crypto.stanford.edu/~dabo/papers/groupsigs.pdf
//reference to yunfeiyanggzq's implementation: github.com/yunfeiyanggzq/BBS04_signature

func (member *MemberKey) VerifyCert(cert *Cert) bool {
	temp1 := member.pairing.NewG2().PowZn(member.g2, cert.a)
	temp2 := member.pairing.NewG2().Mul(member.w, temp1)
	e1 := member.pairing.NewGT().Pair(cert.A, temp2)
	temp1_ := member.pairing.NewG2().PowZn(member.g2, member.x_)
	e2 := member.pairing.NewGT().Pair(member.h_, temp1_)
	ttt1 := member.pairing.NewGT().Mul(e1, e2)
	ttt2 := member.pairing.NewGT().Pair(member.g1, member.g2)
	return ttt1.Equals(ttt2)
}

func Setup(g1, g2 *pbc.Element, pairing_ *pbc.Pairing) *PrivateKey {
	priv := new(PrivateKey)
	priv.Group = new(Group)
	priv.pairing = pairing_
	priv.g1 = g1
	// h1 \in G1
	priv.h = priv.pairing.NewG1().Rand()
	priv.g2 = g2
	// xi_1, xi_2 \in Zr
	priv.xi1 = priv.pairing.NewZr().Rand()
	priv.xi2 = priv.pairing.NewZr().Rand()
	temp1 := priv.pairing.NewZr().Invert(priv.xi1)
	temp2 := priv.pairing.NewZr().Invert(priv.xi2)
	//u^{\xi_1}=v^{\xi_2}=h u,v \in G1
	priv.u = priv.pairing.NewG1().PowZn(priv.h, temp1)
	priv.v = priv.pairing.NewG1().PowZn(priv.h, temp2)
	// gamma \in Zr
	priv.gamma = priv.pairing.NewZr().Rand()
	// w = g2^gamma
	priv.w = priv.pairing.NewG2().PowZn(priv.g2, priv.gamma)
	priv.precompute()
	return priv
}

func (g *Group) precompute() {
	g.ehw = g.pairing.NewGT().Pair(g.h, g.w)
	g.ehg2 = g.pairing.NewGT().Pair(g.h, g.g2)
	t := g.pairing.NewGT().Pair(g.g1, g.g2)
	g.minusEg1g2 = g.pairing.NewGT().Neg(t) //question
}

func (priv *PrivateKey) NewMember() *MemberKey {
	mem := new(MemberKey)
	mem.Group = priv.Group
	mem.x_ = mem.pairing.NewZr().Rand()
	mem.h_ = mem.pairing.NewG1().Rand()
	// A_i(u_) = g_1^{1/(\gamma+x_i)}
	mem.u_ = mem.pairing.NewG1().PowZn(mem.h_, mem.x_)
	return mem
}

func (priv *PrivateKey) Cert(u3 *pbc.Element) *Cert {
	cert := new(Cert)
	cert.Group = priv.Group
	cert.a = priv.pairing.NewZr().Rand()
	temp1 := priv.pairing.NewG1().Invert(u3)
	temp2 := priv.pairing.NewG1().Mul(temp1, priv.g1)
	temp3 := priv.pairing.NewZr().Add(priv.gamma, cert.a)
	temp4 := priv.pairing.NewZr().Invert(temp3)
	cert.A = priv.pairing.NewG1().PowZn(temp2, temp4)
	return cert
}

func (mem *MemberKey) Sign(cert *Cert, c1, c2, c3 *pbc.Element) *Sig {
	sig := new(Sig)
	alpha := mem.pairing.NewZr().Rand()
	beta := mem.pairing.NewZr().Rand()
	//T_1= u^\alpha
	t1 := mem.pairing.NewG1().PowZn(mem.u, alpha)
	//T_2= v^\beta
	t2 := mem.pairing.NewG1().PowZn(mem.v, beta)
	tmp := mem.pairing.NewZr().Add(alpha, beta)
	tmp1 := mem.pairing.NewG1().PowZn(mem.h, tmp)
	//T_3= Ah^{\alpha+\beta}
	t3 := mem.pairing.NewG1().Mul(cert.A, tmp1) //question
	//\delta_1=x\alpha
	delta1 := mem.pairing.NewZr().Mul(cert.a, alpha)
	//\delta_2=x\beta
	delta2 := mem.pairing.NewZr().Mul(cert.a, beta)
	ralpha := mem.pairing.NewZr().Rand()
	rbeta := mem.pairing.NewZr().Rand()
	rx := mem.pairing.NewZr().Rand()
	rdelta1 := mem.pairing.NewZr().Rand()
	rdelta2 := mem.pairing.NewZr().Rand()
	ra := mem.pairing.NewZr().Rand()
	//R_1= u^{r_\alpha}
	r1 := mem.pairing.NewG1().PowZn(mem.u, ralpha)
	//R_2= v^{r_\beta}
	r2 := mem.pairing.NewG1().PowZn(mem.v, rbeta)
	temp1 := mem.pairing.NewGT().Pair(t3, mem.g2)
	r3_e1 := mem.pairing.NewGT().PowZn(temp1, ra)
	uuu := mem.pairing.NewZr().Neg(ralpha)
	www := mem.pairing.NewZr().Neg(rbeta)
	xxx := mem.pairing.NewZr().Add(uuu, www)
	r3_e2 := mem.pairing.NewGT().PowZn(mem.ehw, xxx)
	uuu1 := mem.pairing.NewZr().Neg(rdelta1)
	www1 := mem.pairing.NewZr().Neg(rdelta2)
	xxx1 := mem.pairing.NewZr().Add(uuu1, www1)
	r3_e3 := mem.pairing.NewGT().PowZn(mem.ehg2, xxx1)
	eh3g2 := mem.pairing.NewGT().Pair(mem.h_, mem.g2)
	r3_e4 := mem.pairing.NewGT().PowZn(eh3g2, rx)
	//R_3= e\left(T_3,g_2\right)\ \cdot e\left(h,w\right)^{-r_\alpha-r_\beta}\ \cdot e\left(h,g_2\right)^{-r_{\delta_1}-r_{\delta_2}}
	r3 := mem.pairing.NewGT().Mul(mem.pairing.NewGT().Mul(r3_e1, r3_e2), mem.pairing.NewGT().Mul(r3_e3, r3_e4))
	tt_temp2 := mem.pairing.NewG1().PowZn(t1, ra)
	tt_temp := mem.pairing.NewZr().Neg(rdelta1)
	tt := mem.pairing.NewG1().PowZn(mem.u, tt_temp)
	//R_4= T_1^{r_x}\cdot u^{-r_{\delta_1}}
	r4 := mem.pairing.NewG1().Mul(tt, tt_temp2)
	rr_temp2 := mem.pairing.NewG1().PowZn(t2, ra)
	rr_temp := mem.pairing.NewZr().Neg(rdelta2)
	rr := mem.pairing.NewG1().PowZn(mem.v, rr_temp)
	//R_5= T_2^{r_x}\cdot v^{-r_{\delta_2}}
	r5 := mem.pairing.NewG1().Mul(rr, rr_temp2)
	var s string
	s += t1.String()
	s += t2.String()
	s += t3.String()
	s += r1.String()
	s += r2.String()
	s += r3.String()
	s += r4.String()
	s += r5.String()
	s += c1.String()
	s += c2.String()
	s += c3.String()
	//c= H(M,T_1,T_2,T_3,R_1,R_2,R_3,R_4,R_5) \in Z_p
	c := mem.pairing.NewZr().SetFromStringHash(s, sha256.New())
	sig.c1 = c1
	sig.c2 = c2
	sig.c3 = c3
	sig.c = c
	sig.t1 = t1
	sig.t2 = t2
	sig.t3 = t3
	//s_\alpha=r_\alpha+c\alpha
	sig.salpha = mem.pairing.NewZr().Add(ralpha, mem.pairing.NewZr().Mul(c, alpha))
	//s_\beta=r_\beta+c\beta
	sig.sbeta = mem.pairing.NewZr().Add(rbeta, mem.pairing.NewZr().Mul(c, beta))
	sig.sa = mem.pairing.NewZr().Add(ra, mem.pairing.NewZr().Mul(c, cert.a))
	//s_x=r_x+cx
	sig.sx = mem.pairing.NewZr().Add(rx, mem.pairing.NewZr().Mul(c, mem.x_))
	//s_{\delta_1}=r_{\delta_1}+c\delta_1
	sig.sdelta1 = mem.pairing.NewZr().Add(rdelta1, mem.pairing.NewZr().Mul(c, delta1))
	//s_{\delta_2}=r_{\delta_2}+c\delta_2
	sig.sdelta2 = mem.pairing.NewZr().Add(rdelta2, mem.pairing.NewZr().Mul(c, delta2))
	return sig
}

func (g *Group) Verify(sig *Sig, h3 *pbc.Element) bool {
	//r_1 = u^{s_\alpha}\cdot T_1^{-c}
	r1 := g.pairing.NewG1().Mul(g.pairing.NewG1().PowZn(g.u, sig.salpha), g.pairing.NewG1().PowZn(sig.t1, g.pairing.NewZr().Neg(sig.c)))
	//r_2 = v^{s_\beta}\cdot T_2^{-c}
	r2 := g.pairing.NewG1().Mul(g.pairing.NewG1().PowZn(g.v, sig.sbeta), g.pairing.NewG1().PowZn(sig.t2, g.pairing.NewZr().Neg(sig.c)))
	//******************************************
	//r_3 = e\left(T_3,g_2\right)^{s_x}\cdot e\left(h,w\right)^{-s_\alpha-s_\beta}\cdot e\left(h,g_2\right)^{-s_{\delta_1}-s_{\delta_2}}\cdot\left(e\left(T_3,w\right)/e\left(g_1,g_2\right)\right)^c
	temp1 := g.pairing.NewGT().Pair(sig.t3, g.g2)
	r3_e1 := g.pairing.NewGT().PowZn(temp1, sig.sa)
	uuu := g.pairing.NewZr().Neg(sig.salpha)
	www := g.pairing.NewZr().Neg(sig.sbeta)
	xxx := g.pairing.NewZr().Add(uuu, www)
	r3_e2 := g.pairing.NewGT().PowZn(g.ehw, xxx)
	uuu1 := g.pairing.NewZr().Neg(sig.sdelta1)
	www1 := g.pairing.NewZr().Neg(sig.sdelta2)
	xxx1 := g.pairing.NewZr().Add(uuu1, www1)
	r3_e3 := g.pairing.NewGT().PowZn(g.ehg2, xxx1)
	eh3g2 := g.pairing.NewGT().Pair(h3, g.g2)
	r3_e4 := g.pairing.NewGT().PowZn(eh3g2, sig.sx)
	r3_tep := g.pairing.NewGT().Mul(g.pairing.NewGT().Mul(r3_e1, r3_e2), g.pairing.NewGT().Mul(r3_e3, r3_e4))
	yyy := g.pairing.NewGT().Pair(sig.t3, g.w)
	ggg := g.pairing.NewGT().Pair(g.g1, g.g2)
	hhh := g.pairing.NewGT().Invert(ggg)
	r3 := g.pairing.NewGT().Mul(r3_tep, g.pairing.NewGT().PowZn(g.pairing.NewGT().Mul(yyy, hhh), sig.c))
	//*******************************************
	tt_temp2 := g.pairing.NewG1().PowZn(sig.t1, sig.sa)
	tt_temp := g.pairing.NewZr().Neg(sig.sdelta1)
	tt := g.pairing.NewG1().PowZn(g.u, tt_temp)
	//r_4 = T_1^{s_x}\cdot u^{-s_{\delta_1}}
	r4 := g.pairing.NewG1().Mul(tt, tt_temp2)
	rr_temp2 := g.pairing.NewG1().PowZn(sig.t2, sig.sa)
	rr_temp := g.pairing.NewZr().Neg(sig.sdelta2)
	rr := g.pairing.NewG1().PowZn(g.v, rr_temp)
	//r_5 = T_2^{s_x}\cdot u^{-s_{\delta_2}}
	r5 := g.pairing.NewG1().Mul(rr, rr_temp2)
	var s string
	s += sig.t1.String()
	s += sig.t2.String()
	s += sig.t3.String()
	s += r1.String()
	s += r2.String()
	s += r3.String()
	s += r4.String()
	s += r5.String()
	s += sig.c1.String()
	s += sig.c2.String()
	s += sig.c3.String()
	//c=?H(M,T_1,T_2,T_3,\widetilde{R_1},\widetilde{R_2},\widetilde{R_3},\widetilde{R_4},\widetilde{R_5})
	c_ := g.pairing.NewZr().SetFromStringHash(s, sha256.New())
	return c_.Equals(sig.c)
}

func (priv *PrivateKey) Open(sig *Sig) *pbc.Element {
	//T_3/ (T_1^{\xi_1}\cdot T_2^{\xi_2})
	temp1 := priv.pairing.NewG1().PowZn(sig.t1, priv.xi1)
	temp2 := priv.pairing.NewG1().PowZn(sig.t2, priv.xi2)
	temp3 := priv.pairing.NewG1().Mul(sig.t3, priv.pairing.NewG1().Invert(priv.pairing.NewG1().Mul(temp1, temp2)))
	return temp3
}

func UpdateParams() bool {
	return false
}


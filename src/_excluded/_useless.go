package _excluded

type SimpleChaincode struct {
}
func (r *SignAPI) Setup() *GroupInfo {
	pub, pri, err := ed25519.GenerateKey(rand.Reader)
	member := GroupMember{Cert: utils.NewGenerator(16, "string").RandStr(), Role: "manager"}
	if err != nil {
		panic(err)
	}
	gi := &GroupInfo{
		Gpk:     pub,
		gamma:   pri,
		Members: []GroupMember{member},
	}
	return gi
}

func (r *SignAPI) Join(initParam map[string]string, g GroupInfo) *GroupMember {
	pub, pri, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	var certificate string
	if v, ok := initParam["ID"]; ok {
		certificate = v
	} else {
		certificate = utils.NewGenerator(16, "string").RandStr()
	}

	// Add the member into the group
	member := GroupMember{
		Cert: certificate,
		sk:   pri,
		Role: "member",
	}
	g.Members = append(g.Members, member)
	// Add its management key to group manager
	g.gmsk[certificate] = pub
	return &member
}
func (r *SignAPI) Sign(gm GroupMember, message string) string {
	gmpk := gm.sk
	signature := ed25519.Sign(gmpk, utils.Str2bytes(message))
	return utils.Bytes2str(signature)
}

func (r *SignAPI) Verify(g GroupInfo, sig string, message string) bool {
	gpk := g.GetPubKey()
	return ed25519.Verify(gpk, utils.Str2bytes(message), utils.Str2bytes(sig))
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

func (g *GroupInfo) GetMembers() []GroupMember {
	return g.Members
}

func (g *GroupInfo) GetPubKey() ed25519.PublicKey {
	return g.Gpk
}

func (g *GroupInfo) DumpPrivateKey() {
	log.Println("Group priv =", hex.EncodeToString(g.gamma))
	for k, v := range g.gmsk {
		log.Println("member ", k, "=>", hex.EncodeToString(v))
	}
}
//func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
//	//实现链码初始化或升级时的处理逻辑
//	//编写时可灵活使用stub中的API
//	log.Println("===Init()===")
//	return shim.Success(nil)
//}
//
//func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
//	log.Println("===Invoke()===")
//	//实现链码运行中被调用或查询时的处理逻辑
//	//GetArgs() [][]byte 以byte数组的数组的形式获得传入的参数列表
//	//GetStringArgs() []string 以字符串数组的形式获得传入的参数列表
//	//GetFunctionAndParameters() (string, []string) 将字符串数组的参数分为两部分，数组第一个字是Function，剩下的都是Parameter
//	//GetArgsSlice() ([]byte, error) 以byte切片的形式获得参数列表
//	function,args := stub.GetFunctionAndParameters()
//	log.Println("Invoke() is running ",function)
//	log.Println("Args are ",args)
//	/*
//		switch function{
//		case
//		}
//	*/
//	return shim.Success(nil)
//}

package _excluded

type SimpleChaincode struct {
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

package main

import (
	smartMedicineSystem "ccode/src"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	_ "github.com/hyperledger/fabric-contract-api-go/contractapi/utils"
	"log"
)

/* https://github.com/hyperledger/fabric-samples/tree/release-1.4/chaincode-docker-devmode */
//docker-compose -f docker-compose-simple.yaml up
//docker exec -it chaincode sh
//cd chaincode_example02/go
//go build -o chaincode_example02
//CORE_PEER_ADDRESS=peer:7052 CORE_CHAINCODE_ID_NAME=mycc:0 ./chaincode_example02
//docker exec -it cli bash
//go mod vendor
//go build
func main() {
	defer func() {
		e := recover()
		if e != nil {
			log.Print("Error=", e)
		}
	}()

	smartMedicineContract := new(smartMedicineSystem.MedicalSystem)
	//smartMedicineContract.UnknownTransaction = smartMedicineSystem.UnknownTransactionHandler
	//smartMedicineContract.BeforeTransaction = utils.UndefinedInterface{}

	cc, err := contractapi.NewChaincode(smartMedicineContract)
	if err != nil {
		log.Panicf("Error creating smart medicine system chaincode: %v", err)
	}
	cc.DefaultContract = smartMedicineContract.GetName()

	if err := cc.Start(); err != nil {
		log.Panicf("Error starting smart medicine system chaincode: %v", err)
	}

}

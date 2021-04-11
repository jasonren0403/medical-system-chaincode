package main

import (
	_ "github.com/hyperledger/fabric-contract-api-go/contractapi/utils"
	"log"
)

//fabric-samples -> chaincode-docker-devmode
//docker-compose -f docker-compose-simple.yaml up
//docker exec -it chaincode sh
//go mod vendor
//go build
func main() {
	defer func() {
		e := recover()
		if e != nil {
			log.Print("Error=", e)
		}
	}()

	/*
		smartMedicineContract := new(smartMedicineSystem.MedicalSystem)
		smartMedicineContract.UnknownTransaction = smartMedicineSystem.UnknownTransactionHandler
		//smartMedicineContract.BeforeTransaction = utils.UndefinedInterface{}

		cc, err := contractapi.NewChaincode(smartMedicineContract)
		if err != nil {
			log.Panicf("Error creating smart medicine system chaincode: %v", err)
		}
		cc.DefaultContract = smartMedicineContract.GetName()

		if err := cc.Start(); err != nil {
			log.Panicf("Error starting smart medicine system chaincode: %v", err)
		}
	*/
}

package test

import (
	"ccode/src"
	_ "fmt"
	"testing"

	// Use these two package for testing and mocking contract environment
	_ "github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
)

//func checkXXX(t *testing.T, stub *shimtest.MockStub, args...)
//stub.mockInit(string,args)
//shimtest.newMockStub("",contract)
//stub.mockInvoke(string,args)
func TestInitLedger(t *testing.T) {
	cc := new(smartMedicineSystem.MedicalSystem)
	stub := shimtest.NewMockStub("smartMedicineSystem", cc)
	stub.MockInit("1", nil)
}

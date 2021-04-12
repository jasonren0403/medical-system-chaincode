package smartMedicineSystem

import (
	"ccode/src"
	_ "fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"

	// Use these two package for testing and mocking contract environment
	_ "github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
)

/*
 *	smartMedicineSystem_test.go Testing chaincode(with unit tests)
 *	https://www.cnblogs.com/skzxc/p/12150476.html
 */

//func checkXXX(t *testing.T, stub *shimtest.MockStub, args...)
//stub.mockInit(string,args)
//shimtest.newMockStub("",contract)
//stub.mockInvoke(string,args)
func TestInitLedger(t *testing.T) {
	cc := new(smartMedicineSystem.MedicalSystem)
	stub := shimtest.NewMockStub("smartMedicineSystem", cc)
	stub.MockInit("1", nil)
	log.Println("name=", stub.Name)
	assert.NotNil(t, stub.Name, "Stub's name is nil!")
}

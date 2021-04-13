package smartMedicineSystem

import (
	"ccode/src"
	"ccode/src/asset"
	"encoding/json"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/stretchr/testify/assert"
	"log"
	"strconv"
	"testing"

	// Use these two package for testing and mocking contract environment
	_ "github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
)

/*
 *	smartMedicineSystem_test.go Testing chaincode(with unit tests)
 *	https://www.cnblogs.com/skzxc/p/12150476.html
 */
const (
	test_UUID     = "1"
	internal_name = "smartMedicineSystem"
)

//func checkXXX(t *testing.T, stub *shimtest.MockStub, args...)
//stub.mockInit(string,args)
//shimtest.newMockStub("",contract)
//[][]byte{[]byte("set"), []byte("a"), []byte("100")}
//stub.mockInvoke(string,args)
func TestInitLedger(t *testing.T) {
	cc := new(smartMedicineSystem.MedicalSystem)
	stub := shimtest.NewMockStub(internal_name, cc)
	result := stub.MockInit(test_UUID, nil)
	if result.Status != shim.OK {
		log.Println("result status is not OK but " + strconv.Itoa(int(result.Status)))
		log.Fatalln(result)
	}
	assert.NotNil(t, stub.Name, "Stub's name is nil!")
	assert.EqualValues(t, internal_name, stub.Name, "Stub's name is incorrect!")
}

func TestPatientInfoGet(t *testing.T) {
	patientID := "p1"
	cc := new(smartMedicineSystem.MedicalSystem)
	stub := shimtest.NewMockStub(internal_name, cc)
	stub.MockInit(test_UUID, nil)
	res := stub.MockInvoke(test_UUID, [][]byte{[]byte("GetPatientInfoByPID"), []byte(patientID)})
	// todo:change payload return and use assert.Same for testing
	log.Println("The returning value of getPatientInfobyPid(p1) is", string(res.Payload))
}

func TestIsValidDoctor(t *testing.T) {
	doctor := asset.Doctor{
		Person: asset.Person{
			ID: "d1", Name: "Apple", Age: 24,
		},
		Department: "Dep1",
	}
	binDoctor, _ := json.Marshal(doctor)

	cc := new(smartMedicineSystem.MedicalSystem)
	stub := shimtest.NewMockStub(internal_name, cc)
	stub.MockInit(test_UUID, nil)
	res := stub.MockInvoke(test_UUID, [][]byte{[]byte("isValidDoctor"), binDoctor})
	assert.EqualValues(t, "success", string(res.Payload), "It should be a valid doctor")
}

func TestGetMedicalRecord(t *testing.T) {
	patientID := "p1"
	cc := new(smartMedicineSystem.MedicalSystem)
	stub := shimtest.NewMockStub(internal_name, cc)
	stub.MockInit(test_UUID, nil)
	res := stub.MockInvoke(test_UUID, [][]byte{[]byte("GetMedicalRecord"), []byte(patientID)})
	// todo:change payload return and use assert.Same for testing
	log.Println(string(res.Payload))
}

func TestPatientInfoSet(t *testing.T) {
	patientID := "p1"
	newInfo := make(map[string]interface{})
	newInfo["isMarried"] = false
	binInfo, _ := json.Marshal(newInfo)
	cc := new(smartMedicineSystem.MedicalSystem)
	stub := shimtest.NewMockStub(internal_name, cc)
	stub.MockInit(test_UUID, nil)
	res := stub.MockInvoke(test_UUID, [][]byte{[]byte("SetPatientInfo"), []byte(patientID), binInfo})
	// todo:change payload return and use assert.Same for testing
	log.Println(string(res.Payload))
}

package smartMedicineSystem

import (
	"ccode/src"
	"ccode/src/asset"
	"encoding/json"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
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

var stub *shimtest.MockStub

//func checkXXX(t *testing.T, stub *shimtest.MockStub, args...)
//stub.mockInit(string,args)
//shimtest.newMockStub("",contract)
//[][]byte{[]byte("set"), []byte("a"), []byte("100")}
//stub.mockInvoke(string,args)
func TestMain(m *testing.M) {
	log.Println("===test main===")
	setup()
	exitcode := m.Run() // run all cases
	tearDown()
	os.Exit(exitcode)
}

func setup() {
	log.Println("===setup===")
	cc := new(smartMedicineSystem.MedicalSystem)
	stub = shimtest.NewMockStub(internal_name, cc)
}

func tearDown() {
	log.Println("===tearDown===")
}

func TestInitLedger(t *testing.T) {
	assert.FileExists(t, "../../init.json", "Init file does not exist!")
	result := stub.MockInit(test_UUID, nil)
	assert.EqualValuesf(t, shim.OK, result.Status, "Result status is not OK, get %d", result.Status)
	assert.NotNil(t, stub.Name, "Stub's name is nil!")
	assert.EqualValues(t, internal_name, stub.Name, "Stub's name is incorrect!")
}

func TestPatientInfoGet(t *testing.T) {
	patientID := "p1"
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
	stub.MockInit(test_UUID, nil)
	res := stub.MockInvoke(test_UUID, [][]byte{[]byte("isValidDoctor"), binDoctor})
	assert.EqualValues(t, "success", string(res.Payload), "It should be a valid doctor")
}

func TestGetMedicalRecord(t *testing.T) {
	patientID := "p1"
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
	stub.MockInit(test_UUID, nil)
	res := stub.MockInvoke(test_UUID, [][]byte{[]byte("SetPatientInfo"), []byte(patientID), binInfo})
	// todo:change payload return and use assert.Same for testing
	log.Println(string(res.Payload))
}

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

func TestInitNewRecord(t *testing.T) {
	patientID := "p1"
	rContent := map[string]interface{}{
		"key1": "value1",
	}
	nRecord := asset.Record{
		Type:    "test2",
		Time:    "2021-4-14 9:45:11",
		Content: rContent,
		Signature: asset.Doctor{
			Person: asset.Person{
				ID: "d1", Name: "Apple", Age: 24,
			},
			Department: "Dep1",
		},
	}
	brContent, _ := json.Marshal(rContent)
	bnSign, _ := json.Marshal(nRecord.Signature)
	stub.MockInit(test_UUID, nil)
	res := stub.MockInvoke(test_UUID, [][]byte{[]byte("InitNewRecord"), []byte(patientID), []byte(nRecord.Type),
		[]byte(nRecord.Time), brContent, bnSign})
	var records []asset.Record
	err := json.Unmarshal(res.Payload, &records)
	assert.Nil(t, err, "No problem should appear unmarshalling")
	assert.Len(t, records, 2, "There should be 2 records of patient ", patientID)
	assert.Contains(t, records, nRecord, "The new record should be inserted")
}

func TestPatientInfoGet(t *testing.T) {
	patientID := "p1"
	stub.MockInit(test_UUID, nil)
	res := stub.MockInvoke(test_UUID, [][]byte{[]byte("GetPatientInfoByPID"), []byte(patientID)})
	var pInfo asset.OutPatient
	err := json.Unmarshal(res.Payload, &pInfo)
	assert.Nil(t, err, "there should be no problem unmarshalling returning payload ")
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
	res := stub.MockInvoke(test_UUID, [][]byte{[]byte("IsValidDoctor"), binDoctor})
	assert.EqualValues(t, "true", string(res.Payload), "It should be a valid doctor")
	notExistDoctor := asset.Doctor{
		Person: asset.Person{
			ID: "d1", Name: "Apple", Age: 24,
		},
		Department: "Dep2",
	}
	binDoctor, _ = json.Marshal(notExistDoctor)
	res = stub.MockInvoke(test_UUID, [][]byte{[]byte("IsValidDoctor"), binDoctor})
	assert.EqualValues(t, "false", res.Payload, "Though dep is different, it should return false")
}

func TestGetMedicalRecord(t *testing.T) {
	patientID := "p1"
	stub.MockInit(test_UUID, nil)
	resErr := stub.MockInvoke(test_UUID, [][]byte{[]byte("GetMedicalRecord")})
	assert.EqualValues(t, 500, resErr.Status, "Error message should appear on 0 params")
	assert.Contains(t, resErr.Message, "Support a pid(string) for this call",
		"Following message should be displayed")
	res := stub.MockInvoke(test_UUID, [][]byte{[]byte("GetMedicalRecord"), []byte(patientID)})
	var rec []asset.Record
	err := json.Unmarshal(res.Payload, &rec)
	assert.Nil(t, err, "Error is not nil! Error is ", err)
	assert.Len(t, rec, 1, "There should be 1 record, found ", len(rec))
}

func TestPatientInfoSet(t *testing.T) {
	var pInfo asset.OutPatient
	patientID := "p1"
	newInfo := make(map[string]interface{})
	newInfo["isMarried"] = false
	binInfo, _ := json.Marshal(newInfo)
	stub.MockInit(test_UUID, nil)
	res := stub.MockInvoke(test_UUID, [][]byte{[]byte("SetPatientInfo"), []byte(patientID), binInfo})
	err := json.Unmarshal(res.Payload, &pInfo)
	assert.Nil(t, err, "Nothing wrong happens to unmarshalling")
	assert.False(t, pInfo.IsMarried, "This has successfully changed")
}

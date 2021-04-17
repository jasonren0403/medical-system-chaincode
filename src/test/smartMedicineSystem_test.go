package smartMedicineSystem

import (
	"bytes"
	"ccode/src"
	"ccode/src/asset"
	"encoding/json"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
	"time"

	// Use these two package for testing and mocking contract environment
	"github.com/hyperledger/fabric-chaincode-go/shim"
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

// -- Helpers -- //
func AddRecord(pid string, rec asset.Record, rContent map[string]interface{}) peer.Response {
	var brContent []byte
	if rContent == nil {
		brContent, _ = json.Marshal(rec.Content)
	} else {
		brContent, _ = json.Marshal(rContent)
	}
	bnSign, _ := json.Marshal(rec.Signature)
	return stub.MockInvoke(test_UUID, [][]byte{[]byte("InitNewRecord"), []byte(pid), []byte(rec.Type),
		[]byte(rec.Time), brContent, bnSign})
}

func GetRecord(pid string) peer.Response {
	if len(pid) == 0 {
		return stub.MockInvoke(test_UUID, [][]byte{[]byte("GetMedicalRecord")})
	}
	return stub.MockInvoke(test_UUID, [][]byte{[]byte("GetMedicalRecord"), []byte(pid)})
}

func GetPatientInfo(pid string) peer.Response {
	return stub.MockInvoke(test_UUID, [][]byte{[]byte("GetPatientInfoByPID"), []byte(pid)})
}

func SetPatientInfo(pid string, newInfo map[string]interface{}) peer.Response {
	binInfo, _ := json.Marshal(newInfo)
	return stub.MockInvoke(test_UUID, [][]byte{[]byte("SetPatientInfo"), []byte(pid), binInfo})
}

func IsValidDoctor(doctor asset.Doctor) peer.Response {
	binDoctor, _ := json.Marshal(doctor)
	return stub.MockInvoke(test_UUID, [][]byte{[]byte("IsValidDoctor"), binDoctor})
}

// -- Tests -- //
func TestInitLedger(t *testing.T) {
	assert.FileExists(t, "../../init.json", "Init file does not exist!")
	result := stub.MockInit(test_UUID, nil)
	assert.EqualValuesf(t, shim.OK, result.Status, "Result status is not OK, get %d", result.Status)
	assert.NotNil(t, stub.Name, "Stub's name is nil!")
	assert.EqualValues(t, internal_name, stub.Name, "Stub's name is incorrect!")
}

func TestInitNewRecord(t *testing.T) {
	patientID := "p1"
	cnt := map[string]interface{}{
		"keystr":  "value1",
		"keybool": true,
	}
	nRecord := asset.Record{
		Type:    "test2",
		Time:    "2021-4-14 9:45:11",
		Content: cnt,
		Signature: asset.Doctor{
			Person: asset.Person{
				ID: "d1", Name: "Apple", Age: 24,
			},
			Department: "Dep1",
		},
	}
	stub.MockInit(test_UUID, nil)
	res := AddRecord(patientID, nRecord, nil)
	var records []asset.Record
	dec := json.NewDecoder(bytes.NewBuffer(res.Payload))
	dec.UseNumber()
	err := dec.Decode(&records)
	assert.NoError(t, err, "No problem should appear unmarshalling")
	assert.Len(t, records, 2, "There should be 2 records of patient ", patientID)
	assert.Contains(t, records, nRecord, "The new record should be inserted")
	// another one
	res = AddRecord(patientID, asset.Record{
		Type: "test3",
		Time: time.Now().Format("2006-1-2 15:04:05"),
		Signature: asset.Doctor{
			Person: asset.Person{
				ID: "d1", Name: "Apple", Age: 24,
			},
			Department: "Dep1",
		},
	}, map[string]interface{}{
		"keybool": false,
		"keynum":  67,
	})
	dec = json.NewDecoder(bytes.NewBuffer(res.Payload))
	dec.UseNumber()
	err = dec.Decode(&records)
	assert.NoError(t, err, "No problem should appear unmarshalling")
	assert.Len(t, records, 3, "There should be 3 records of patient ", patientID)
}

func TestPatientInfoGet(t *testing.T) {
	patientID := "p1"
	stub.MockInit(test_UUID, nil)
	res := GetPatientInfo(patientID)
	var pInfo asset.OutPatient
	err := json.Unmarshal(res.Payload, &pInfo)
	assert.NoError(t, err, "there should be no problem unmarshalling returning payload ")
}

func TestIsValidDoctor(t *testing.T) {
	doctor := asset.Doctor{
		Person: asset.Person{
			ID: "d1", Name: "Apple", Age: 24,
		},
		Department: "Dep1",
	}
	notExistDoctor := asset.Doctor{
		Person: asset.Person{
			ID: "d1", Name: "Apple", Age: 24,
		},
		Department: "Dep2",
	}
	stub.MockInit(test_UUID, nil)
	res := IsValidDoctor(doctor)
	assert.EqualValues(t, "true", string(res.Payload), "It should be a valid doctor")
	res = IsValidDoctor(notExistDoctor)
	assert.EqualValues(t, "false", res.Payload, "Though dep is different, it should return false")
}

func TestGetMedicalRecord(t *testing.T) {
	patientID := "p1"
	stub.MockInit(test_UUID, nil)
	resErr := GetRecord("")
	assert.EqualValues(t, 500, resErr.Status, "Error message should appear on 0 params")
	assert.Contains(t, resErr.Message, "Support a pid(string) for this call",
		"Following message should be displayed")
	res := GetRecord(patientID)
	var rec []asset.Record
	err := json.Unmarshal(res.Payload, &rec)
	assert.NoError(t, err, "Error is not nil! Error is ", err)
	assert.Len(t, rec, 1, "There should be 1 record, found ", len(rec))
}

func TestPatientInfoSet(t *testing.T) {
	var pInfo asset.OutPatient
	patientID := "p1"
	stub.MockInit(test_UUID, nil)
	res := SetPatientInfo(patientID, map[string]interface{}{
		"isMarried": false,
	})
	err := json.Unmarshal(res.Payload, &pInfo)
	assert.NoError(t, err, "Nothing wrong happens to unmarshalling")
	assert.False(t, pInfo.IsMarried, "This has successfully changed")
}

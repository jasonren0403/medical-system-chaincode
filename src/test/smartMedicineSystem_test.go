package smartMedicineSystem

import (
	"bytes"
	"ccode/src"
	"ccode/src/asset"
	"ccode/src/utils"
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
	PRINTRES      = true
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
	_ = stub.MockInit(test_UUID, nil)
}

func tearDown() {
	log.Println("===tearDown===")
}

// -- Helpers -- //

func NewPatient(pat asset.OutPatient) peer.Response {
	bin, _ := json.Marshal(pat)
	return stub.MockInvoke(test_UUID, [][]byte{[]byte("NewPatient"), bin})
}

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

func GetAllPatients() peer.Response {
	return stub.MockInvoke(test_UUID, [][]byte{[]byte("GetAllPatients")})
}

func IsValidDoctor(doctor asset.Doctor) peer.Response {
	binDoctor, _ := json.Marshal(doctor)
	return stub.MockInvoke(test_UUID, [][]byte{[]byte("IsValidDoctor"), binDoctor})
}

func QueryDoctorByID(dID string) peer.Response {
	return stub.MockInvoke(test_UUID, [][]byte{[]byte("QueryDoctorByID"), []byte(dID)})
}

func QueryMRByStartToEndDate(startdate time.Time, enddate time.Time, pid string) peer.Response {
	return stub.MockInvoke(test_UUID, [][]byte{[]byte("GetMRByDate"), []byte(pid),
		[]byte(startdate.Format("2006-1-2 15:01:05")),
		[]byte(enddate.Format("2006-1-2 15:01:05"))})
}

func GetAllDoctors() peer.Response {
	return stub.MockInvoke(test_UUID, [][]byte{[]byte("GetAllDoctors")})
}

func AddCollaborator(member asset.Doctor, pID string,
	filterRec map[string]interface{}) peer.Response {
	bdoctor, _ := json.Marshal(member)
	bf, _ := json.Marshal(filterRec)
	return stub.MockInvoke(test_UUID, [][]byte{[]byte("AddCollaborator"), []byte(pID), bdoctor, bf})
}

func RemoveCollaborator(member asset.Doctor, pID string,
	filterRec map[string]interface{}) peer.Response {
	bdoctor, _ := json.Marshal(member)
	bf, _ := json.Marshal(filterRec)
	return stub.MockInvoke(test_UUID, [][]byte{[]byte("RemoveCollaborator"), []byte(pID), bdoctor, bf})
}

// -- Tests -- //
func TestInitLedger(t *testing.T) {
	if !assert.FileExists(t, "../../init.json", "Init file does not exist!") {
		t.FailNow()
	}
	result := stub.MockInit(test_UUID, nil)
	if PRINTRES {
		str, err := utils.IndentedJson(result, utils.INDENT_SPACE)
		assert.NoError(t, err, "")
		log.Println(str)
	}
	assert.EqualValuesf(t, shim.OK, result.Status, "Result status is not OK, get %d", result.Status)
	assert.NotNil(t, stub.Name, "Stub's name is nil!")
	assert.EqualValues(t, internal_name, stub.Name, "Stub's name is incorrect!")
}

func TestInitNewRecord(t *testing.T) {
	patientID := "p3"
	cnt := map[string]interface{}{
		"keystr":  "value1",
		"keybool": true,
	}
	nRecord := asset.Record{
		ID:      patientID,
		Type:    "test2",
		Time:    "2021-4-14 9:45:11",
		Content: cnt,
		Signature: asset.Doctor{
			Person: asset.Person{
				ID: "doct2", Name: "Banana", Age: 25,
			},
			Department: "Dep1",
		},
	}
	res := AddRecord(patientID, nRecord, nil)
	if PRINTRES {
		str, err := utils.IndentedJson(res, utils.INDENT_SPACE)
		assert.NoError(t, err, "")
		log.Println(str)
	}
	var records []asset.Record
	dec := json.NewDecoder(bytes.NewBuffer(res.Payload))
	dec.UseNumber()
	err := dec.Decode(&records)
	if assert.NoError(t, err, "No problem should appear unmarshalling") {
		if assert.Len(t, records, 1, "There should be 1 record of patient ", patientID) {
			if assert.Len(t, records[0].Collaborators, 1, "There should be 1 collaborator by default") {
				assert.EqualValues(t, records[0].Collaborators[0].Role, "manager", "The first collaborator should be manager")
			}
		}
	}
	// another one
	res = AddRecord(patientID, asset.Record{
		Type: "test3",
		Time: time.Now().Format("2006-1-2 15:04:05"),
		Signature: asset.Doctor{
			Person: asset.Person{
				ID: "doct3", Name: "Catt", Age: 26,
			},
			Department: "Dep2",
		},
	}, map[string]interface{}{
		"keybool": false,
		"keynum":  67,
	})
	if PRINTRES {
		str, err := utils.IndentedJson(res, utils.INDENT_SPACE)
		assert.NoError(t, err, "")
		log.Println(str)
	}
	dec = json.NewDecoder(bytes.NewBuffer(res.Payload))
	dec.UseNumber()
	err = dec.Decode(&records)
	if assert.NoError(t, err, "No problem should appear unmarshalling") {
		assert.Lenf(t, records, 2, "There should be 2 records of patient %s", patientID)
	}
}

func TestPatientInfoGet(t *testing.T) {
	patientID := "p1"
	res := GetPatientInfo(patientID)
	if PRINTRES {
		str, err := utils.IndentedJson(res, utils.INDENT_SPACE)
		assert.NoError(t, err, "")
		log.Println(str)
	}
	var pInfo asset.OutPatient
	err := json.Unmarshal(res.Payload, &pInfo)
	assert.NoError(t, err, "there should be no problem unmarshalling returning payload ")
}

func TestQueryDoctorByID(t *testing.T) {
	existingID := "doct1"
	res := QueryDoctorByID(existingID)
	if PRINTRES {
		str, err := utils.IndentedJson(res, utils.INDENT_SPACE)
		assert.NoError(t, err, "")
		log.Println(str)
	}
	assert.NotEmpty(t, res, existingID, "should be found at state map")
	res = QueryDoctorByID("notexist")
	if PRINTRES {
		str, err := utils.IndentedJson(res, utils.INDENT_SPACE)
		assert.NoError(t, err, "")
		log.Println(str)
	}
	assert.EqualValues(t, 500, res.Status, "invalid doctor should not exist")
}

func TestIsValidDoctor(t *testing.T) {
	doctor := asset.Doctor{
		Person: asset.Person{
			ID: "doct1", Name: "Apple", Age: 24,
		},
		Department: "Dep1",
	}
	notExistDoctor := asset.Doctor{
		Person: asset.Person{
			ID: "doct1", Name: "Apple", Age: 24,
		},
		Department: "Dep2",
	}
	res := IsValidDoctor(doctor)
	if PRINTRES {
		str, err := utils.IndentedJson(res, utils.INDENT_SPACE)
		assert.NoError(t, err, "")
		log.Println(str, string(res.Payload))
	}
	trueJSON, _ := json.Marshal(struct {
		Val bool `json:"val"`
	}{true})
	falseJSON, _ := json.Marshal(struct {
		Val bool `json:"val"`
	}{false})
	assert.JSONEq(t, string(trueJSON), string(res.Payload), "It should be a valid doctor")
	res = IsValidDoctor(notExistDoctor)
	if PRINTRES {
		str, err := utils.IndentedJson(res, utils.INDENT_SPACE)
		assert.NoError(t, err, "")
		log.Println(str)
	}
	assert.JSONEq(t, string(falseJSON), string(res.Payload), "Though dep is different, it should return false")
}

func TestGetMedicalRecord(t *testing.T) {
	patientID := "p1"
	resErr := GetRecord("")
	if PRINTRES {
		str, err := utils.IndentedJson(resErr, utils.INDENT_SPACE)
		assert.NoError(t, err, "")
		log.Println(str)
	}
	assert.NotEqualValues(t, "null", string(resErr.Payload),
		"Payload is 'null'")
	res := GetRecord(patientID)
	if PRINTRES {
		str, err := utils.IndentedJson(res, utils.INDENT_SPACE)
		assert.NoError(t, err, "")
		log.Println(str)
	}
	var rec []asset.Record
	err := json.Unmarshal(res.Payload, &rec)
	if assert.NoError(t, err, "Error is not nil! Error is ", err) {
		assert.Lenf(t, rec, 3, "There should be 3 records, found %d", len(rec))
	}
}

func TestPatientInfoSet(t *testing.T) {
	var pInfo asset.OutPatient
	patientID := "p1"
	res := SetPatientInfo(patientID, map[string]interface{}{
		"isMarried": false,
	})
	err := json.Unmarshal(res.Payload, &pInfo)
	if assert.NoError(t, err, "Nothing wrong happens to unmarshalling") {
		assert.False(t, pInfo.IsMarried, "This has successfully changed")
	}
}

func TestGetAllPatients(t *testing.T) {
	var p []asset.OutPatient
	res := GetAllPatients()
	err := json.Unmarshal(res.Payload, &p)
	if assert.NoError(t, err, "Nothing wrong happens to unmarshalling") {
		assert.Len(t, p, 3, "There are 3 patients overall")
	}
	np := asset.OutPatient{
		Person: asset.Person{
			ID:   "p23512345",
			Name: "testtest",
			Age:  75,
		},
		Birthday: time.Now().Format("2006-1-2"),
	}
	res = NewPatient(np)
	if PRINTRES {
		str, err := utils.IndentedJson(res, utils.INDENT_SPACE)
		assert.NoError(t, err, "")
		log.Println(str)
	}
	err = json.Unmarshal(res.Payload, &p)
	if assert.NoError(t, err, "Nothing wrong happens to unmarshalling") {
		assert.Len(t, p, 4, "There are 4 patients overall")
	}
}

func TestGetMRBydate(t *testing.T) {
	pid := "p1"
	start := time.Date(2021, 4, 8, 0, 0, 0, 0, time.Local)
	end := time.Now()
	res := QueryMRByStartToEndDate(start, end, pid)
	var rec []asset.Record
	err := json.Unmarshal(res.Payload, &rec)
	if assert.NoError(t, err, "Error is not nil! Error is", err) {
		assert.Len(t, rec, 2, "There should be 2 records")
	}
	start2 := time.Date(2021, 4, 10, 11, 45, 14, 0, time.Local)
	res2 := QueryMRByStartToEndDate(start2, end, "p2")
	err = json.Unmarshal(res2.Payload, &rec)
	if assert.NoError(t, err, "Error is not nil! Error is", err) {
		assert.Len(t, rec, 2, "There should be 2 records")
	}
}

// fixme: failed test
func TestCollaborator(t *testing.T) {
	res1 := GetRecord("p1")
	var precord []asset.Record
	err := json.Unmarshal(res1.Payload, &precord)
	if assert.NoError(t, err, "Nothing wrong happens to unmarshalling") {
		if assert.GreaterOrEqualf(t, len(precord), 1, "There should be at least one record") {
			assert.NotEmpty(t, precord[0].Collaborators, "There should be at least one collaborators")
		} else {
			t.FailNow()
		}
	}
	res2 := AddCollaborator(asset.Doctor{
		Person: asset.Person{
			Name: "Catt",
			ID:   "doct3",
			Age:  26,
		},
		Department: "Dep2",
	}, "p1", map[string]interface{}{
		"type": "Type1",
	})
	jsonsuccess, _ := json.Marshal(map[string]interface{}{
		"success": true,
		"error":   "null",
	})
	assert.JSONEq(t, string(jsonsuccess), string(res2.Payload), "Operation should be success")

	std, _ := time.ParseInLocation("2006-1-2 15:04:05", "2021-4-7 13:00:00", time.Local)
	end, _ := time.ParseInLocation("2006-1-2 15:04:05", "2021-4-8 08:00:00", time.Local)
	res3 := QueryMRByStartToEndDate(std, end, "p1")
	err = json.Unmarshal(res3.Payload, &precord)
	if assert.NoError(t, err, "Nothing wrong happens to unmarshalling") {
		log.Println("res3:", string(res3.Payload))
		if !assert.Len(t, precord, 1, "There should be 1 result") {
			t.FailNow()
		} else {
			assert.Len(t, precord[0].Collaborators, 3, "There should be 3 collaborators")
		}
	}
	res5 := RemoveCollaborator(asset.Doctor{
		Person: asset.Person{
			Name: "Apple",
			ID:   "doct1",
			Age:  24,
		},
		Department: "Dep1",
	}, "p1", map[string]interface{}{
		"type": "Type1",
	})
	assert.JSONEq(t, string(jsonsuccess), string(res5.Payload), "Operating should be success")
	res6 := RemoveCollaborator(asset.Doctor{
		Person: asset.Person{
			Name: "Banana",
			ID:   "doct2",
			Age:  25,
		},
		Department: "Dep1",
	}, "p1", map[string]interface{}{
		"type": "Type1",
	})
	m, err := utils.JsonToMap(string(res6.Payload))
	if assert.NoError(t, err, "No error transferring to map") {
		assert.False(t, m["success"].(bool), "Operation should fail")
		assert.Contains(t, m["error"], "cannot remove manager", "")
	}
}

func TestGetAllDoctors(t *testing.T) {
	var doc []asset.Doctor
	res := GetAllDoctors()
	if PRINTRES {
		str, err := utils.IndentedJson(res, utils.INDENT_SPACE)
		assert.NoError(t, err, "")
		log.Println(str)
	}
	err := json.Unmarshal(res.Payload, &doc)
	if assert.NoError(t, err, "Error is not nil! Error is", err) {
		assert.Len(t, doc, 3, "There should be 3 doctors overall")
	}
}
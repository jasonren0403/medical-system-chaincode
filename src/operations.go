package smartMedicineSystem

import (
	"bytes"
	"ccode/src/asset"
	"ccode/src/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/goinggo/mapstructure"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-protos-go/peer"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"
	"time"
)

var Name = "smartMedicineSystem"

type MedicalSystem struct {
	contractapi.Contract
}

// Init returns {"status": 200}
func (s *MedicalSystem) Init(stub shim.ChaincodeStubInterface) peer.Response {
	log.Println("Init() called")
	str, _ := os.Getwd()
	log.Println("pwd:", str)
	var (
		content []byte
		err     error
	)
	if strings.Contains(str, "test") {
		content, err = ioutil.ReadFile("../../init.json")
	} else {
		content, err = ioutil.ReadFile("init.json")
	}
	if err != nil {
		return shim.Error(err.Error())
	}
	var m map[string][]interface{}
	err = json.Unmarshal(content, &m)
	if err != nil {
		return shim.Error(err.Error())
	}
	doctors, res := m["doctors"]
	if res {
		for _, v := range doctors {
			var doctor asset.Doctor
			err := mapstructure.Decode(v, &doctor)
			if err != nil {
				log.Println("error on decoding doctor info,", err)
				return shim.Error(err.Error())
			}
			log.Println("Decoded doctor info:", doctor)
			dbyte, err := json.Marshal(doctor)
			if err != nil {
				return shim.Error(err.Error())
			}
			nkey, err := stub.CreateCompositeKey(utils.DOCTOR_STATE_KEY_PREFIX, []string{doctor.ID})
			log.Println("[Init] put doctor info state key", nkey)
			err = stub.PutState(nkey, dbyte)
			if err != nil {
				return shim.Error(err.Error())
			}
		}
	}
	log.Println("Doctor info init success")
	patients, res := m["patients"]
	if res {
		for _, v := range patients {
			var rec asset.MedicalRecord
			_byte, _ := json.Marshal(v)
			err = json.Unmarshal(_byte, &rec)
			if err != nil {
				log.Println("error on decoding patient json,", err)
				return shim.Error(err.Error())
			}
			log.Println("Decoded patient info:", rec)
			opbyte, err := json.Marshal(rec.PatientInfo)
			prbyte, err := json.Marshal(rec.PatientRecord)
			if err != nil {
				return shim.Error(err.Error())
			}
			nkey, err := stub.CreateCompositeKey(utils.PATIENT_INFO_STATE_KEY_PREFIX, []string{rec.PatientInfo.ID})
			log.Printf("[Init] Put patient info state key %s\n", nkey)
			err = stub.PutState(nkey, opbyte)
			for _, v := range rec.PatientRecord {
				nkey, err = stub.CreateCompositeKey(utils.PATIENT_RECORD_STATE_KEY_PREFIX, []string{rec.PatientInfo.ID, v.Time})
				log.Printf("[Init] Put record state key %s\n", nkey)
				err = stub.PutState(nkey, prbyte)
			}
			if err != nil {
				return shim.Error(err.Error())
			}
		}
	}
	log.Println("Patient info and records init success")
	return shim.Success(nil)
}

func (s *MedicalSystem) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fun, params := stub.GetFunctionAndParameters()
	log.Printf("Invoke() called >> \n\t%-10s%s\n\t%-10s%s", "Params:", params, "Function:", fun)
	var payload []byte
	switch fun {
	case "NewPatient":
		if len(params) < 1 {
			return shim.Error("Not enough param for new patient info init!")
		}
		info := params[0]
		var p asset.OutPatient
		_ = json.Unmarshal([]byte(info), &p)
		err := s.NewPatient(stub, p)
		if err != nil {
			return shim.Error("error creating new patient")
		}
		ps, err := s.GetAllPatients(stub)
		if err != nil {
			return shim.Error("error fetching all patient list")
		}
		payload, _ = json.Marshal(ps)
	case "InitNewRecord":
		if len(params) < 5 {
			return shim.Error("Not enough param for new patient record init!")
		}
		pID := params[0]
		ty := params[1]
		tm := params[2]
		cnt, _ := utils.JsonToMap(params[3])
		var doc asset.Doctor
		_ = json.Unmarshal([]byte(params[4]), &doc)
		err := s.InitNewRecord(stub, pID, ty, tm, cnt, doc)
		if err != nil {
			return shim.Error("Error in init new record for patientID " + pID + ":" + err.Error())
		}
		var rs []asset.Record
		rs, err = s.GetMedicalRecord(stub, pID)
		if err != nil {
			return shim.Error("Error fetching new records for patientID " + pID)
		}
		payload, _ = json.Marshal(rs)
	case "SetPatientInfo":
		if len(params) < 2 {
			return shim.Error("Not enough param for patient info setting!")
		}
		pid := params[0]
		m, err := utils.JsonToMap(params[1])
		if err != nil {
			log.Println(params[1])
			return shim.Error("Error params -> map!")
		}
		err = s.SetPatientInfo(stub, pid, m)
		if err != nil {
			log.Println(err)
			return shim.Error("Set patient info error!")
		}
		ninfo, err := s.GetPatientInfoByPID(stub, pid)
		if err != nil {
			log.Println(err)
			return shim.Error("Reget new patient info error!")
		}
		payload, _ = json.Marshal(ninfo)
	case "GetMedicalRecord":
		var pmrs []asset.Record
		var err error
		if len(params) == 0 {
			pmrs, err = s.GetMedicalRecord(stub, "")
		} else {
			pmrs, err = s.GetMedicalRecord(stub, params[0])
		}
		if err != nil {
			return shim.Error(err.Error())
		}
		payload, _ = json.Marshal(pmrs)
	case "IsValidDoctor":
		doc := params[0]
		var d asset.Doctor
		err := json.Unmarshal([]byte(doc), &d)
		log.Printf("Test if %s is a valid doctor", d.String())
		if err != nil {
			return shim.Error(err.Error())
		}
		res := s.IsValidDoctor(stub, d)
		payload, _ = json.Marshal(struct {
			Val bool `json:"val"`
		}{res})
	case "QueryDoctorByID":
		dID := params[0]
		res, err := s.QueryDoctorByID(stub, dID)
		if err != nil {
			return shim.Error(err.Error())
		}
		payload, _ = json.Marshal(res)
	case "GetPatientInfoByPID":
		if len(params) == 0 {
			return shim.Error("Support a pid(string) for this call!")
		}
		log.Println("Get patient info by pid called", "(pid =", params[0], ")")
		pInfo, err := s.GetPatientInfoByPID(stub, params[0])
		if err != nil {
			return shim.Error(err.Error())
		}
		payload, _ = json.Marshal(pInfo)
	case "GetAllPatients":
		r, err := s.GetAllPatients(stub)
		if err != nil {
			return shim.Error(err.Error())
		}
		payload, _ = json.Marshal(r)
	case "GetAllDoctors":
		r, err := s.GetAllDoctors(stub)
		if err != nil {
			return shim.Error(err.Error())
		}
		payload, _ = json.Marshal(r)
	case "GetMRByDate":
		if len(params) < 3 {
			return shim.Error("Not enough params to get mr by date")
		}
		pID := params[0]
		st := params[1]
		ed := params[2]
		r, err := s.GetMRByDate(stub, pID, st, ed)
		if err != nil {
			return shim.Error(err.Error())
		}
		payload, _ = json.Marshal(r)
	case "AddCollaborator":
		if len(params) < 3 {
			return shim.Error("Not enough params")
		}
		pID := params[0]
		doc := params[1]
		var d asset.Doctor
		_ = json.Unmarshal([]byte(doc), &d)
		m, _ := utils.JsonToMap(params[2])
		err := s.AddCollaborator(stub, d, pID, m)
		var errstring string = "null"
		if err != nil {
			errstring = err.Error()
		}
		payload, _ = json.Marshal(map[string]interface{}{
			"success": err == nil,
			"error":   errstring,
		})
	case "RemoveCollaborator":
		if len(params) < 3 {
			return shim.Error("Not enough params")
		}
		pID := params[0]
		doc := params[1]
		var d asset.Doctor
		_ = json.Unmarshal([]byte(doc), &d)
		m, _ := utils.JsonToMap(params[2])
		err := s.RemoveCollaborator(stub, d, pID, m)
		var errstring string = "null"
		if err != nil {
			errstring = err.Error()
		}
		payload, _ = json.Marshal(map[string]interface{}{
			"success": err == nil,
			"error":   errstring,
		})
	default:
		log.Println("Unknown function ", fun, "called")
		return shim.Error("Nothing has called")
	}
	return shim.Success(payload)
}

// NewPatient /* put a new patient to the world state */
func (s *MedicalSystem) NewPatient(stub shim.ChaincodeStubInterface, patient asset.OutPatient) error {
	p2, err := s.GetPatientInfoByPID(stub, patient.ID)
	if p2 != nil {
		if err != nil {
			return fmt.Errorf("patient already exists <origin error:%s>", err.Error())
		}
	}
	p, _ := json.Marshal(patient)
	dkey, _ := stub.CreateCompositeKey(utils.PATIENT_INFO_STATE_KEY_PREFIX, []string{patient.ID})
	log.Println("new patient key:", dkey)
	return stub.PutState(dkey, p)
}

// IsValidDoctor /* check if the doctor exists in the world state */
func (s *MedicalSystem) IsValidDoctor(stub shim.ChaincodeStubInterface, doctor asset.Doctor) bool {
	log.Println("checking if", doctor.ID, "is valid in world state")
	nkey, _ := stub.CreateCompositeKey(utils.DOCTOR_STATE_KEY_PREFIX, []string{doctor.ID})
	dbyte, err := stub.GetState(nkey)
	if err != nil {
		log.Println(err)
		return false
	}
	var _doctor asset.Doctor
	_ = json.Unmarshal(dbyte, &_doctor)
	return dbyte != nil && doctor == _doctor
}

// GetAllDoctors /* Return all the doctors in this system */
func (s *MedicalSystem) GetAllDoctors(stub shim.ChaincodeStubInterface) ([]asset.Doctor, error) {
	query, _ := stub.GetStateByPartialCompositeKey(utils.DOCTOR_STATE_KEY_PREFIX, []string{})
	var ps []asset.Doctor
	for query.HasNext() {
		t, _ := query.Next()
		var r1 asset.Doctor
		dec := json.NewDecoder(bytes.NewBuffer(t.GetValue()))
		dec.UseNumber()
		_ = dec.Decode(&r1)
		ps = append(ps, r1)
	}
	_ = query.Close()
	return ps, nil
}

// QueryDoctorByID /* find doctor info by doctor ID */
func (s *MedicalSystem) QueryDoctorByID(stub shim.ChaincodeStubInterface,
	dID string) (*asset.Doctor, error) {
	dKey, _ := stub.CreateCompositeKey(utils.DOCTOR_STATE_KEY_PREFIX, []string{dID})
	existing, err := stub.GetState(dKey)
	if err != nil {
		return nil, errors.New("Unable to interact with world state")
	}
	if existing == nil {
		return nil, fmt.Errorf("Current doctor dID=%s does not exist", dID)
	}
	var doc asset.Doctor
	err = json.Unmarshal(existing, &doc)
	return &doc, err
}

// InitNewRecord /* append a new record to the patient's records */
func (s *MedicalSystem) InitNewRecord(stub shim.ChaincodeStubInterface, patientID string,
	_type string, time string, content interface{}, signature asset.Doctor) error {
	records, err := s.GetMedicalRecord(stub, patientID)
	if err != nil {
		return err
	}
	if !s.IsValidDoctor(stub, signature) {
		return errors.New("not a valid doctor in database")
	}
	newRec := asset.Record{
		ID:        patientID,
		Type:      _type,
		Time:      time,
		Content:   content,
		Signature: signature,
	}
	if len(newRec.ID) == 0 {
		newRec.ID = patientID
	}
	if len(newRec.RecordID) == 0 {
		newRec.RecordID = utils.NewGenerator(32, "string").RandStr()
	}
	// the first doctor with one patient defaults to manager role
	newRec.Collaborators = append(newRec.Collaborators, asset.Collaborator{
		Doc:  signature,
		Role: "manager",
	})
	records = append(records, newRec)
	rec, _ := json.Marshal(records)
	dkey, _ := stub.CreateCompositeKey(utils.PATIENT_RECORD_STATE_KEY_PREFIX, []string{patientID, time})
	log.Println("new record dkey=", dkey)
	return stub.PutState(dkey, rec)
}

// RemoveCollaborator /* remove a doctor from patient's record's collaborator list */
func (s *MedicalSystem) RemoveCollaborator(stub shim.ChaincodeStubInterface, member asset.Doctor, pID string,
	filterRec map[string]interface{}) error {
	log.Println("Removing", member, "out of", pID, "'s record's collaborator list")
	res := s.IsValidDoctor(stub, member)
	if !res {
		return errors.New("not a valid doctor")
	}
	mrs, err := s.GetMedicalRecord(stub, pID)
	if err != nil {
		return err
	}
	var restoreTime string
Loop:
	for indexOuter, value := range mrs {
		for _, v := range filterRec {
			valof := reflect.ValueOf(value)
			for i := 0; i < valof.NumField(); i++ {
				if valof.Field(i).Interface() == v {
					cols := value.Collaborators
					for index, doc := range cols {
						if doc.Role == "manager" {
							return errors.New("cannot remove manager")
						}
						if member == doc.Doc {
							mrs[indexOuter].Collaborators = append(cols[:index], cols[index+1:]...)
							restoreTime = value.Time
							break Loop
						}
					}
				}
			}
		}
	}
	if len(restoreTime) == 0 {
		return errors.New("no valid record found")
	}
	dkey, _ := stub.CreateCompositeKey(utils.PATIENT_RECORD_STATE_KEY_PREFIX,
		[]string{pID, restoreTime})
	log.Println("after delete:", mrs)
	mrbytes, _ := json.Marshal(mrs)
	return stub.PutState(dkey, mrbytes)
}

// AddCollaborator /* add a new doctor to patient's record's collaborator list */
func (s *MedicalSystem) AddCollaborator(stub shim.ChaincodeStubInterface, member asset.Doctor, pID string,
	filterRec map[string]interface{}) error {
	log.Println("Adding", member, "to", pID, "'s record's collaborator list")
	res := s.IsValidDoctor(stub, member)
	if !res {
		return errors.New("not a valid doctor")
	}
	mrs, err := s.GetMedicalRecord(stub, pID)
	if err != nil {
		return err
	}
	var restoreTime string
Loop:
	for index, value := range mrs {
		for _, v := range filterRec {
			valof := reflect.ValueOf(value)
			for i := 0; i < valof.NumField(); i++ {
				if valof.Field(i).Interface() == v {
					for _, doc := range value.Collaborators {
						if member == doc.Doc {
							return errors.New("the doctor is already in this list")
						}
						mrs[index].Collaborators = append(value.Collaborators,
							asset.Collaborator{Doc: member, Role: "member"})
						restoreTime = value.Time
						break Loop
					}
				}
			}
		}
	}
	if len(restoreTime) == 0 {
		return errors.New("no valid record found")
	}
	dkey, _ := stub.CreateCompositeKey(utils.PATIENT_RECORD_STATE_KEY_PREFIX,
		[]string{pID, restoreTime})
	mrbytes, _ := json.Marshal(mrs)
	log.Println("after add:", string(mrbytes))
	return stub.PutState(dkey, mrbytes)
}

// SetPatientInfo /* Set the patient's info using key and values */
/* We define that the patient's ID and name cannot be changed */
func (s *MedicalSystem) SetPatientInfo(stub shim.ChaincodeStubInterface, ID string,
	kvs map[string]interface{}) error {
	pinfo, err := s.GetPatientInfoByPID(stub, ID)
	// reject if the patient's info is nil or error occurs
	if err != nil {
		return err
	}
	if pinfo == nil {
		return errors.New("PInfo is nil")
	}
	for k, v := range kvs {
		switch k {
		case "country":
			pinfo.Country = v.(string)
		case "region":
			pinfo.Region = v.(string)
		case "birthday":
			pinfo.Birthday = v.(string)
		case "isMarried":
			pinfo.IsMarried = v.(bool)
		case "career":
			pinfo.Career = v.(string)
		case "address":
			pinfo.Address = v.(string)
		case "age":
			pinfo.Age = v.(int)
		case "id", "name":
			return errors.New("cannot change ID or name")
		default:
			return errors.New("bad field to change")
		}
	}

	rec, _ := json.Marshal(pinfo)
	dkey, _ := stub.CreateCompositeKey(utils.PATIENT_INFO_STATE_KEY_PREFIX, []string{ID})
	log.Println("[SetPatientInfo]Composite key:", dkey)
	return stub.PutState(dkey, rec)
}

// GetPatientInfoByPID /* Get the patient's info by patient's ID */
func (s *MedicalSystem) GetPatientInfoByPID(stub shim.ChaincodeStubInterface,
	patientID string) (*asset.OutPatient, error) {
	dkey, _ := stub.CreateCompositeKey(utils.PATIENT_INFO_STATE_KEY_PREFIX, []string{patientID})
	existing, err := stub.GetState(dkey)
	if err != nil {
		return nil, errors.New("Unable to interact with world state")
	}
	if existing == nil {
		return nil, fmt.Errorf("current patient <PID=%s> does not exist", patientID)
	}
	var patient asset.OutPatient
	err = json.Unmarshal(existing, &patient)
	return &patient, err
}

// GetAllPatients return all valid patients in this system
func (s *MedicalSystem) GetAllPatients(stub shim.ChaincodeStubInterface) ([]asset.OutPatient, error) {
	query, _ := stub.GetStateByPartialCompositeKey(utils.PATIENT_INFO_STATE_KEY_PREFIX, []string{})
	var ps []asset.OutPatient
	for query.HasNext() {
		t, _ := query.Next()
		var r1 asset.OutPatient
		dec := json.NewDecoder(bytes.NewBuffer(t.GetValue()))
		dec.UseNumber()
		_ = dec.Decode(&r1)
		ps = append(ps, r1)
	}
	_ = query.Close()
	return ps, nil
}

// GetMedicalRecord /* Get the patient's record(s) by patient's ID */
func (s *MedicalSystem) GetMedicalRecord(stub shim.ChaincodeStubInterface,
	patientID string) ([]asset.Record, error) {
	var queryRes shim.StateQueryIteratorInterface
	var err error
	if len(patientID) == 0 {
		queryRes, err = stub.GetStateByPartialCompositeKey(utils.PATIENT_RECORD_STATE_KEY_PREFIX,
			[]string{})
	} else {
		queryRes, err = stub.GetStateByPartialCompositeKey(utils.PATIENT_RECORD_STATE_KEY_PREFIX,
			[]string{patientID})
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state %s", err.Error())
	}
	var r []asset.Record
	var lastvalue []byte
	for queryRes.HasNext() {
		t, _ := queryRes.Next()
		if len(patientID) == 0 {
			log.Println("queryAll", string(t.GetValue()))
		} else {
			log.Println("queryByID:", patientID, string(t.GetValue()))
		}
		if lastvalue == nil {
			goto newRecord
		}
		if bytes.Equal(lastvalue, t.GetValue()) {
			log.Println("Equal value, skipped")
			continue
		}
	newRecord:
		var rs []asset.Record
		dec := json.NewDecoder(bytes.NewBuffer(t.GetValue()))
		dec.UseNumber()
		_ = dec.Decode(&rs)
		if len(r) == 0 {
			lastvalue = t.GetValue()
			r = append(r, rs...)
			log.Println("first append")
			continue
		}
		for i := 0; i < len(rs); i++ {
			for j := 0; j < len(r); j++ {
				if !asset.RecordEquals(rs[i], r[j]) {
					// rs[i] not in r --> append
					for _, v := range rs[i].Collaborators {

						if !asset.InCollaboratorList(r[j].Collaborators, v) {
							r = append(r, rs[i])
						}
					}
				}
			}
		}

		lastvalue = t.GetValue()
	}
	_ = queryRes.Close()
	return r, nil
}

// GetMRByDate /* Get the patient's record(s) by patient's ID and date */
func (s *MedicalSystem) GetMRByDate(stub shim.ChaincodeStubInterface, patientID string,
	startDate string, endDate string) ([]asset.Record, error) {
	log.Println("Find records of", patientID, "that is from", startDate, "to", endDate)
	all, err := s.GetMedicalRecord(stub, patientID)
	//log.Println(all)
	var record []asset.Record
	if err != nil {
		return record, err
	}
	startTime, _ := time.ParseInLocation("2006-1-2 15:04:05", startDate, time.Local)
	endTime, _ := time.ParseInLocation("2006-1-2 15:04:05", endDate, time.Local)
	for _, v := range all {
		thisTime, _ := time.ParseInLocation("2006-1-2 15:04:05", v.Time, time.Local)
		//log.Println("v.Time=",v.Time)
		//log.Println("startTime=",startTime,"thisTime=",thisTime,"endTime=",endTime)
		if thisTime.After(startTime) && thisTime.Before(endTime) {
			record = append(record, v)
		}
	}
	return record, nil
}

func (s *MedicalSystem) GetName() string {
	return Name
}

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
	"strings"
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
			err = stub.PutState(nkey, opbyte)
			nkey, err = stub.CreateCompositeKey(utils.PATIENT_RECORD_STATE_KEY_PREFIX, []string{rec.PatientInfo.ID})
			err = stub.PutState(nkey, prbyte)
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
		log.Println("Get patient info by pid called", "( pid =", params[0], ")")
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
	default:
		log.Println("Unknown function ", fun, "called")
		return shim.Error("Nothing has called")
	}
	return shim.Success(payload)
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
	records = append(records, newRec)
	rec, _ := json.Marshal(records)
	dkey, _ := stub.CreateCompositeKey(utils.PATIENT_RECORD_STATE_KEY_PREFIX, []string{patientID})
	return stub.PutState(dkey, rec)
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
		return nil, fmt.Errorf("Current patient <PID=%s> does not exist", patientID)
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
	if len(patientID) == 0 {
		queryRes, err := stub.GetStateByPartialCompositeKey(utils.PATIENT_RECORD_STATE_KEY_PREFIX,
			[]string{})
		if err != nil {
			return nil, fmt.Errorf("failed to read from world state %s", err.Error())
		}
		var r []asset.Record
		for queryRes.HasNext() {
			t, _ := queryRes.Next()
			var r1 []asset.Record
			dec := json.NewDecoder(bytes.NewBuffer(t.GetValue()))
			dec.UseNumber()
			_ = dec.Decode(&r1)
			r = append(r, r1...)
		}
		_ = queryRes.Close()
		return r, nil
	}
	dkey, _ := stub.CreateCompositeKey(utils.PATIENT_RECORD_STATE_KEY_PREFIX, []string{patientID})
	mr, err := stub.GetState(dkey)
	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if mr == nil {
		return nil, fmt.Errorf("%s does not exist", patientID)
	}

	var record []asset.Record
	// keep the precision
	dec := json.NewDecoder(bytes.NewBuffer(mr))
	dec.UseNumber()
	_ = dec.Decode(&record)
	return record, nil
}

func (s *MedicalSystem) GetName() string {
	return Name
}

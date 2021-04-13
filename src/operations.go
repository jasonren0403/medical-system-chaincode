package smartMedicineSystem

import (
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
	"time"
)

var Name = "smartMedicineSystem"

type MedicalSystem struct {
	contractapi.Contract
}

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
				fmt.Println("error on decoding doctor info,", err)
				return shim.Error(err.Error())
			}
			log.Println("Decoded doctor info:", doctor)
			dbyte, err := json.Marshal(doctor)
			if err != nil {
				return shim.Error(err.Error())
			}
			err = stub.PutState("Doctor"+doctor.ID, dbyte)
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
				fmt.Println("error on decoding patient json,", err)
				return shim.Error(err.Error())
			}
			log.Println("Decoded patient info:", rec)
			opbyte, err := json.Marshal(rec.PatientInfo)
			prbyte, err := json.Marshal(rec.PatientRecord)
			if err != nil {
				return shim.Error(err.Error())
			}
			err = stub.PutState("Patient"+rec.PatientInfo.ID, opbyte)
			err = stub.PutState("PatientRecord"+rec.PatientInfo.ID, prbyte)
			if err != nil {
				return shim.Error(err.Error())
			}
		}
	}
	log.Println("Patient info and records init success")
	return shim.Success(nil)
}

func (s *MedicalSystem) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	log.Println("Invoke() called")
	fun, params := stub.GetFunctionAndParameters()
	log.Println("Params:", params)
	var payload []byte
	switch fun {
	case "InitNewRecord":
		log.Println("Init new record called")
		//payload,_ = json.Marshal(res)
	case "SetPatientInfo":
		log.Println("Set patient info called")
	case "GetMedicalRecord":
		log.Println("Get medical record called")
	case "isValidDoctor":
		doc := params[0]
		log.Printf("Test if %s is a valid doctor", doc)
		var d asset.Doctor
		err := json.Unmarshal([]byte(doc), &d)
		if err != nil {
			return shim.Error(err.Error())
		}
		res := s.IsValidDoctor(stub, d)
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
		log.Println(pInfo)
		payload, _ = json.Marshal(pInfo)
	default:
		log.Println("Unknown function ", fun, "called")
		return shim.Error("Nothing has called")
	}
	return shim.Success(payload)
}

/* check if the doctor exists in the world state */
func (s *MedicalSystem) IsValidDoctor(stub shim.ChaincodeStubInterface, doctor asset.Doctor) bool {
	dbyte, err := stub.GetState(utils.CreateDoctorKey(doctor.ID))
	if err != nil {
		return false
	}
	return dbyte != nil
}

/* append a new record to the patient's records */
func (s *MedicalSystem) InitNewRecord(stub shim.ChaincodeStubInterface, patientID string,
	_type string, time string, content interface{}, signature asset.Doctor) error {
	records, err := s.GetMedicalRecord(stub, patientID)
	if err != nil {
		return err
	}
	newRec := asset.Record{
		Type:      _type,
		Time:      time,
		Content:   content,
		Signature: signature,
	}
	records = append(records, newRec)
	rec, _ := json.Marshal(records)
	return stub.PutState(utils.CreatePatientRecordKey(patientID), rec)
}

/* Set the patient's info using key and values */
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
		switch strings.ToLower(k) {
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

	return stub.PutState("Patient"+ID, rec)
}

/* Get the patient's info by patient's ID */
func (s *MedicalSystem) GetPatientInfoByPID(stub shim.ChaincodeStubInterface,
	patientID string) (*asset.OutPatient, error) {
	existing, err := stub.GetState(utils.CreatePatientInfoKey(patientID))
	if err != nil {
		return nil, errors.New("Unable to interact with world state")
	}
	if existing == nil {
		return nil, fmt.Errorf("Current patient <PID=%s> does not exist", patientID)
	}
	var patient *asset.OutPatient
	err = json.Unmarshal(existing, patient)
	return patient, err
}

/* Get the patient's record(s) by patient's ID */
func (s *MedicalSystem) GetMedicalRecord(stub shim.ChaincodeStubInterface,
	patientID string) ([]asset.Record, error) {
	mr, err := stub.GetState(utils.CreatePatientRecordKey(patientID))
	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if mr == nil {
		return nil, fmt.Errorf("%s does not exist", patientID)
	}

	record := make([]asset.Record, 100)
	_ = json.Unmarshal(mr, &record)
	return record, nil
}

/* Get the patient's record(s) by date range, [startDate,endDate) */
func (s *MedicalSystem) GetMRbyDateRange(stub shim.ChaincodeStubInterface,
	patientID string, startDate time.Time, endDate time.Time) ([]asset.MedicalRecord, error) {
	return nil, nil
}

func (s *MedicalSystem) GetName() string {
	return Name
}

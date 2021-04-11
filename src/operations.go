package smartMedicineSystem

import (
	"ccode/src/asset"
	"encoding/json"
	"fmt"
	"github.com/goinggo/mapstructure"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-protos-go/peer"
	"io/ioutil"
	"log"
	"time"
)

var Name = "smartMedicineSystem"

type MedicalSystem struct {
	contractapi.Contract
}

func (s *MedicalSystem) Init(stub shim.ChaincodeStubInterface) peer.Response {
	content, err := ioutil.ReadFile("init.json")
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
			fmt.Println(doctor)
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
			fmt.Println(rec)
			pbyte, err := json.Marshal(rec)
			if err != nil {
				return shim.Error(err.Error())
			}
			err = stub.PutState("Patient"+rec.PatientInfo.ID, pbyte)
			if err != nil {
				return shim.Error(err.Error())
			}
		}
	}
	log.Println("Patient info and records init success")
	return shim.Success(nil)
}

func (s *MedicalSystem) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (s *MedicalSystem) InitNewRecord(ctx contractapi.TransactionContextInterface, record *asset.MedicalRecord) error {
	rec, _ := json.Marshal(record)
	return ctx.GetStub().PutState(record.PatientInfo.ID, rec)
}

func (s *MedicalSystem) SetPatientInfo(ctx contractapi.TransactionContextInterface, ID string,
	field string, nvalue interface{}) error {
	m := make(map[string]interface{})
	m["ID"] = ID
	pinfo, err := s.GetPatientInfo(ctx, m)
	if err != nil || pinfo == nil {
		return err
	}

	// todo: change the record according to 'field' param
	rec, _ := json.Marshal(pinfo)

	return ctx.GetStub().PutState(ID, rec)
}

func (s *MedicalSystem) GetPatientInfo(ctx contractapi.TransactionContextInterface,
	query_map map[string]interface{}) (*asset.MedicalRecord, error) {
	return nil, nil
}

func (s *MedicalSystem) GetMedicalRecord(ctx contractapi.TransactionContextInterface, ID string) (*asset.MedicalRecord, error) {
	mr, err := ctx.GetStub().GetState(ID)
	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if mr == nil {
		return nil, fmt.Errorf("%s does not exist", ID)
	}

	record := new(asset.MedicalRecord)
	_ = json.Unmarshal(mr, record)
	return record, nil
}

func (s *MedicalSystem) GetMRbyDateRange(ctx contractapi.TransactionContextInterface,
	ID string, startDate time.Time, endDate time.Time) ([]asset.MedicalRecord, error) {
	return nil, nil
}

func (s *MedicalSystem) Transfer(ctx contractapi.TransactionContextInterface) (bool, error) {
	// Transfer a patient to another hospital
	return false, nil
}

func (s *MedicalSystem) GetName() string {
	return Name
}

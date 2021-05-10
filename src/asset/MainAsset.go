package asset

import (
	"encoding/json"
	"strings"
	"time"
)

type MRTime time.Time

type MedicalRecord struct {
	PatientInfo   OutPatient `json:"patient_info"`
	PatientRecord []Record   `json:"patient_records"`
}

type Collaborator struct {
	Doc  Doctor `json:"doctor"`
	Role string `json:"role" validate:"required,oneof=manager member"`
}

type Record struct {
	Collaborators []Collaborator `json:"collaborators"`
	ID            string         `json:"id" validate:"required"`
	Type          string         `json:"type"`
	Time          string         `json:"time" validate:"required,datetime=2006-1-2 15:04:05"`
	Content       interface{}    `json:"content"`
	Signature     Doctor         `json:"signed_by"`
}

func (m MedicalRecord) SetPInfo(patient OutPatient) {
	m.PatientInfo = patient
}

func (m MedicalRecord) SetPRecord(rec []Record) {
	m.PatientRecord = rec
}

func (r *MRTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		return err
	}
	*r = MRTime(t)
	return nil
}

func (r *MRTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(r)
}

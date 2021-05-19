package asset

import (
	"encoding/json"
	"fmt"
	"reflect"
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
	RecordID      string         `json:"record_id" validate:"required"`
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

func (c *Collaborator) String() string {
	return fmt.Sprintf("<collaborator doctor=%s role=%s>", c.Doc.ID, c.Role)
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

func RecordEquals(r1, r2 Record) bool {
	return reflect.DeepEqual(r1.Collaborators, r2.Collaborators) && r1.ID == r2.ID &&
		r1.Signature == r2.Signature && r1.Type == r2.Type && r1.Time == r2.Time && r1.RecordID == r2.RecordID
}

/**
"collaborators":[
{"doctor":{"person":{"id":"doct1","name":"Apple","age":24},"department":"Dep1"},"role":"member"},
{"doctor":{"person":{"id":"doct2","name":"Banana","age":25},"department":"Dep1"},"role":"manager"},
{"doctor":{"person":{"id":"doct3","name":"Catt","age":26},"department":"Dep2"},"role":"member"}
]
*/

func InCollaboratorList(haystack []Collaborator, needle Collaborator) bool {
	for i := 0; i < len(haystack); i++ {
		if needle == haystack[i] {
			return true
		}
	}
	return false
}

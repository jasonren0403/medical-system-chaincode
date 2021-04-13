package utils

import "fmt"

// In this system, the patient's state used 'Patient'+ <patientID> as key,
// and the doctor's state used 'Doctor' + <doctorID> as key.
// Patient records' state: 'PatientRecord' + <patientID>
const (
	PATIENT_INFO_STATE_KEY_PREFIX   = "Patient"
	DOCTOR_STATE_KEY_PREFIX         = "Doctor"
	PATIENT_RECORD_STATE_KEY_PREFIX = "PatientRecord"
)

func CreatePatientInfoKey(patientID string) string {
	return fmt.Sprintf("%s_%s", PATIENT_INFO_STATE_KEY_PREFIX, patientID)
}

func CreatePatientRecordKey(patientID string) string {
	return fmt.Sprintf("%s_%s", PATIENT_RECORD_STATE_KEY_PREFIX, patientID)
}

func CreateDoctorKey(doctorID string) string {
	return fmt.Sprintf("%s_%s", DOCTOR_STATE_KEY_PREFIX, doctorID)
}

package utils

import "fmt"

// In this system, the patient's state used 'Patient'+ <patientID> as key,
// and the doctor's state used 'Doctor' + <doctorID> as key.
// Patient records' state: 'PatientRecord' + <patientID>
// 21/5/9: replace these functions with Composite Key API provided by shim.ChaincodeStubInterface
const (
	PATIENT_INFO_STATE_KEY_PREFIX   = "Patient~ID"
	DOCTOR_STATE_KEY_PREFIX         = "Doctor~ID"
	PATIENT_RECORD_STATE_KEY_PREFIX = "PatientRecord~PID"
)

// Deprecated: Use composite key API instead
func CreatePatientInfoKey(patientID string) string {
	return fmt.Sprintf("%s_%s", PATIENT_INFO_STATE_KEY_PREFIX, patientID)
}

// Deprecated: Use composite key API instead
func CreatePatientRecordKey(patientID string) string {
	return fmt.Sprintf("%s_%s", PATIENT_RECORD_STATE_KEY_PREFIX, patientID)
}

// Deprecated: Use composite key API instead
func CreateDoctorKey(doctorID string) string {
	return fmt.Sprintf("%s_%s", DOCTOR_STATE_KEY_PREFIX, doctorID)
}

package _excluded

import (
	"ccode/src/asset"
	"github.com/shopspring/decimal"
	"time"
)

const (
	Food = iota
	Greens
	Pets
	Insects
)

type SingleAllergyRecord struct{
	Type string `json:"allergyType"`
	Time time.Time `json:"allergyStartTime"`
}

type SingleHospitalizationRecord struct{
	Patient             asset.Person `json:"patient"`
	HospitalizationID   int          `json:"hospitalization_id"`
	BedID               string       `json:"bed_id"`
	Department          string       `json:"department"`
	HospitalizationDate time.Time    `json:"hospitalization_date"`
	History             struct{
		Past string `json:"past"`
		Personal string `json:"personal"`
		Marriage string `json:"marriage"`
		Family string `json:"family"`
	} `json:"history"`
	BodyTests []SingleTestSheet `json:"body_tests"`
}

type SingleTemperatureRecord struct {
	Temperature decimal.Decimal `json:"temperature"`
	Pulse int16 `json:"pulse"`
	Time time.Time `json:"time"`
	BloodPressure struct{
		Upper int16 `json:"upper"`
		Lower int16 `json:"lower"`
	} `json:"blood_pressure"`
	Weight decimal.Decimal      `json:"weight"`
	OtherRecords []SimpleRecord `json:"others"`
}

type SingleMedicalOrder struct{
	Starts struct{
		Time time.Time           `json:"time"`
		Order string             `json:"order_content"`
		Signature []asset.Doctor `json:"signed_by"`
	}
	Ends struct{
		Time time.Time           `json:"time"`
		Signature []asset.Doctor `json:"signed_by"`
	}
}

type SingleTestSheet struct{
	ID int8 `json:"id"`
	Item string `json:"item"`
	Value interface{} `json:"value"`
	Unit string `json:"unit"`
	ReferenceValue interface{} `json:"ref_value"` //“参考范围/值”
}

type SingleConsultationRecord struct{
	Application      `json:"application"`
	ConsultationView `json:"consultationView"`
}

type Application struct{
	Patient           asset.Person `json:"patient"`
	Reason            string       `json:"reason"`
	Goal              string       `json:"goal"`
	RequestDoctorFrom asset.Doctor `json:"request_doctor"`
}
type ConsultationView struct{
	Time           time.Time    `json:"time"`
	Suggestion     string       `json:"suggestion"`
	AssignedDoctor asset.Doctor `json:"assigned_doctor"`
}

type SimpleRecord struct{
	Type string `json:"type"`
	Content string `json:"content"`
}

type Advice struct{
	PersonFrom   asset.Person `json:"from"`
	Advice       string       `json:"advice"`
	Relationship string       `json:"relationship"`
	Time         time.Time    `json:"time"`
}

type SingleTranscript struct{
	Date          time.Time    `json:"discussion_date"`
	Location      string       `json:"location"`
	Patient       asset.Person `json:"patient"`
	DiscussResult string       `json:"discuss_result"`
	Discussion    string       `json:"discuss_contents"`
}

type SingleNormalCourseRecord struct{
	Condition string `json:"condition_changes"`
	ImplementationRecord struct{
		Process  string       `json:"process"`
		Result   string       `json:"result"`
		Reaction string       `json:"reaction"`
		Others   SimpleRecord `json:"others"`
	} `json:"implementation"`
	OperationRecord struct{
		Phenomenon string `json:"phenomenon"`
		Operations []string `json:"operations"`
		PatientReaction string `json:"patient_reaction"`
	} `json:"operation"`
	AdviceFromOthers []Advice `json:"advices"`
}

// define a struct for recording patient's info
type MedicalRecord struct {
	Patient asset.Person `json:"patient"`
	/**
	 * According to 《医疗事故处理条例》, a medical record should consist of
	 * subjective and objective medical materials
	 */
	//①死亡病例讨论记录;②疑难病例讨论记录;③上级医师查房记录;④会诊意见;⑤病程记录。
	Subjective struct{
		DeathMedicalRecords []SingleTranscript     `json:"deathMedicalRecords"`
		DifficultMedicalRecords []SingleTranscript `json:"difficultMedicalRecords"`
		RoundRecords []SimpleRecord                `json:"roundRecords"`
		Consultations []SingleConsultationRecord   `json:"consultations"`
		Course struct{
			FirstCourse struct{
				PatientInfo asset.Person `json:"patient_info"`
				PreChecks   struct{}     `json:"pre_checks"`
				Diagnosis   string       `json:"diagnosis"`
				Suggestion  string       `json:"suggestion"`
			} `json:"firstCourse"`
			NormalCourse []SingleNormalCourseRecord `json:"normalCourse"`
		} `json:"courseRecord"`
	} `json:"subjective"`
	//①门诊病历;②入院记录;③体温单;④医嘱单;⑤化验单(检验报告);
	//⑥医学影像检查资料;⑦特殊检查同意书、手术同意书;⑧手术及麻醉记录单;⑨病理资料;⑩护理记录。
	Objective struct{
		OutPatient2 struct{
			asset.OutPatient
			Allergies []SingleAllergyRecord `json:"allergyHistories"`
		} `json:"outPatient"`
		Hospitalization []SingleHospitalizationRecord `json:"hospitalizationHistories"`
		TemperatureRecord []SingleTemperatureRecord   `json:"temperatureRecords"`
		MedicalOrders []SingleMedicalOrder            `json:"medicalOrders"`
		TestSheet []SingleTestSheet                   `json:"testSheets"`
		ImageMaterials []SimpleRecord                 `json:"imageMaterials"`
		Agreement struct{
			Content string `json:"content"`
			Agreed bool `json:"isAgreed"`
		}
		Others []SimpleRecord `json:"others"`
	} `json:"objective"`
}

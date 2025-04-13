package domain

import "time"

type Task struct {
	ID              string          `json:"id"`
	CreatedAt       time.Time       `json:"created_at"`
	Processed       bool            `json:"processed"`
	NeedRetry       bool            `json:"need_retry"`
	Error           string          `json:"error"`
	Url             string          `json:"url"`
	Status          string          `json:"status"`
	DetectionResult DetectionResult `json:"detection_result,omitempty"`
}

type DetectionResult struct {
	Detections []Detection `json:"detections"`
}

type Detection struct {
	Class DetectionClass  `json:"class"`
	Score float64 `json:"score"`
	Box   []int   `json:"box"`
}

type DetectionClass string

const (
	DetectionClassFemaleGenitaliaCovered DetectionClass = "FEMALE_GENITALIA_COVERED"
	DetectionClassFaceFemale             DetectionClass = "FACE_FEMALE"
	DetectionClassButtocksExposed        DetectionClass = "BUTTOCKS_EXPOSED"
	DetectionClassFemaleBreastExposed    DetectionClass = "FEMALE_BREAST_EXPOSED"
	DetectionClassFemaleGenitaliaExposed DetectionClass = "FEMALE_GENITALIA_EXPOSED"
	DetectionClassMaleBreastExposed      DetectionClass = "MALE_BREAST_EXPOSED"
	DetectionClassAnusExposed            DetectionClass = "ANUS_EXPOSED"
	DetectionClassFeetExposed            DetectionClass = "FEET_EXPOSED"
	DetectionClassBellyCovered           DetectionClass = "BELLY_COVERED"
	DetectionClassFeetCovered            DetectionClass = "FEET_COVERED"
	DetectionClassArmpitsCovered         DetectionClass = "ARMPITS_COVERED"
	DetectionClassArmPitsExposed         DetectionClass = "ARMPITS_EXPOSED"
	DetectionClassFaceMale               DetectionClass = "FACE_MALE"
	DetectionClassBellyExposed           DetectionClass = "BELLY_EXPOSED"
	DetectionClassMaleGenitaliaExposed   DetectionClass = "MALE_GENITALIA_EXPOSED"
	DetectionClassAnusCovered            DetectionClass = "ANUS_COVERED"
	DetectionClassFemaleBreastCovered    DetectionClass = "FEMALE_BREAST_COVERED"
	DetectionClassButtocksCovered        DetectionClass = "BUTTOCKS_COVERED"
)

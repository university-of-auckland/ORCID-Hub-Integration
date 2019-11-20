package main

import (
	"errors"
	"strconv"
	"strings"
)

// Degree API student-v1 degree response message.
type Degree struct {
	ID               string `json:"id"`
	StudentDegNbr    string `json:"studentDegNbr"`
	Code             string `json:"degreeCode"`
	Desc             string `json:"degreeDesc"`
	AcadCareer       string `json:"degAcadCareer"`
	ConferDate       string `json:"degreeConferDate"`
	HonorsPrefix     string `json:"honorsPrefix"`
	HonorsSuffix     string `json:"honorsSuffix"`
	AcadDegreeStatus string `json:"degAcadDegreeStatus"`
	ProspectusCode   string `json:"prospectusCode"`
	Plans            []struct {
		AcadPlanCode        string `json:"acadPlanCode"`
		AcadPlanDesc        string `json:"acadPlanDesc"`
		DgpAcadCareer       string `json:"dgpAcadCareer"`
		StudentCareerNbr    int    `json:"studentCareerNbr"`
		DgpAcadDegreeStatus string `json:"dgpAcadDegreeStatus"`
		DegreeStatusDate    string `json:"degreeStatusDate"`
		AcadProgCode        string `json:"acadProgCode"`
		AcadProgGroupCode   int    `json:"acadProgGroupCode"`
		AcadProgGroup       string `json:"acadProgGroup"`
		AcadProgLevelCode   string `json:"acadProgLevelCode"`
		AcadProgLevel       string `json:"acadProgLevel"`
		AcadOrgCode         string `json:"acadOrgCode"`
		AcadGroupDesc       string `json:"acadGroupDesc"`
	} `json:"degreePlans"`
}

// Degrees - array of degrees
type Degrees []Degree

// Qualification - external organisations-qualification entry
type Qualification struct {
	Type        string `json:"type"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Country     string `json:"country"`
}

// Qualifications - array of qualifications
type Qualifications []Qualification

// propagateToHub adds degree/education records to the current affiliation task.
func (degrees Degrees) propagateToHub(email, orcid string) (count int, err error) {

	count = len(degrees)
	if count == 0 {
		return 0, errors.New("no degree entry")
	}

	records := make([]Record, count)
	for i, d := range degrees {

		degreeName, ok := qualifications[d.Code]
		if !ok {
			degreeName, ok = degreeCodes[strings.ToUpper(d.Desc)]
			if !ok {
				degreeName = d.Desc
			}
		}
		date := strings.Split(d.ConferDate, "T")[0]
		records[i] = Record{
			AffiliationType: "education",
			EndDate:         date,
			StartDate:       date,
			LocalID:         d.ID + "/" + d.StudentDegNbr,
			Email:           email,
			Orcid:           orcid,
			Role:            degreeName,
			IsActive:        true,
		}
	}
	// Make sure the task set-up is comlete

	var task Task
	err = oh.patch("api/v1/affiliations/"+strconv.Itoa(taskID), Task{ID: taskID, Records: records}, &task)
	if err != nil {
		log.Error("failed to update the taks: ", err)
		return
	}
	taskRecordCountMutex.Lock()
	taskRecordCount += count
	taskRecordCountMutex.Unlock()

	return
}

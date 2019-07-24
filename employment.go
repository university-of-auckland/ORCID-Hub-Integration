package main

type Employment struct {
	AcademicStaffFTE int    `json:"academicStaffFTE"`
	EmployeeID       string `json:"employeeID"`
	Job              []struct {
		Company                              string `json:"company"`
		CostCentre                           string `json:"costCentre"`
		DepartmentDescription                string `json:"departmentDescription"`
		DepartmentID                         string `json:"departmentID"`
		EffectiveDate                        string `json:"effectiveDate"`
		EffectiveSequence                    int    `json:"effectiveSequence"`
		EmployeeRecord                       int    `json:"employeeRecord"`
		EmployeeStatus                       string `json:"employeeStatus"`
		EmployeeType                         string `json:"employeeType"`
		FullTimeEquivalent                   int    `json:"fullTimeEquivalent"`
		HrStatus                             string `json:"hrStatus"`
		JobCode                              string `json:"jobCode"`
		JobCodeDescription                   string `json:"jobCodeDescription"`
		JobEndDate                           string `json:"jobEndDate"`
		JobGrade                             string `json:"jobGrade"`
		JobIndicator                         string `json:"jobIndicator"`
		JobStartDate                         string `json:"jobStartDate"`
		LastHRaction                         string `json:"lastHRaction"`
		Location                             string `json:"location"`
		LocationDescription                  string `json:"locationDescription"`
		OrganizationalRelation               string `json:"organizationalRelation"`
		ParentDepartmentDescription          string `json:"parentDepartmentDescription"`
		PoiType                              string `json:"poiType"`
		PositionDescription                  string `json:"positionDescription"`
		PositionNumber                       string `json:"positionNumber"`
		PrimaryActivityCentreDeptDescription string `json:"primaryActivityCentreDeptDescription"`
		PrimaryActivityCentreDeptID          string `json:"primaryActivityCentreDeptID"`
		ReportsToPosition                    string `json:"reportsToPosition"`
		SalAdminPlan                         string `json:"salAdminPlan"`
		StandardHours                        int    `json:"standardHours"`
		SupervisorID                         string `json:"supervisorID"`
		UpdatedDateTime                      string `json:"updatedDateTime"`
	} `json:"job"`
	ProfessionalStaffFTE int    `json:"professionalStaffFTE"`
	RequestTimeStamp     string `json:"requestTimeStamp"`
	UniServicesFTE       int    `json:"uniServicesFTE"`
}

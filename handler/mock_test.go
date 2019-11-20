//+build test

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"unicode"
)

// isValidID validates employment/student ID
func isValidID(uid string) bool {
	if l := len(uid); l < 8 || l > 10 {
		return false
	}
	for _, r := range uid {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func getenv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func createMockHandler(t *testing.T) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		url, ru := r.URL, r.URL.RequestURI()
		switch {
		case ru == "/ping":
			w.WriteHeader(http.StatusNoContent)
		case ru == "/oauth/token":
			io.WriteString(w, `{"access_token": "7jsxDZceygy2xNbK2M23sD5eyHimtx", "expires_in": 86400, "token_type": "Bearer", "scope": ""}`)
		case ru == "/api/v1/tasks?type=AFFILIATION&status=INACTIVE":
			if withTasks {
				io.WriteString(w, `[
	{"created-at":"2019-07-24T08:47:09","filename":"UOA-OH-INTEGRATION-TASK-pv51ql.json","id":781,"records":[],"status":"ACTIVE","task-type":"AFFILIATION","updated-at":"2019-07-24T09:29:24"},
	{"created-at":"2019-07-25T00:34:08","filename":"UOA-OH-INTEGRATION-TASK-pv69kw.json","id":787,"records":[],"task-type":"AFFILIATION","updated-at":"2019-07-25T01:32:36"}`)
			} else {
				io.WriteString(w, "[")
			}
			if withTasks && withAnIncomleteTask {
				io.WriteString(w, ",")
			}

			if withAnIncomleteTask {
				io.WriteString(w, `
	{"created-at":"2019-07-24T08:47:09","filename":"UOA-OH-INTEGRATION-TASK-this-should-get-activated.json","id":892,
		"records":[{}, {}],"task-type":"AFFILIATION","updated-at":"2019-07-24T09:29:24"},
	{"created-at":"2099-07-25T00:34:08","filename":"UOA-OH-INTEGRATION-TASK-pv69kZ.json","id":888,"records":[{},{}],"task-type":"AFFILIATION"}`)
			}
			io.WriteString(w, "]")
		case strings.HasPrefix(ru, "/api/v1/tokens/"):
			var id = strings.TrimPrefix(ru, "/api/v1/tokens/")
			if id == "rad42@mailinator.com" || id == "0000-0001-8228-7153" || id == "rcir178@auckland.ac.nz" {
				io.WriteString(w, `[{
				"access_token":"ecf16b31-ad54-4ba2-ae55-e97fb90e211a", 
				"email":"rad42@mailinator.com", 
				"eppn":"443469635@auckland.ac.nz", 
				"expires_in":631138518, 
				"issue_time":"2019-07-18T03:13:35", 
				"orcid":"0000-0001-8228-7153", 
				"refresh_token":"a6c9da20-31be-442a-9faa-73f1d92fac45",
				"scopes":"/read-limited,/activities/update"
			}]`)
			} else if id == "roshan.pawar@auckland.ac.nz" || id == "0000-0003-1255-9023" || id == "rpaw053@auckland.ac.nz" {
				io.WriteString(w, `[{
				"access_token": "89537453-7d6c-41a1-b619-895374524b76",
				"expires_in": 631138518,
				"issue_time": "2033-08-01T08:00:13",
				"refresh_token": "89537459-f057-477a-8a57-538953745f37",
				"scopes": "/read-limited,/activities/update",
				"email": "roshan.pawar@auckland.ac.nz",
				"eppn": "rpaw053@auckland.ac.nz",
				"orcid": "0000-0003-1255-9023"
			}]`)
			} else if id == "rcir178@auckland.ac.nz" {
				io.WriteString(w, `[]`)
			} else if id == "dthn666@auckland.ac.nz" {
				io.WriteString(w, `[{
				"access_token":"ecf16b31-7777-4ba2-ae55-e97fb90e211a", 
				"email":"dthn666@mailinator.com", 
				"eppn":"66666666@auckland.ac.nz", 
				"expires_in":631138518, 
				"issue_time":"2069-07-18T03:13:35", 
				"orcid":"0000-0001-8888-7153", 
				"refresh_token":"a6c9da20-31be-8888-9faa-73f1d92fac45",
				"scopes":"/read-limited"
			}]`)
			} else {
				w.WriteHeader(http.StatusNotFound)
				io.WriteString(w, `{"error": "User with specified identifier 'rcir178ABC@auckland.ac.nz' not found."}`)
			}
		case strings.HasPrefix(ru, "/api/v1/tasks/"):
			id := strings.TrimPrefix(ru, "/api/v1/tasks/")
			if r.Method == "POST" {
				filename := url.Query()["filename"][0]
				io.WriteString(w, `{
				"id":99999,
				"created-at":"2032-08-25T02:07:28",
				"filename":"`+filename+`",
				"task-type":"AFFILIATION",
				"updated-at":"2032-07-25T02:23:32"
			}`)
			} else if id == "12345" {
				io.WriteString(w, `{
				"created-at":"2032-08-25T02:07:28",
				"filename":"UOA-OH-INTEGRATION-TASK-pv6xi6.json",
				"id":`+id+`,
				"status":"ACTIVE",
				"task-type":"AFFILIATION",
				"updated-at":"2032-07-25T02:23:32",
				"records": [{}, {}]
			}`)
			} else if id == "54321" {
				io.WriteString(w, `{
				"created-at":"2002-08-25T02:07:28",
				"filename":"UOA-OH-INTEGRATION-TASK-this-should-get-activatged.json",
				"id":`+id+`,
				"task-type":"AFFILIATION",
				"updated-at":"2002-07-25T02:23:32",
				"records": [{}, {}]
			}`)
			} else {
				io.WriteString(w, `{
				"created-at":"2032-08-25T02:07:28",
				"filename":"UOA-OH-INTEGRATION-TASK-pv6xi6.json",
				"id":`+id+`,
				"status":"ACTIVE",
				"task-type":"AFFILIATION",
				"updated-at":"2032-07-25T02:23:32"
			}`)
			}
		case strings.HasPrefix(ru, "/api/v1/affiliations?filename="):
			var filename = strings.TrimPrefix(ru, "/api/v1/affiliations?filename=")
			io.WriteString(w, `{
				"created-at":"2019-07-25T02:23:32",
				"filename":"`+filename+`",
				"id":999,
				"records":[],
				"status":null,
				"task-type":"AFFILIATION",
				"updated-at":"2019-07-25T02:23:32"
			}`)
		case strings.HasPrefix(ru, "/api/v1/affiliations/"):
			var taskID = strings.TrimPrefix(ru, "/api/v1/affiliations/")
			io.WriteString(w, `{
				"id": `+taskID+`,
				"created-at": "2019-07-31T02:53:03",
				"filename": "UOA-OH-INTEGRATION-TASK-pvhk0f.json",
				"task-type": "AFFILIATION",
				"records": [
					{
						"id": 11441,
						"affiliation-type": "employment",
						"department": "Enterprise Architecture",
						"email": "radomirs.cirskis@auckland.ac.nz",
						"end-date": "2019-12-09",
						"orcid": "0000-0001-8228-7153",
						"role": "Project Architect",
						"start-date": "2018-08-09"
					},
					{
						"id": 11442,
						"affiliation-type": "employment",
						"department": "Enterprise Architecture",
						"email": "roshan.pawar@auckland.ac.nz",
						"end-date": "2019-11-15",
						"orcid": "0000-0003-1255-9023",
						"role": "Developer",
						"start-date": "2018-07-16"
					},
					{
						"id": 11443,
						"affiliation-type": "employment",
						"department": "Cent Learning \u0026 Rsch Higher Ed",
						"email": "roshan.pawar@auckland.ac.nz",
						"end-date": "2018-04-28",
						"orcid": "0000-0003-1255-9023",
						"role": "Professional Casual Staff",
						"start-date": "2016-06-15"
					},
					        {
						"id": 11449,
						"affiliation-type": "employment",
						"department": "Info Systems \u0026 Operations Mgmt",
						"email": "daniel.jimenez@auckland.ac.nz",
						"end-date": "2019-07-13",
						"orcid": "0000-0002-3008-0422",
						"role": "Teaching Assistant",
						"start-date": "2019-03-04"
					},
					{
						"id": 11448,
						"affiliation-type": "employment",
						"department": "Enterprise Architecture",
						"email": "daniel.jimenez@auckland.ac.nz",
						"orcid": "0000-0002-3008-0422",
						"role": "Professional Casual Staff",
						"start-date": "2017-03-01"
					},
					{
						"id": 11447,
						"affiliation-type": "employment",
						"department": "App Dev and QA",
						"email": "daniel.jimenez@auckland.ac.nz",
						"end-date": "2017-03-03",
						"orcid": "0000-0002-3008-0422",
						"role": "Intern",
						"start-date": "2016-11-16"
					},
					{
						"id": 11446,
						"affiliation-type": "employment",
						"department": "Enterprise Architecture",
						"email": "radomirs.cirskis@auckland.ac.nz",
						"end-date": "2019-12-09",
						"orcid": "0000-0001-8228-7153",
						"role": "Project Architect",
						"start-date": "2018-08-09"
					},
					{
						"id": 11445,
						"affiliation-type": "employment",
						"department": "Enterprise Architecture",
						"email": "roshan.pawar@auckland.ac.nz",
						"end-date": "2019-11-15",
						"orcid": "0000-0003-1255-9023",
						"role": "Developer",
						"start-date": "2018-07-16"
					},
					{
						"id": 11444,
						"affiliation-type": "employment",
						"department": "Cent Learning \u0026 Rsch Higher Ed",
						"email": "roshan.pawar@auckland.ac.nz",
						"end-date": "2018-04-28",
						"orcid": "0000-0003-1255-9023",
						"role": "Professional Casual Staff",
						"start-date": "2016-06-15"
					}
				]
			}`)
		case strings.HasPrefix(ru, "/service/student/integrations/v1/student/"):
			parts := strings.Split(ru, "/")
			uid := parts[6]
			switch uid {
			case "dthn666", "66666666", "dthn777", "77777777", "484378182", "4306445":
				io.WriteString(w, `[]`)
			case "jken016", "8524255":
				io.WriteString(w, `[
					{
						"id": "8524255",
						"studentDegNbr": "01",
						"degreeCode": "MSC-DG",
						"degreeDesc": "MSc",
						"degAcadCareer": "UC01",
						"degreeConferDate": "1990-05-03T12:00:00.000Z",
						"honorsPrefix": " ",
						"honorsSuffix": "H1",
						"degAcadDegreeStatus": "A",
						"prospectusCode": "AU0408",
						"degreePlans": [
						{
							"acadPlanCode": "PSYC-MSC",
							"acadPlanDesc": "Psychology",
							"dgpAcadCareer": "UC01",
							"studentCareerNbr": 1,
							"dgpAcadDegreeStatus": "A",
							"degreeStatusDate": "2001-02-22T11:00:00.000Z",
							"acadProgCode": "MSC",
							"acadProgGroupCode": 30,
							"acadProgGroup": "Degree",
							"acadProgLevelCode": "40",
							"acadProgLevel": "Postgraduate",
							"acadOrgCode": "SCIFAC",
							"acadGroupDesc": "Science"
						}
						]
					},
					{
						"id": "8524255",
						"studentDegNbr": "02",
						"degreeCode": "BSC-DG",
						"degreeDesc": "BSc",
						"degAcadCareer": "UC01",
						"degreeConferDate": "1989-05-03T12:00:00.000Z",
						"honorsPrefix": " ",
						"honorsSuffix": " ",
						"degAcadDegreeStatus": "A",
						"prospectusCode": "AU0087",
						"degreePlans": [
						{
							"acadPlanCode": "PSYC-BSC",
							"acadPlanDesc": "Psychology",
							"dgpAcadCareer": "UC01",
							"studentCareerNbr": 2,
							"dgpAcadDegreeStatus": "A",
							"degreeStatusDate": "2001-02-22T11:00:00.000Z",
							"acadProgCode": "BSC",
							"acadProgGroupCode": 30,
							"acadProgGroup": "Degree",
							"acadProgLevelCode": "20",
							"acadProgLevel": "Undergraduate",
							"acadOrgCode": "SCIFAC",
							"acadGroupDesc": "Science"
						}
						]
					}
				]`)
			case "208013283", "rpaw053":
				io.WriteString(w, `[
					{
						"id": "208013283",
						"studentDegNbr": "01",
						"degreeCode": "MESTU-DG",
						"degreeDesc": "MEngSt",
						"degAcadCareer": "UC01",
						"degreeConferDate": "2016-09-26T11:00:00.000Z",
						"honorsPrefix": " ",
						"honorsSuffix": " ",
						"degAcadDegreeStatus": "A",
						"prospectusCode": "AU4067",
						"degreePlans": [
							{
								"acadPlanCode": "SOFT-MESTU",
								"acadPlanDesc": "Software Engineering",
								"dgpAcadCareer": "UC01",
								"studentCareerNbr": 0,
								"dgpAcadDegreeStatus": "A",
								"degreeStatusDate": "2016-10-02T11:00:00.000Z",
								"acadProgCode": "MESTU",
								"acadProgGroupCode": 30,
								"acadProgGroup": "Degree",
								"acadProgLevelCode": "40",
								"acadProgLevel": "Postgraduate",
								"acadOrgCode": "ENGFAC",
								"acadGroupDesc": "Engineering"
							}
						]
					},
					{
						"id": "208013283",
						"studentDegNbr": "02",
						"degreeCode": "MESTU-DG",
						"degreeDesc": "MEngSt-NOT-EXISTING",
						"degAcadCareer": "UC01",
						"degreeConferDate": "2016-09-26T11:00:00.000Z",
						"honorsPrefix": " ",
						"honorsSuffix": " ",
						"degAcadDegreeStatus": "A",
						"prospectusCode": "AU4067",
						"degreePlans": [
							{
								"acadPlanCode": "SOFT-MESTU",
								"acadPlanDesc": "Software Engineering",
								"dgpAcadCareer": "UC01",
								"studentCareerNbr": 0,
								"dgpAcadDegreeStatus": "A",
								"degreeStatusDate": "2016-10-02T11:00:00.000Z",
								"acadProgCode": "MESTU",
								"acadProgGroupCode": 30,
								"acadProgGroup": "Degree",
								"acadProgLevelCode": "40",
								"acadProgLevel": "Postgraduate",
								"acadOrgCode": "ENGFAC",
								"acadGroupDesc": "Engineering"
							}
						]
					}
				]`)
			case "477579437", "djim087":
				io.WriteString(w, `[
					{
						"id": "477579437",
						"studentDegNbr": "01",
						"degreeCode": "BEHON-DG",
						"degreeDesc": "BE(Hons)",
						"degAcadCareer": "UC01",
						"degreeConferDate": "2019-05-02T12:00:00.000Z",
						"honorsPrefix": " ",
						"honorsSuffix": "H22",
						"degAcadDegreeStatus": "A",
						"prospectusCode": " ",
						"degreePlans": [
							{
								"acadPlanCode": "ELEC-BEHON",
								"acadPlanDesc": "Electrical and Electronic Eng",
								"dgpAcadCareer": "UC01",
								"studentCareerNbr": 0,
								"dgpAcadDegreeStatus": "A",
								"degreeStatusDate": "2019-05-09T12:00:00.000Z",
								"acadProgCode": "BEHON",
								"acadProgGroupCode": 30,
								"acadProgGroup": "Degree",
								"acadProgLevelCode": "20",
								"acadProgLevel": "Undergraduate",
								"acadOrgCode": "ENGFAC",
								"acadGroupDesc": "Engineering"
							}
						]
					}
				]`)
			default:
				w.WriteHeader(http.StatusBadRequest)
				io.WriteString(w, `{
					"timestamp": "2019-09-13T04:21:11.049Z",
					"status": 400,
					"error": "Bad Request",
					"message": "Incorrect or not supported id: `+uid+`",
					"path": "/student/`+uid+`/degree/"
				}`)
			}

		case strings.HasPrefix(ru, "/service/identity/integrations/v3/identity/"):
			var uid = strings.TrimPrefix(ru, "/service/identity/integrations/v3/identity/")
			switch uid {
			case "dthn666", "66666666":
				io.WriteString(w, `{}`)
			case "dthn777", "77777777":
				io.WriteString(w, `{"upi":"dthn777", "id":77777777}`)
			case "jken016", "8524255":
				io.WriteString(w, `{
    "emailAddress": "jeff.kennedy@auckland.ac.nz",
    "emails": [
        {
            "email": "jeff_is_dead_at_last_in_spite_of_all@hotmail.com",
            "lastUpdated": "2016-03-17T21:06:21.000+0000",
            "type": "Other",
            "typeId": "Other",
            "verified": false
        },
        {
            "email": "jken016@aucklanduni.ac.nz",
            "lastUpdated": "2015-07-27T22:16:33.000+0000",
            "type": "Student",
            "typeId": "Student",
            "verified": true
        },
        {
            "email": "honeylarkin@gmail.com",
            "lastUpdated": "2015-07-27T22:16:33.000+0000",
            "type": "Work",
            "typeId": "Business",
            "verified": true
        },
        {
            "email": "jeff.kennedy@auckland.ac.nz",
            "lastUpdated": "2016-03-17T21:06:21.000+0000",
            "type": "University",
            "typeId": "Campus",
            "verified": true
        },
        {
            "email": "belacqua66@hotmail.com",
            "lastUpdated": "2019-03-01T06:10:56.000+0000",
            "type": "Personal",
            "typeId": "Home",
            "verified": false
        }
    ],
    "extIds": [
        {
            "id": "38713",
            "type": "Advancement"
        },
        {
            "id": "20517852425502",
            "type": "IDCard"
        },
        {
            "id": "138256828",
            "type": "NSN"
        },
        {
            "id": "http://orcid.org/0000-0002-8982-6444",
            "type": "ORCID"
        },
        {
            "id": "23817",
            "type": "UID"
        }
    ],
    "id": 8524255,
    "upi": "jken016"
}`)
			case "4306445", "yyan161":
				io.WriteString(w, `{
    "emailAddress": "jasmine_yinyin@hotmail.com",
    "emails": [
        {
            "email": "jasmine_yinyin@hotmail.com",
            "lastUpdated": "2018-08-27T23:51:28.000+0000",
            "type": "Other",
            "typeId": "Other",
            "verified": false
        }
    ],
    "extIds": [
        {
            "id": "118009562",
            "type": "NSN"
        },
        {
            "id": "340776",
            "type": "UID"
        }
    ],
    "id": 4306445,
    "upi": "yyan161"
	}`)
			case "rpaw053", "208013283":
				io.WriteString(w, `{
   "emailAddress":"roshan_pawarasjdfkasdjfajs_@auckland.ac.nz",
   "emails":[
      {
         "email":"rwrrwe3343@auckland.ac.nz",
         "lastUpdated":"2017-01-13T17:12:23.000+0000",
         "typeId":"Campus",
         "type":"University",
         "verified":false
      },
      {
         "email":"rpfds434@aucklanduni.ac.nz",
         "lastUpdated":"2017-01-13T17:12:24.000+0000",
         "typeId":"Student",
         "type":"Student",
         "verified":true
      },
      {
         "email":"getconfjsdlkajfalhan@gmail.com",
         "lastUpdated":"2017-01-13T17:12:24.000+0000",
         "typeId":"Business",
         "type":"Work",
         "verified":false
      }
   ],
   "extIds":[
      {
         "id":"2121820801328312",
         "type":"IDCard"
      },
      {
         "id":"149928464",
         "type":"NSN"
      },
      {
         "id":"http://orcid.org/0000-0002-9398-4322",
         "type":"ORCID"
      },
      {
         "id":"2490528",
         "type":"UID"
      }
   ],
   "id":208013283,
   "upi":"rpaw053",
   "previousIds":[]}`)
			case "rcir178", "484378182":
				io.WriteString(w, `{
    "emailAddress": "sjdfkjd9444353@auckland.ac.nz",
    "emails": [
		{
			"email": "rrrr4353@auckland.ac.nz",
			"lastUpdated": "2017-08-24T23:25:18.000+0000",
			"type": "University",
			"typeId": "Campus",
			"verified": false
		},
		{
			"email": "nad2000@gmail.com",
			"lastUpdated": "2017-08-24T23:25:22.000+0000",
			"type": "Other",
			"typeId": "Other",
			"verified": true
		},
		{
			"email": "rad@nowitworks.eu",
			"lastUpdated": "2017-08-24T23:25:22.000+0000",
			"type": "Work",
			"typeId": "Business",
			"verified": false
		}
    ],
    "extIds": [
		{
			"id": "2011948437818225",
			"type": "IDCard"
		},
		{
			"id": "154244310",
			"type": "NSN"
		},
		{
			"id": "2594016",
			"type": "UID"
		}
    ],
    "id": 484378182,
    "mobile": "+64221221442",
    "upi": "rcir178"
}`)
			case "477579437", "djim087":
				io.WriteString(w, `{
    "emailAddress": "daniel.jimenez@auckland.ac.nz",
    "emails": [
        {
            "email": "daniel.jimenez@auckland.ac.nz",
            "lastUpdated": "2017-05-05T03:31:13.000+0000",
            "type": "University",
            "typeId": "Campus",
            "verified": false
        },
        {
            "email": "dan.kiwi@live.com",
            "lastUpdated": "2017-05-05T03:31:14.000+0000",
            "type": "Other",
            "typeId": "Other",
            "verified": true
        },
        {
            "email": "djim087@aucklanduni.ac.nz",
            "lastUpdated": "2017-05-05T03:31:15.000+0000",
            "type": "Student",
            "typeId": "Student",
            "verified": false
        }
    ],
    "extIds": [
        {
            "id": "2121847757943760",
            "type": "IDCard"
        },
        {
            "id": "130768622",
            "type": "NSN"
        },
        {
            "id": "2456801",
            "type": "UID"
        }
    ],
	"id": 477579437,
    "upi": "djim087"
}`)

			case "rad42", "non-existing-upi-error":
				t.Log("NOT FOUND .... ", uid)
				w.WriteHeader(http.StatusNotFound)
				io.WriteString(w, `{
					"timestamp":"2019-07-25T06:34:50.211+0000",
					"status":404,
					"error":"Not Found",
					"message":"Identity not found","path":"/identity/`+uid+`"
				}`)
			default:
				w.WriteHeader(http.StatusBadRequest)
				io.WriteString(w, `{"timestamp":"2019-07-25T02:23:32.668+0000","status":400,"error":"Bad Request","message":"Incorrect or not supported id","path":"/identity/`+uid+`"}`)
			}
		case strings.HasPrefix(ru, "/service/employment/integrations/v1/employee/"):
			var upiOrID = strings.TrimPrefix(ru, "/service/employment/integrations/v1/employee/")
			switch upiOrID {
			case "477579437", "djim087":
				io.WriteString(w, `{
    "employeeID":"477579437",
    "professionalStaffFTE":0,
    "academicStaffFTE":0,
    "uniServicesFTE":0,"requestTimeStamp":"2019-07-24T03:40:53.000Z",
    "job":[
        {"employeeRecord":0,"effectiveDate":"2017-03-04","effectiveSequence":0,"organizationalRelation":"EMP","departmentID":"ITGADEVQA","departmentDescription":"App Dev and QA","jobCode":"H00055","jobGrade":"G4S","positionNumber":"55561662","positionDescription":"Intern","hrStatus":"I","employeeStatus":"T","lastHRaction":"TER","location":"435","locationDescription":"58 Symonds Street","standardHours":37.5,"employeeType":"Fixed Term","salAdminPlan":"GS1","fullTimeEquivalent":1,"jobIndicator":"S","supervisorID":"","poiType":"","jobStartDate":"2016-11-16","jobEndDate":"2017-03-03","jobCodeDescription":"Analyst/Developer","parentDepartmentDescription":"Application Development and Quality Assurance","primaryActivityCentreDeptID":"","primaryActivityCentreDeptDescription":"","reportsToPosition":"","company":"UOA","costCentre":"8854","updatedDateTime":"2017-03-03T11:10:40.000Z"},
        {"employeeRecord":1,"effectiveDate":"2019-07-15","effectiveSequence":0,"organizationalRelation":"EMP","departmentID":"ITARCHIT","departmentDescription":"Enterprise Architecture","jobCode":"B00029","jobGrade":"G4S","positionNumber":"00005285","positionDescription":"Professional Casual Staff","hrStatus":"A","employeeStatus":"A","lastHRaction":"DTA","location":"409","locationDescription":"Information Technology Centre","standardHours":0.01,"employeeType":"Casual","salAdminPlan":"GS1","fullTimeEquivalent":0,"jobIndicator":"P","supervisorID":"8524255","poiType":"","jobStartDate":"2017-03-01","jobCodeDescription":"Professional Casual Staff","parentDepartmentDescription":"Enterprise Architecture","primaryActivityCentreDeptID":"CDO","primaryActivityCentreDeptDescription":"Chief Digital Officer's Office","reportsToPosition":"00012578","company":"UOA","costCentre":"8848","updatedDateTime":"2019-07-15T01:59:49.000Z"},
        {"employeeRecord":2,"effectiveDate":"2019-07-14","effectiveSequence":0,"organizationalRelation":"EMP","departmentID":"ISOM","departmentDescription":"Info Systems & Operations Mgmt","jobCode":"A00363","jobGrade":"TAS","positionNumber":"00009299","positionDescription":"Teaching Assistant","hrStatus":"I","employeeStatus":"T","lastHRaction":"TER","location":"260","locationDescription":"Owen G Glenn Building","standardHours":0.01,"employeeType":"Casual","salAdminPlan":"AS1","fullTimeEquivalent":0,"jobIndicator":"S","supervisorID":"2450582","poiType":"","jobStartDate":"2019-03-04","jobEndDate":"2019-07-13","jobCodeDescription":"Teaching Assistant","parentDepartmentDescription":"Information Systems and Operations Management","primaryActivityCentreDeptID":"BUSEC","primaryActivityCentreDeptDescription":"Business and Economics","reportsToPosition":"55560561","company":"UOA","costCentre":"1545","updatedDateTime":"2019-07-15T01:59:48.000Z"}
    ]
}`)
			case "208013283", "rpaw053":
				io.WriteString(w, `{
    "employeeID": "208013283",
    "professionalStaffFTE": 1,
    "academicStaffFTE": 0,
    "uniServicesFTE": 0,
    "requestTimeStamp": "2019-07-31T01:39:11.000Z",
    "job": [
        {
            "employeeRecord": 0,
            "effectiveDate": "2018-04-29",
            "effectiveSequence": 0,
            "organizationalRelation": "EMP",
            "departmentID": "CLEAR",
            "departmentDescription": "Cent Learning & Rsch Higher Ed",
            "jobCode": "B00029",
            "jobGrade": "G2S",
            "positionNumber": "00004741",
            "positionDescription": "Professional Casual Staff",
            "hrStatus": "I",
            "employeeStatus": "T",
            "lastHRaction": "TER",
            "location": "804",
            "locationDescription": "Fisher Building",
            "standardHours": 0.01,
            "employeeType": "Casual",
            "salAdminPlan": "GS1",
            "fullTimeEquivalent": 0,
            "jobIndicator": "S",
            "supervisorID": "8986011",
            "poiType": "",
            "jobStartDate": "2016-06-15",
            "jobEndDate": "2018-04-28",
            "jobCodeDescription": "Professional Casual Staff",
            "parentDepartmentDescription": "Centre for Learning and Research in Higher Education",
            "primaryActivityCentreDeptID": "EDUFAC",
            "primaryActivityCentreDeptDescription": "Education and Social Work",
            "reportsToPosition": "55561014",
            "company": "UOA",
            "costCentre": "7000",
            "updatedDateTime": "2018-05-02T19:54:47.000Z"
        },
        {
            "employeeRecord": 1,
            "effectiveDate": "2019-02-01",
            "effectiveSequence": 0,
            "organizationalRelation": "EMP",
            "departmentID": "ITARCHIT",
            "departmentDescription": "Enterprise Architecture",
            "jobCode": "H00028",
            "jobGrade": "G5S",
            "positionNumber": "55561720",
            "positionDescription": "Developer",
            "hrStatus": "A",
            "employeeStatus": "A",
            "lastHRaction": "PAY",
            "location": "435",
            "locationDescription": "58 Symonds Street",
            "standardHours": 37.5,
            "employeeType": "Fixed Term",
            "salAdminPlan": "GS1",
            "fullTimeEquivalent": 1,
            "jobIndicator": "P",
            "supervisorID": "8524255",
            "poiType": "",
            "jobStartDate": "2018-07-16",
            "jobEndDate": "2019-11-15",
            "jobCodeDescription": "IT Analyst",
            "parentDepartmentDescription": "Enterprise Architecture",
            "primaryActivityCentreDeptID": "CDO",
            "primaryActivityCentreDeptDescription": "Chief Digital Officer's Office",
            "reportsToPosition": "00012578",
            "company": "UOA",
            "costCentre": "8848",
            "updatedDateTime": "2019-02-06T23:33:42.000Z"
        }
    ]
}`)
			case "484378182", "rcir178":
				io.WriteString(w, `{
    "employeeID": "484378182",
    "professionalStaffFTE": 1,
    "academicStaffFTE": 0,
    "uniServicesFTE": 0,
    "requestTimeStamp": "2019-07-31T01:40:45.000Z",
    "job": [
        {
            "employeeRecord": 0,
            "effectiveDate": "2019-02-01",
            "effectiveSequence": 0,
            "organizationalRelation": "EMP",
            "departmentID": "ITARCHIT",
            "departmentDescription": "Enterprise Architecture",
            "jobCode": "H00028",
            "jobGrade": "G6S",
            "positionNumber": "55561722",
            "positionDescription": "Project Architect",
            "hrStatus": "A",
            "employeeStatus": "A",
            "lastHRaction": "PAY",
            "location": "435",
            "locationDescription": "58 Symonds Street",
            "standardHours": 37.5,
            "employeeType": "Fixed Term",
            "salAdminPlan": "GS1",
            "fullTimeEquivalent": 1,
            "jobIndicator": "P",
            "supervisorID": "8524255",
            "poiType": "",
            "jobStartDate": "2018-08-09",
            "jobEndDate": "2019-12-09",
            "jobCodeDescription": "IT Analyst",
            "parentDepartmentDescription": "Enterprise Architecture",
            "primaryActivityCentreDeptID": "CDO",
            "primaryActivityCentreDeptDescription": "Chief Digital Officer's Office",
            "reportsToPosition": "00012578",
            "company": "UOA",
            "costCentre": "8848",
            "updatedDateTime": "2019-02-06T23:56:56.000Z"
        }
    ]
}`)
			default:
				if isValidUPI(upiOrID) || isValidID(upiOrID) {
					w.WriteHeader(http.StatusNotFound)
					io.WriteString(w, `{
				"timestamp": "2029-07-31T01:30:15.565Z",
				"status": 404,
				"error": "Not Found",
				"exception": "nz.ac.auckland.exceptions.ApiException",
				"message": "User is not found in LDAP",
				"path": "/employment/v1/employee/`+upiOrID+`"
			}`)
				} else {
					w.WriteHeader(http.StatusBadRequest)
					io.WriteString(w, `{
				"timestamp": "2019-07-31T01:34:06.957Z",
				"status": 400,
				"error": "Bad Request",
				"exception": "nz.ac.auckland.exceptions.ApiException",
				"message": "Incorrect or not supported id: `+upiOrID+`",
				"path": "/employment/v1/employee/`+upiOrID+`"
			}`)
				}
			}
		case ru == "/external-organisations/v1/qualifications":
			io.WriteString(w, `[
  {
    "type": "tertiary",
    "code": "DPURB-DP",
    "description": "Diploma in Urban Valuation",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MED-DG",
    "description": "Master of Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BPLAN-DG",
    "description": "Bachelor of Planning",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MPLAN-DG",
    "description": "Master of Planning",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MLITT-DG",
    "description": "Master of Literature",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "LLD-DG",
    "description": "Doctor of Laws",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DJUR-DG",
    "description": "Doctor of Jurisprudence",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPDRA-DP",
    "description": "Diploma in Drama",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "LG",
    "description": "NZ Certificate in Computer Technology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "LX",
    "description": "National Certificate in Business Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MOR-DG",
    "description": "Master of Operations Research",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "CTHSC-CT",
    "description": "Certificate in Health Sciences",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPTED-DP",
    "description": "Diploma in Teacher Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPMUS-DP",
    "description": "Diploma in Music",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GDMUA-DP",
    "description": "Graduate Diploma in Music",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MUD-DG",
    "description": "Master of Urban Design",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BSW-DG",
    "description": "Bachelor of Social Work",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GDSW-DP",
    "description": "Graduate Diploma in Social Work",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "CTYS-CT",
    "description": "Certificate in Youth Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PROFL",
    "description": "Professional Legal Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDET-DP",
    "description": "Postgraduate Diploma in Education (Teaching)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GDCOM-DP",
    "description": "Graduate Diploma in Commerce",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MCPA-DG",
    "description": "Master of Creative and Performing Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDMS-DP",
    "description": "Postgraduate Diploma in Medical Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BTHEC-DG",
    "description": "Bachelor of Theology (Conjoint)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDHP-DP",
    "description": "Postgraduate Diploma in Health Psychology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MPHPR-DG",
    "description": "Master of Pharmacy Practice",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPART-DP",
    "description": "Diploma in Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDSC-DP",
    "description": "Postgraduate Diploma in Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDSR-DP",
    "description": "Postgraduate Diploma in Social Research",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDPR-DP",
    "description": "Postgraduate Diploma in Property",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPTSE-DP",
    "description": "Diploma in Teaching (Secondary)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPOMG-DP",
    "description": "Diploma in Obstetrics and Medical Gynaecology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPMNH-DP",
    "description": "Diploma in Mental Health",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "CTGET-CT",
    "description": "Certificate in Geothermal Energy Techn",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "TO",
    "description": "Bachelor of Law/Bachelor of Commerce Hons",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "UE",
    "description": "Bachelor of Medicine and Surgery",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "UF",
    "description": "Master of Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "VB",
    "description": "Law Professionals",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "WC",
    "description": "Postgraduate Diploma in Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "WV",
    "description": "Unitech Certificate",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "YD",
    "description": "Doctor of Philosophy",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDEL-DP",
    "description": "Postgraduate Diploma in Educational Leadership",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "KA",
    "description": "Association of Chartered Accountants",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "XB",
    "description": "B.C.A:",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "WO",
    "description": "Interest Only",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC",
    "description": "Licence of English",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "XE",
    "description": "work at Teachers College",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EBBH",
    "description": "Bachelor of Engineering Technology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EDDV",
    "description": "Diploma in Maori Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EDED",
    "description": "Diploma in Visual Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EGDFO",
    "description": "Graduate Diploma in Journalism",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ENCJS",
    "description": "National Diploma in Journalism",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ENCJV",
    "description": "National Diploma in Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EHBC",
    "description": "Bachelor of Dental Technology (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EHBK",
    "description": "Bachelor of Health Science (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EBCB",
    "description": "Bachelor of Performing and Screen Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EHCG",
    "description": "Bachelor of Social Work (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ECDA",
    "description": "Certificate in Sports Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EDDH",
    "description": "Diploma in Applied Interior Design",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EBCO",
    "description": "Bachelor of Veterinary Technology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMGL",
    "description": "Master of Business and Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMHB",
    "description": "Master of Entrepreneurship",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPMED-DP",
    "description": "Diploma in Mathematics Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPSCE-DP",
    "description": "Diploma in Science Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NCFLM-CT",
    "description": "National Certificate in First Line Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGCLL-CT",
    "description": "Postgraduate Certificate in Law",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPCCM-DP",
    "description": "Diploma in Care Co-ordination and Management (Intellectual Disability)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "HCTLE-CT",
    "description": "Higher Certificate in Language and Literacy",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPIT-DP",
    "description": "Diploma of Information Technology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "COPEX-CT",
    "description": "Certificate of Proficiency for Exchange",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDPS-DP",
    "description": "Postgraduate Diploma in Professional Supervision",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPDME",
    "description": "Postgraduate Diploma in Ministry",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPCKT",
    "description": "Postgraduate Certificate in Midwifery",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPCLE",
    "description": "Postgraduate Certificate in Teaching",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPDMR",
    "description": "Postgraduate Diploma in Social Welfare",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPDLW",
    "description": "Postgraduate Diploma in Higher Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPDMH",
    "description": "Postgraduate Diploma in Performance and Media Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPCKJ",
    "description": "Postgraduate Certificate",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPDNC",
    "description": "Postgraduate Diploma in Clinical Research",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0218",
    "description": "National Certificate in Service Sector",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0224",
    "description": "National Certificate in Solid Waste",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0227",
    "description": "National Certificate in Specialist Rescue",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0228",
    "description": "National Certificate in Sport Officiating",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0234",
    "description": "National Certificate in Surveying",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "C00018",
    "description": "Certificate in Language Teaching",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "C00021",
    "description": "Certificate in Music",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "D00004",
    "description": "Diploma in Environmental Technology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0003",
    "description": "National Certificate in Adult Education and Training",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0017",
    "description": "National Certificate in Animal Care",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0026",
    "description": "National Certificate in Beauty Services",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0029",
    "description": "National Certificate in Blaster Coating",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0045",
    "description": "National Certificate in Civil Defense",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0047",
    "description": "National Certificate in Civil Infrastructure Health Safety and Environment",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0063",
    "description": "National Certificate in Core Rigging",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0067",
    "description": "National Certificate in Credit Administration",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0090",
    "description": "National Certificate in Employment Skills",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0108",
    "description": "National Certificate in Food and Related Processing",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0132",
    "description": "National Certificate in Industrial Machine",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0133",
    "description": "National Certificate in Industrial Measurement and Control",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0137",
    "description": "National Certificate in Infrastructure Works",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0148",
    "description": "National Certificate in Maori (Te Ngutu Awa)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0166",
    "description": "National Certificate in Music",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0169",
    "description": "National Certificate in Offender Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0172",
    "description": "National Certificate in Pacific Island Social Services",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0175",
    "description": "National Certificate in Passive Fire Protection",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EDDQ",
    "description": "Diploma in Interior Design Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EGCET",
    "description": "Graduate Certificate in Design",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EGDFI",
    "description": "Graduate Diploma in Commerce",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EGDFM",
    "description": "Graduate Diploma in Engineering",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMGF",
    "description": "Master of Arts and Design",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMGG",
    "description": "Master of Arts Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPDLJ",
    "description": "Postgraduate Diploma in Architecture",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MINBS-DG",
    "description": "Master of International Business",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0237",
    "description": "National Certificate in Te Matauranga Maori me te Whakangungu ",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0242",
    "description": "National Certificate in Transportation",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0009",
    "description": "National Diploma in Applied Journalism",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0010",
    "description": "National Diploma in Aviation",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0040",
    "description": "National Diploma in Infrastructure Asset Mmgt",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0043",
    "description": "National Diploma in Joinery",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0050",
    "description": "National Diploma in Outdoor Recreation",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0060",
    "description": "National Diploma in Resource Efficiency",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0069",
    "description": "National Diploma in Sports Turf Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0073",
    "description": "National Diploma in Textiles Manufacture",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0080",
    "description": "National Diploma in Workplace Emergency Risk Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0081",
    "description": "National Diploma in Zero Waste and Res Recovery",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ECCW",
    "description": "Certificate in Foundation Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EFCEO",
    "description": "Foundation Certificate (any other certificate)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "COPWP-CT",
    "description": "Certificate of Proficiency - Winter Programme",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MPHL-DG",
    "description": "Master of Philosophy",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MMIB-DG",
    "description": "Master of Māori and Indigenous Business",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MARUD-DG",
    "description": "Master of Architecture (Professional) and Urban Design",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPPIE-DP",
    "description": "Diploma in Pacific Islands Early Childhood Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MBS-DG",
    "description": "Master of Biomedical Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "FCTTM-CT",
    "description": "Foundation Certificate Tohu Tuapapa Matauranga",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GDHSC",
    "description": "Graduate Diploma in Health Sciences",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MPHIL-DG",
    "description": "Master of Philosophy",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDPU-DP",
    "description": "Postgraduate Diploma in Public Policy",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BMIH-DG",
    "description": "Bachelor of Medical Imaging (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BNUBS-DG",
    "description": "BNurs/BSc Conjoint",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MMGT-DG",
    "description": "Master of Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MEDE-DG",
    "description": "Master of Education",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "CPFS",
    "description": "CPIT Certificate in Foundation Studies",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "MCPU",
    "description": "Certificate of University Preparation (Massey Uni)",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "FSBM",
    "description": "Foundation Studies (Business)",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "VUFS",
    "description": "Certificate in Foundation Studies",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "UOCP",
    "description": "Certificate in Preparation Studies",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "MITP",
    "description": "MIT Preparation courses all programmes",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "CPCP",
    "description": "CPIT Certificate of Proficiency Fdn",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "UOAF",
    "description": "Bradford College Foundation Studies Programme",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "VUCP",
    "description": "Certificate of University Preparation",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "TFSS",
    "description": "Foundation Social Science",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "LUCF",
    "description": "Certificate in Foundation Studies",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "SEC",
    "description": "Secondary Education Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "DHC",
    "description": "Diplome des Humanites Completes",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "TAY",
    "description": "Taylors Foundation Year",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "AGCS",
    "description": "Advanced General Cert of Secondary Education",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "SLCE",
    "description": "School Leaving Certificate Examination",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "TSC",
    "description": "Tonga School Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "BACC",
    "description": "Baccalaureat",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "TER",
    "description": "Tertiary Entrance Rank (TER)",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "KCE",
    "description": "Kenya Certificate of Education",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "EHEE",
    "description": "Ethiopean Higher Education Entrance Examination",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "FF7",
    "description": "Fiji Seventh Form Examination",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "ZAC",
    "description": "ZIMSEC A Level Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "HSSC",
    "description": "Higher Secondary School Certificates - India",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GV",
    "description": "Bachelor of Home Science/Bachelor of Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GW",
    "description": "Bachelor of Phys Ed/Bachelor of Commerce",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "HA",
    "description": "Registered Comprehensive Nurse",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "HV",
    "description": "Diploma in Music",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "IA",
    "description": "Diploma in Clinical Chemistry",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "IB",
    "description": "Diploma in Radiography",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ID",
    "description": "Diploma in Nuclear Medicine",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "IL",
    "description": "Diploma in Diagnostic Radiography",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "IX",
    "description": "Overseas Diploma (unclassified)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "IY",
    "description": "Diploma in Quality Assurance",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "IZ",
    "description": "Diploma in Rehabilitation",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "JD",
    "description": "Certificate in Graphic Design",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "JE",
    "description": "Technicians Certificate (Telephone)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "JI",
    "description": "Advanced Diploma in Business Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "JN",
    "description": "Executant Music Diploma",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "JQ",
    "description": "Certificate in Te Reo Maori",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "KO",
    "description": "Diploma in Teaching English as a 2nd Lang",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MM",
    "description": "Bachelor of Music",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MS",
    "description": "Bachelor of Medicine and Bachelor of Surgery",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NO",
    "description": "Diploma in Valuation & Property Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PG",
    "description": "Master of Commerce",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PM",
    "description": "Master of Social Welfare",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PT",
    "description": "Master of Horticultural Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "QK",
    "description": "Bachelor of Engineering (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "QT",
    "description": "Bachelor of Communication Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "RB",
    "description": "Bachelor of Human Biology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "RU",
    "description": "Bachelor of Property Administration",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "SI",
    "description": "Diploma in English Language Teaching",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "SK",
    "description": "Diploma in Agricultural Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ST",
    "description": "Diploma in Criminology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "TI",
    "description": "Bachelor of Law Honours/Bachelor of Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MHMGT-DG",
    "description": "Master of Health Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "FR",
    "description": "Diploma in Community Health",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "FT",
    "description": "Diploma in Developmental Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GE",
    "description": "Bachelor of Arts/Bachelor of Theology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "YM",
    "description": "Doctor of Literature",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "TBS",
    "description": "Titulo de Bachilleratio Scientifico",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "BEC",
    "description": "Basic Education Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "EDC",
    "description": "Diplome d'Ecole de Culture Generale",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "OTH",
    "description": "Other Secondary Qualification",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "ZSC",
    "description": "Zambia School Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "WAHC",
    "description": "West African Higher School Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "ERE",
    "description": "Erettsegi/Matura",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "ICSE",
    "description": "Indian Certificate of Secondary Education",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "BSSC",
    "description": "Bermuda Secondary School Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "FSLC",
    "description": "Federal Secondary School Leaving Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "VFG",
    "description": "Vitnemal fra Grunnskolen",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "FYTC",
    "description": "Fiji Year 13 Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "CES",
    "description": "Certificat d'Enseignement Secondaire",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "BDC",
    "description": "Brevet des Colleges",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "MAH",
    "description": "Reifprufung/Matura",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "CIE3",
    "description": "Cambridge International Examinations - A Levels",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "GCAS",
    "description": "GCE/CIE AS-Levels",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "GSEC",
    "description": "General Secondary Education Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "GCAS",
    "description": "General Certificate of Education (AS Level)",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "AUS8",
    "description": "Western Australia Certificate of Education (WACE)",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "COSC",
    "description": "Cambridge Overseas School Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "RSC",
    "description": "Religious Secondary Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "WSSC",
    "description": "Western Samoa School Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "FPSP",
    "description": "Foundation Programme (USP)",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "CPU",
    "description": "Cambridge Pre-U",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "KMT",
    "description": "Kantonale Maturität",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "VVS",
    "description": "Vitnemal fra den Vidergaende Skole",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPARC-DP",
    "description": "Diploma in Architecture",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPFAH-DP",
    "description": "Diploma in Fine Arts (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPOPT-DP",
    "description": "Diploma in Optometry",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MTP-DG",
    "description": "Master of Town Planning",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PHD-DG",
    "description": "Doctor of Philosophy",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BSCED-DG",
    "description": "Bachelor of Science Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BMUPH-DG",
    "description": "Bachelor of Music (Performance) (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BACJT-DG",
    "description": "Bachelor of Arts (Conjoint)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "LLBHC-DG",
    "description": "Bachelor of Laws (Honours) (Conjoint)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MARCH-DG",
    "description": "Master of Architecture",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "LLM-DG",
    "description": "Master of Laws",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "LH",
    "description": "Certificate in Applied Social Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "LJ",
    "description": "Certificate in Computer Engineering",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "LP",
    "description": "NZ Certificate in Land Surveying",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "LZ",
    "description": "Higher National Certificate",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "COPOS-CT",
    "description": "Certificate of Proficiency for Overseas University",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPENM-DP",
    "description": "Diploma in Environmental Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MENVL-DG",
    "description": "Master of Environmental Legal Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MCW-DG",
    "description": "Master of Creative Writing",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GCSPD-CT",
    "description": "Graduate Certificate in Supervision and Professional Development",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPESV-DP",
    "description": "Diploma in Education of Students with Vision Impairment",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPEDM-DP",
    "description": "Diploma of Education Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPEI-DP",
    "description": "Diploma in Early Intervention",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPDCE-DP",
    "description": "Diploma of Dance Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGCMS-CT",
    "description": "Postgraduate Certificate in Medical Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGCHS-CT",
    "description": "Postgraduate Certificate in Health Sciences",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDHM-DP",
    "description": "Postgraduate Diploma in Health Science (Mental Health)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DMLED-DP",
    "description": "Diploma in Mathematical Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDCM-DP",
    "description": "Postgraduate Diploma in Community Emergency Medicine",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPBRC-DP",
    "description": "Diploma in Broadcast Communication",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPPRE-DP",
    "description": "Diploma in Professional Ethics",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDLS-DP",
    "description": "Postgraduate Diploma in Legal Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDCP-DP",
    "description": "Postgraduate Diploma in Clinical Psychology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPSPE-DP",
    "description": "Diploma in Special Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPCSC-DP",
    "description": "Diploma in Computer Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPCEM-DP",
    "description": "Diploma in Community Emergency Medicine",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "UA",
    "description": "Master of Business Administration",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "UB",
    "description": "Master of Fine Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "UH",
    "description": "Master of Forestry Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "UT",
    "description": "Master of Social Sciences",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "UV",
    "description": "Master of Regional Resource Planning",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "UX",
    "description": "Overseas Masters Degree (unclassified)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "WD",
    "description": "Postgraduate Diploma in Business Administration",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GDLAW-DP",
    "description": "Graduate Diploma in Law",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPTES-DP",
    "description": "Diploma of TESSOL",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MARPF-DG",
    "description": "Master of Architecture (Professional)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BSWH-DG",
    "description": "Bachelor of Social Work (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BDANS-DG",
    "description": "Bachelor of Dance Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMHD",
    "description": "Master of Environmental Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMHF",
    "description": "Master of Food Technology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMII",
    "description": "Master of Resource and Environmental Planning",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMIL",
    "description": "Master of Social Work",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMIS",
    "description": "Master of Theatre Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EDCEG",
    "description": "Doctor of Clinical Dentistry",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMHC",
    "description": "Master of Environmental Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMHY",
    "description": "Master of Nursing",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ENCJH",
    "description": "National Certificate in Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ENCJR",
    "description": "National Diploma in Hospitality Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ENCJT",
    "description": "National Diploma in Quantity Surveying",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ENCJX",
    "description": "National Diploma in Surveying",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EBAC",
    "description": "Bachelor of Applied Economics",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EBAI",
    "description": "Bachelor of Aviation",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EBAX",
    "description": "Bachelor of Creative Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EBBB",
    "description": "Bachelor of Dental Technology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EBBY",
    "description": "Bachelor of Natural Sciences",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EHCI",
    "description": "Bachelor of Surveying (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EBAB",
    "description": "Bachelor of Accountancy",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EBAD",
    "description": "Bachelor of Applied Economics (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EDDL",
    "description": "Diploma in Contemporary Craft",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EDDN",
    "description": "Diploma in Contemporary Photography",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GF",
    "description": "Bachelor of Arts/Cert of Proficiency",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "TS",
    "description": "Bachelor of Law/Cert. of Proficiency",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "FE",
    "description": "Higher Diploma of Teaching",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "SA",
    "description": "Diploma in Dairy Technology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "WJ",
    "description": "Doctor of Engineering",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MSWP-DG",
    "description": "Master of Social Work (Professional)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BDSHN-DG",
    "description": "Bachelor of Dance Studies (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GCIE-CT",
    "description": "Graduate Certificate in Innovation and Entrepreneurship",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EBCP",
    "description": "Bachelor of Visual Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMGS",
    "description": "Master of Creative Technologies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMGW",
    "description": "Master of Dental Technology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMGZ",
    "description": "Master of Electronic Commerce",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NCHUS-CT",
    "description": "National Certificate in Human Services",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "HCTEM-CT",
    "description": "Higher Certificate in Educational Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MCLED-DG",
    "description": "Master of Clinical Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GDSED-DP",
    "description": "Graduate Diploma in Special Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BFAH-DG",
    "description": "Bachelor of Fine Arts (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "CTSCI-CT",
    "description": "Certificate in Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BFAHC-DG",
    "description": "Bachelor of Fine Arts (Honours) (Conjoint)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPDMA",
    "description": "Postgraduate Diploma in Landscape Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPDMQ",
    "description": "Postgraduate Diploma in Social Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EDKE",
    "description": "NZIM Diploma in Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPDMS",
    "description": "Postgraduate Diploma in Social Work",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPCKL",
    "description": "Postgraduate Certificate in Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPDMK",
    "description": "Postgraduate Diploma in Physiotherapy",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPDMO",
    "description": "Postgraduate Diploma in Rehabilitation",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPDNA",
    "description": "Postgraduate Diploma of Computer Graphic Design",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGCMT-CT",
    "description": "Postgraduate Certificate in Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MPH",
    "description": "Master of Health Psychology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPDMU",
    "description": "Postgraduate Diploma in Sport and  Exercise",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0232",
    "description": "National Certificate in Stevedoring",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "C00003",
    "description": "Certificate in Animal Welfare Investigations",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "C00008",
    "description": "Certificate in Communication and Media Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "C00012",
    "description": "Certificate in Employment and Community Skills",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "C00019",
    "description": "Certificate in Liasion Interpreting",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0004",
    "description": "National Certificate in Adult Literacy and Numeracy Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0012",
    "description": "National Certificate in Air Traffic Services",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0023",
    "description": "National Certificate in Aviation",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0031",
    "description": "National Certificate in Border Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0032",
    "description": "National Certificate in Bovine Leather Manufacture",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0037",
    "description": "National Certificate in Career Practice",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0051",
    "description": "National Certificate in Clothing Manufacture",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0053",
    "description": "National Certificate in Community Recreation",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0054",
    "description": "National Certificate in Community Support Services",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0057",
    "description": "National Certificate in Composites",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0058",
    "description": "National Certificate in Compost",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0059",
    "description": "National Certificate in Concrete",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0071",
    "description": "National Certificate in Deer Farming",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0077",
    "description": "National Certificate in Drainlaying",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0089",
    "description": "National Certificate in Electronics Technology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0092",
    "description": "National Certificate in Equine",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0095",
    "description": "National Certificate in Fellmongery",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0098",
    "description": "National Certificate in Fibreboard ",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0103",
    "description": "National Certificate in Fitness",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0109",
    "description": "National Certificate in Footwear",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0126",
    "description": "National Certificate in Health, Disability, and Aged Support ",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0127",
    "description": "National Certificate in Heating Ventilating and Air Conditioning",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0129",
    "description": "National Certificate in Horse Trek Guiding",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0139",
    "description": "National Certificate in Intermediate Scaffolding",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0142",
    "description": "National Certificate in Joinery",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0147",
    "description": "National Certificate in Manufacturing",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0152",
    "description": "National Certificate in Marina Operations",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0160",
    "description": "National Certificate in Mechanical Engineering",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0161",
    "description": "National Certificate in Mediation",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0163",
    "description": "National Certificate in Metal Casting",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0187",
    "description": "National Certificate in Print Industry Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0189",
    "description": "National Certificate in Project Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0191",
    "description": "National Certificate in Public Sector Compliance",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0192",
    "description": "National Certificate in Public Sector Services",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0201",
    "description": "National Certificate in Refrigeration and Air Conditioning",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0207",
    "description": "National Certificate in Road Marking",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0211",
    "description": "National Certificate in Sales",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0212",
    "description": "National Certificate in Scaffolding",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EDCEK",
    "description": "Doctor of Laws",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EGCEZ",
    "description": "Graduate Certificate in Resource Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EGDFE",
    "description": "Graduate Diploma in Applied Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EGDFX",
    "description": "Graduate Diploma in Tourism Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMGE",
    "description": "Master of Architecture",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMGH",
    "description": "Master of Aviation",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMGN",
    "description": "Master of Communication Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMGP",
    "description": "Master of Conservation",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPCKW",
    "description": "Postgraduate Certificate in Physiotherapy",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPCKY",
    "description": "Postgraduate Certificate in Rehabilitation",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPDLG",
    "description": "Postgraduate Diploma (any other field)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MMKT-DG",
    "description": "Master of Marketing",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0244",
    "description": "National Certificate in Veterinary Nursing",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0257",
    "description": "National Certificate in Zero Waste and Resource Recovery",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0007",
    "description": "National Diploma in Airport Operations",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0015",
    "description": "National Diploma in Casino Gaming",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0016",
    "description": "National Diploma in Civil Engineering",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0022",
    "description": "National Diploma in Counselling",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0023",
    "description": "National Diploma in Diving",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0033",
    "description": "National Diploma in Forestry",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0034",
    "description": "National Diploma in Funeral Directing",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0042",
    "description": "National Diploma in Iwi/Maori Social Services",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0048",
    "description": "National Diploma in Mental Health",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0049",
    "description": "National Diploma in Nga Mahi a te Whare Pora",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0053",
    "description": "National Diploma in Print Industry Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0059",
    "description": "National Diploma in Public Sector Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0061",
    "description": "National Diploma in Resource Recovery",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0063",
    "description": "National Diploma in Seafood Maori",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MAFIN-DG",
    "description": "Master of Applied Finance",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDDS-DG",
    "description": "Postgraduate Diploma in Dance Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BUP-DG",
    "description": "Bachelor of Urban Planning",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BBMH-DG",
    "description": "Bachelor of Biomedical Science (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MEG-DG",
    "description": "Master of Engineering Geology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BEDHM-DG",
    "description": "Bachelor of Education (Teaching)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BDES-DG",
    "description": "Bachelor of Design",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BEDE-DG",
    "description": "Bachelor of Education (Teaching)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BEDRP-DG",
    "description": "Bachelor of Education (Teaching)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BEDEJ-DG",
    "description": "BEd(Tchg) Joint",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "NAUT",
    "description": "AUT International Foundation Certificate",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "OPTS",
    "description": "Open Polytech Cert in Tertiary Stdy Skills",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "FSCT",
    "description": "Foundation Science",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "UOCP",
    "description": "Certificate in Universitity Preparation",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "WITT",
    "description": "Western Inst of Tech Cert in Tertiary Studies",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "FSCV",
    "description": "Foundation Science",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "SSLD",
    "description": "Matura / Secondary School-Leaving Diploma",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "STN1",
    "description": "Steiner School Certificate Level 1",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "CII3",
    "description": "Cambridge International Examinations - A Levels",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "TNSC",
    "description": "Tonga National Form Seven Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "ASMC",
    "description": "Atestat de Studii Medii de Cultura Generala",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "HLT",
    "description": "Habilitacao Litersarias",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "ESEC",
    "description": "Eritrean Secondary Education Cert Examination",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "STN3",
    "description": "Steiner School Certificate Level 3",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "HSDS",
    "description": "High School Diploma (Natl Uni of Singapore)",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "UAI",
    "description": "Universities Admission Index (UAI)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "IM",
    "description": "Diploma in Educ of Deaf",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "JB",
    "description": "Diploma in Industrial Design",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "JC",
    "description": "Associate of Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "JO",
    "description": "Certificate in the Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "JT",
    "description": "Diploma in Fine Art",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "JV",
    "description": "Certificate in Clothing and Textiles",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "JY",
    "description": "NZ Institute of Management Mgt Diploma",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "KI",
    "description": "Diploma in Sheep Farming",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "KN",
    "description": "Diploma in Training and Development",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "KP",
    "description": "Diploma in Theological Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MA",
    "description": "Bachelor of Commerce (Agricultural)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MC",
    "description": "Bachelor of Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MP",
    "description": "Bachelor of Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MQ",
    "description": "Bachelor of Parks and Recreation Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MW",
    "description": "Bachelor of Veterinary Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND",
    "description": "Bachelor of Divinity",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NP",
    "description": "Diploma in Dairy Science & Technology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "OC",
    "description": "Fine Arts Intermediate",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "OH",
    "description": "Forestry Intermediate",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "OM",
    "description": "Veterinary Intermediate",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "OW",
    "description": "Town Planning Intermediate",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PK",
    "description": "Master of Agricultural Economics",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "QO",
    "description": "Bachelor of Technology (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "QQ",
    "description": "Bachelor of Arts/Bachelor of Management St",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "RM",
    "description": "Bachelor of Horticultural Science Honours",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "RO",
    "description": "Bachelor of Medical Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "RS",
    "description": "Bachelor of Planning",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "SL",
    "description": "Diploma in Psychology (Clinical)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "SQ",
    "description": "Diploma in Educational Psychology (PG)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "FN",
    "description": "Advanced Technical Teacher's Certificate",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GK",
    "description": "B. of Commerce Honours/Bachelor of Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "YT",
    "description": "Doctor of Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ZF",
    "description": "Certificate in General Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ZG",
    "description": "Certificate in Humanities",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ZH",
    "description": "Certificate in Industrial Relations",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BPA-DG",
    "description": "Bachelor of Property Administration",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "CPG",
    "description": "Certificato Primero Grau",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "CDE",
    "description": "Certificat de Maturite",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "MAT",
    "description": "Matura",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "HSDI",
    "description": "High School Diploma (international USA school)",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "CEM",
    "description": "Certificado de Ensino Medio",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "MAC",
    "description": "Maturita Cantonale",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "SSC",
    "description": "Senior School Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "TES",
    "description": "TES",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "LCA",
    "description": "Leaving Certificate Applied",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "IGSE",
    "description": "Cambridge International Examinations - IGSE",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "HSDS",
    "description": "High School Diploma from a Specialised High School",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "NSC",
    "description": "Senior/National Senior Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "MCE",
    "description": "Matriculation Certificate Exam",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "HGCS",
    "description": "Higher Intl. General Cert. of Secondary Education",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "UACE",
    "description": "Uganda Advanced Certificate of Education",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "GED",
    "description": "General Educational Development Test",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "UPUS",
    "description": "University Preparatory Year (NUS)",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "CBE",
    "description": "Certificate of Basic Education",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "CPE",
    "description": "Certificate of Prepatory Education",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "BAG",
    "description": "Begrut (Matriculation) or Mechina",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "B2P",
    "description": "Baccalauréat 2ème partie",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "AISC",
    "description": "All India Senior School Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "BPC",
    "description": "Brevet d'Etudes du Premier Cycle",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "GCGS",
    "description": "Certificat de Absolvire a Liceului",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "CUSE",
    "description": "Certificate of Unified State Examination",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "SCGN",
    "description": "Singapore/Cambridge GCE N levels",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "CANS",
    "description": "Completion Grade 12 Standing/Division IV Standing",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "SQCA",
    "description": "Scottish Qulification Certificate Advanced Highers",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "AUS2",
    "description": "Higher School Certificate (HSC)",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "AUS4",
    "description": "Queensland Senior Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "DPSE",
    "description": "Diploma of Completed Secondary Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPCMP-DP",
    "description": "Diploma in Computational Mathematics",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EXDPM-DP",
    "description": "Executant Diploma in Music",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BECJT-DG",
    "description": "Bachelor of Engineering (Conjoint)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BEDTE-DG",
    "description": "Bachelor of Education (Teaching)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BMUS-DG",
    "description": "Bachelor of Music",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "LLBCJ-DG",
    "description": "Bachelor of Laws (Conjoint)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MPH-DG",
    "description": "Master of Public Health",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "LD",
    "description": "NZ Certificate in Building",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "LM",
    "description": "NZ Certificate in Fire Technology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "LR",
    "description": "NZ Certificate in Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "LV",
    "description": "NZ Certificate in Town Planning",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "LW",
    "description": "NZ Certificate in Town and Country Planning Draughting",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "CRTFS-CT",
    "description": "Certificate in Foundation Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPHAR-DG",
    "description": "Doctor of Pharmacy",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GDCOU-DP",
    "description": "Graduate Diploma in Counselling",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GDTSE-DP",
    "description": "Graduate Diploma in Teaching (Secondary)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDTS-DP",
    "description": "Postgraduate Diploma in Translation Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDAM-DP",
    "description": "Postgraduate Diploma in Arts Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDGT-DP",
    "description": "Postgraduate Diploma in Geothermal Energy Technology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "HDPTC-DP",
    "description": "Higher Diploma of Teaching",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GDTP-DP",
    "description": "Graduate Diploma of Teaching (Primary)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BPE-DG",
    "description": "Bachelor of Physical Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "FCED-CT",
    "description": "Foundation Certificate Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPDD-DP",
    "description": "Diploma of Dance and Drama in Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "CTSSS-CT",
    "description": "Certificate in Support Services in Schools",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DCLPS-DG",
    "description": "Doctor of Clinical Psychology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPVLH-DP",
    "description": "Diploma in Valuation (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "TRANC-CT",
    "description": "Transitional Certificate",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "COPMU-CT",
    "description": "Certificate of Proficiency for Massey University",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BVA-DG",
    "description": "Bachelor of Visual Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MSLTP-DG",
    "description": "Master of Speech Language Therapy Practice",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDFA-DP",
    "description": "Postgraduate Diploma in Fine Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDAP-DP",
    "description": "Postgraduate Diploma in Applied Psychology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDFS-DP",
    "description": "Postgraduate Diploma in Forensic Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPEDP-DP",
    "description": "Diploma in Educational Psychology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "UI",
    "description": "Master of Commerce and Administration",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "UM",
    "description": "Master of Guidance and Counselling",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "VC",
    "description": "Medicine Professionals",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "WG",
    "description": "Postgraduate Diploma in Parks, Recreation and Tourism",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "WI",
    "description": "Doctor of Dental Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "WY",
    "description": "TeTimatanga Hou Programme",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "WZ",
    "description": "Transitional Certificate",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BHSH-DG",
    "description": "Bachelor of Health Sciences (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGCEN-CT",
    "description": "Postgraduate Certificate in Engineering",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MEDLD-DG",
    "description": "Master of Educational Leadership",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MCOUN-DG",
    "description": "Master of Counselling",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDCT-DP",
    "description": "Postgraduate Diploma in Counselling Theory",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDSW-DP",
    "description": "Postgraduate Diploma in Social Work",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "XD",
    "description": "Canterbury Passes",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "TY",
    "description": "Engineering Intermediate/Bachelor of Law",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "KB",
    "description": "Meat Inspectors Certificate",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "TX",
    "description": "Medicine Intermediate/Pharmacy Intermediat",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "WT",
    "description": "Personal Interest",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "CP",
    "description": "University of the South Pacific Foundation Year",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMHH",
    "description": "Master of General Practice",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMHK",
    "description": "Master of Information Sciences",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMHV",
    "description": "Master of Medical Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMHX",
    "description": "Master of Ministry",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMHZ",
    "description": "Master of Ophthalmology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMIC",
    "description": "Master of Planning",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EDEC",
    "description": "Diploma in Veterinary Nursing",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EGCES",
    "description": "Graduate Certificate in Commerce",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMGM",
    "description": "Master of Commerce and Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MPHC",
    "description": "Master of Primary Health Care",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ENCJB",
    "description": "National Certificate in Early Childhood Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ENCJP",
    "description": "National Diploma in Design",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EBAM",
    "description": "Bachelor of Business",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EBAU",
    "description": "Bachelor of Computing Systems",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EBBA",
    "description": "Bachelor of Dance",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EBBI",
    "description": "Bachelor of Environmental Planning",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EBBL",
    "description": "Bachelor of Information Sciences",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EBBU",
    "description": "Bachelor of Mathematical Sciences",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EBBZ",
    "description": "Bachelor of Nursing",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EBCJ",
    "description": "Bachelor of Teaching",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EHCK",
    "description": "Bachelor of Teaching (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EHCN",
    "description": "Bachelor of Tourism (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EBCQ",
    "description": "Bachelor of Viticulture and Oenology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ECDE",
    "description": "Certificate of Higher Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EBAG",
    "description": "Bachelor of Applied Social Work",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ECCX",
    "description": "Certificate in Health Promotion",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EDDK",
    "description": "Diploma in Business Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGCEC-CT",
    "description": "Postgraduate Certificate in Commercialisation and Entrepreneurship",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EAAA",
    "description": "Associate of Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EHAK",
    "description": "Bachelor of Aviation Management (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EGCEX",
    "description": "Graduate Certificate in Mathematical Sciences",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMGQ",
    "description": "Master of Construction Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMHA",
    "description": "Master of Engineering Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDPA-DP",
    "description": "Postgraduate Diploma in Creative and Performing Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "CTCOL-CT",
    "description": "Certificate in Computer Literacy",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "CTLSS-CT",
    "description": "Certificate in Library Services in Schools",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDLI-DP",
    "description": "Postgraduate Diploma in Early Literacy Intervention",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPHUS-DP",
    "description": "Diploma in Human Services",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGCAP-CT",
    "description": "Postgraduate Certificate in Academic Practice",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "HCTAE-CT",
    "description": "Higher Certificate in Art Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "HCTTS-CT",
    "description": "Higher Certificate in TESSOL",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPTPD-DP",
    "description": "Diploma in Teaching People with Disabilities",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDSS-DP",
    "description": "Postgraduate Diploma in Social Science Research Methods",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BNURC-DG",
    "description": "Bachelor of Nursing (Conjoint)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MHP",
    "description": "Master of Health Psychology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EDKC",
    "description": "NZ Diploma in Business",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPCKP",
    "description": "Postgraduate Certificate in Engineering",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPDMZ",
    "description": "Postgraduate Diploma in Tourism",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPDNB",
    "description": "Postgraduate Diploma of International Hospitality Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPDNF",
    "description": "Postgraduate Diploma in Marketing",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPCKR",
    "description": "Postgraduate Certificate in Higher Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPDLH",
    "description": "Postgraduate Diploma in Agricultural Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPDLK",
    "description": "Postgraduate Diploma in Business",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPDMD",
    "description": "Postgraduate Diploma in Medical Technology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPCKS",
    "description": "Postgraduate Certificate in Medical Technology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPDMB",
    "description": "Postgraduate Diploma in Maori and Indigenous Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EDCOKH",
    "description": "Overseas Doctoral Degree",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPDLO",
    "description": "Postgraduate Diploma in Dentistry",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0215",
    "description": "National Certificate in Seafood Maori",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0222",
    "description": "National Certificate in Social Services",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0235",
    "description": "National Certificate in Tamaki Ora - Well Child Services",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "C00006",
    "description": "Certificate in Business (Introductory)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "C00013",
    "description": "Certificate in Employment Skills",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "D00003",
    "description": "Diploma in Design Media",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "D00006",
    "description": "Diploma in Information and Library Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0008",
    "description": "National Certificate in Aeronautical Engineering",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0019",
    "description": "National Certificate in Apiculture",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0035",
    "description": "National Certificate in Business Administration and Computing",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0039",
    "description": "National Certificate in Carpet Manufacture",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0040",
    "description": "National Certificate in Casino",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0041",
    "description": "National Certificate in Caving",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0060",
    "description": "National Certificate in Conservation",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0062",
    "description": "National Certificate in Contact Centre Operations",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0064",
    "description": "National Certificate in Corrugated Case Converting",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0065",
    "description": "National Certificate in Crane Operations",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0073",
    "description": "National Certificate in Design",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0079",
    "description": "National Certificate in Drilling",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0081",
    "description": "National Certificate in Dry Cleaning",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0082",
    "description": "National Certificate in Electrical Apparatus in Explosive Atmospheres",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0088",
    "description": "National Certificate in Electronics Manufacturing",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0096",
    "description": "National Certificate in Fencing",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0113",
    "description": "National Certificate in Frame and Truss Manufacture",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0118",
    "description": "National Certificate in Glass Container Manufacturing",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0124",
    "description": "National Certificate in Hauora (Maori Health)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0125",
    "description": "National Certificate in Hazardous Waste",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0146",
    "description": "National Certificate in Locksmithing",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0150",
    "description": "National Certificate in Maori Environmental Processing",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0153",
    "description": "National Certificate in Marine Sales and Services",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0157",
    "description": "National Certificate in Mathematics",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0159",
    "description": "National Certificate in Meat Retail",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0165",
    "description": "National Certificate in Museum Practice",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0174",
    "description": "National Certificate in Passenger Service",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0178",
    "description": "National Certificate in Pharmaceutical and Allied Products Manufacturing",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0186",
    "description": "National Certificate in Precast Concrete",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EDDS",
    "description": "Diploma in Language",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EDDU",
    "description": "Diploma in Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EDDZ",
    "description": "Diploma in Software",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EGCEV",
    "description": "Graduate Certificate in Landscape Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EGDFH",
    "description": "Graduate Diploma in Business Administration",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EGDFS",
    "description": "Graduate Diploma in Professional Accounting",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EGDFU",
    "description": "Graduate Diploma in Social Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMGD",
    "description": "Master of Applied Finance",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0239",
    "description": "National Certificate in TESOL",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0241",
    "description": "National Certificate in Timber Machining",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0245",
    "description": "National Certificate in Wahi Tapu",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0249",
    "description": "National Certificate in Welding",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0252",
    "description": "National Certificate in Wood Fibre Manufacturing",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0256",
    "description": "National Certificate in Youth Work",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0017",
    "description": "National Diploma in Community Recreation",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0024",
    "description": "National Diploma in Drilling",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0030",
    "description": "National Diploma in Equine (Farriery)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0035",
    "description": "National Diploma in Hauora (Maori Health)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0037",
    "description": "National Diploma in Hearing Therapy",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0041",
    "description": "National Diploma in Intelligence Analysis",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0055",
    "description": "National Diploma in Project Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0065",
    "description": "National Diploma in Security",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0068",
    "description": "National Diploma in Social Services",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0078",
    "description": "National Diploma in Whakairo",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MDISM-DG",
    "description": "Master of Disaster Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGCDM-CT",
    "description": "Postgraduate Certificate in Disaster Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDDS-DP",
    "description": "Postgraduate Diploma in Dance Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGCSL-DG",
    "description": "Postgraduate Certificate in Social and Community Leadership",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MUPUD-DG",
    "description": "Master of Urban Planning (Professional) and Urban Design",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MCD-DG",
    "description": "Master of Community Dance",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGCMI-CT",
    "description": "Postgraduate Certificate in Māori and Indigenous Business",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "FSTCT-CT",
    "description": "Foundation Studies Certificate",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BMEDT-DG",
    "description": "BMusEd/Dip(Tchg)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGCPM-CT",
    "description": "Postgraduate Certificate in Engineering Project Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BDESC-DG",
    "description": "Bachelor of Design (Conjoint)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDHI-DP",
    "description": "Postgraduate Diploma in Higher Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGCHE-CT",
    "description": "Postgraduate Certificate in Heritage Conservation",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MEDJB-DG",
    "description": "Master of Education",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "FSSM",
    "description": "Foundation Social Science",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "SFSS",
    "description": "Foundation Social Science",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "UWAF",
    "description": "University of Western Australia Fdn Programme ",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "FSCM",
    "description": "Foundation Science",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "LUCP",
    "description": "Certificate in Preparation Studies",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "AIFC",
    "description": "International Foundation Certificate",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "MUFS",
    "description": "Certificate in Foundation Studies",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "KFSS",
    "description": "Foundation Social Science",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "FSSN",
    "description": "Foundation Social Science",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "UOFS",
    "description": "Foundation Studies Certificate",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "UQFY",
    "description": "University of Queensland Cert IV in Uni Prep",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "UOWA",
    "description": "Certificate of Attainment in Foundation Studies",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "UAFS",
    "description": "University of Auckland Cert in Foundation Studies",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "SGEO",
    "description": "Sri Lankan General Cert of Education (Ordinary)",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "ESC",
    "description": "Ensino Secundario Complementar",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "CCSE",
    "description": "Certificate of Completed Secondary Education",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "SCN",
    "description": "Senior Cert. (no Matriculation Endorsement)",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "SHSC",
    "description": "Senior High School Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "AOL",
    "description": "Apolytirio of Lykeio",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "CIE2",
    "description": "Cambridge International Examinations - AS Level",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "ICCE",
    "description": "ICCE Advanced (Academic) Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "MCT",
    "description": "Matura Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "GSC",
    "description": "General Secondary Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "GCEA",
    "description": "General Certificate of Education (Advanced Level)",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "SQCH",
    "description": "Scottish Qulification Certificate Highers",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "COM",
    "description": "Certificate of Maturity",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "ELC",
    "description": "Established Leaving Certificate",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GY",
    "description": "Bachelor of Science/Bachelor of Commerce",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "HC",
    "description": "Registered General Nurse",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "HE",
    "description": "Diploma Horticultural Fruit",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "HG",
    "description": "Diploma in Horticultural Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "HI",
    "description": "Advanced Diploma in Nursing",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "HM",
    "description": "Diploma of Hamilton Teachers College",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "HQ",
    "description": "Diploma in Journalism",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "HZ",
    "description": "Diploma in Orchard Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "IF",
    "description": "Diploma in Occupation Safety and Health",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "IO",
    "description": "Diploma in Physical Education/Bach Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "IR",
    "description": "Certificate for Orthopaedic Technicians",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "JG",
    "description": "Technicians Certificate (Electrical)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "JM",
    "description": "National Certificate in Horticulture",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "JR",
    "description": "Higher Certificate in Data Processing",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "KE",
    "description": "Postgraduate Certificate in Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MN",
    "description": "Bachelor of Technology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MU",
    "description": "Bachelor of Social Sciences",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NZ",
    "description": "Diploma in Landscape Technology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "OF",
    "description": "Engineering Intermediate",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PJ",
    "description": "Master of Business Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PR",
    "description": "Master of Environmental Design",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "QM",
    "description": "Bachelor of Laws (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "QR",
    "description": "Bach. of Science/Bach. of Management St",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "QS",
    "description": "Bach. Social Sci/Bach. Management Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "RF",
    "description": "Bachelor of Building Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "RQ",
    "description": "Bachelor of Speech and Language Therapy",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "SB",
    "description": "Diploma in Business Administration",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "SP",
    "description": "Diploma in Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPBUS-DP",
    "description": "Diploma in Business",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "FH",
    "description": "NZ Kindergarten Teacher's Diploma",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GN",
    "description": "Bachelor of Commerce/Bachelor of Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ZC",
    "description": "Certificate in Comm Work Practical",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ZL",
    "description": "Certificate in Maori Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ZQ",
    "description": "Certificate in Fitness Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BCOMH-DG",
    "description": "Bachelor of Commerce (Honours)",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "AU12",
    "description": "Queensland Certificate of Education",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "PSSC",
    "description": "Pacific Senior Secondary Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "VWO",
    "description": "Voorbereidend Wetenschappelijk Onderwijs",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "AM",
    "description": "Advanded Matriculation",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "WAEC",
    "description": "Senior School Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "NCE1",
    "description": "NCEA Level 1",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "BFEM",
    "description": "Brevet de Fin d'Etudes Moyennes",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "APL",
    "description": "Apolytirio of Lykeio",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "MSCE",
    "description": "Malawi School Cerificate of Education",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "CANY",
    "description": "Senior Secondary Graduation Diploma (Yukon)",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "DEUG",
    "description": "Diplome d'Etudes Universitaires Generales",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "NSCH",
    "description": "Namibia Senior Secondary Cert. (Higher Level)",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "DSO",
    "description": "Diploma van Secundair Onderwijs",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "CAN4",
    "description": "High School Graduation Diploma (New Brunswick)",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "ADM",
    "description": "Attestato di Maturita",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "BDE",
    "description": "Baccalaureat de l'Enseignement du Second Degre",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "DDM",
    "description": "Diploma di Maturita",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "GLO",
    "description": "Gumnaasiumi Ioputunnistus (Secondary School Cert)",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "SCM",
    "description": "Senior Cert. (with Matriculation Endorsement)",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "CTC",
    "description": "Certificat du Tronc Commun",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "NZIM",
    "description": "NZ Institute of Mgmt Certificate in Management",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "CAN9",
    "description": "High School Graduation Dip (Prince Edward Island)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BTP-DG",
    "description": "Bachelor of Town Planning",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPCRI-DP",
    "description": "Diploma in Criminology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPTP-DP",
    "description": "Diploma in Town Planning",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPVAL-DP",
    "description": "Diploma in Valuation",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BSCC-DG",
    "description": "Bachelor of Science (Conjoint)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BFA-DG",
    "description": "Bachelor of Fine Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BTHEO-DG",
    "description": "Bachelor of Theology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BSC-DG",
    "description": "Bachelor of Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BSCHC-DG",
    "description": "Bachelor of Science (Honours) (Conjoint)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MTAXS-DG",
    "description": "Master of Taxation Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MPROP-DG",
    "description": "Master of Property",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "LLMEN-DG",
    "description": "Master of Laws in Environmental Law",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "LB",
    "description": "NZ Certificate in Civil Draughting",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "LF",
    "description": "NZ Certificate in Commerce",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "LS",
    "description": "NZ Certificate in Statistics",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "LXA",
    "description": "National Certificate in Business",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "LY",
    "description": "National Diploma of Medical Diagnostic Imaging",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GDART-DP",
    "description": "Graduate Diploma in Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGCHM-CT",
    "description": "PG Certificate in Health (Mental Health Nursing)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BAHED-DG",
    "description": "Bachelor of Adult and Higher Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BHSCC-DG",
    "description": "Bachelor of Health Sciences (Conjoint)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BHSC-DG",
    "description": "Bachelor of Health Sciences",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BPERF-DG",
    "description": "Bachelor of Performing Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BEDTL-DG",
    "description": "Bachelor of Education (TESOL)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GDENG-DP",
    "description": "Graduate Diploma in Engineering",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGCMH-CT",
    "description": "Postgraduate Certificate in Maori Health",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPSW-DP",
    "description": "Diploma in Social Work",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BSS-DG",
    "description": "Bachelor of Social Sciences (Human Services)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BTHEH-DG",
    "description": "Bachelor of Theology (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "LITTD-DG",
    "description": "Doctor of Literature",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GDTHE-DP",
    "description": "Graduate Diploma in Theology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDCO-DP",
    "description": "Postgraduate Diploma in Commerce",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDLT-DP",
    "description": "Postgraduate Diploma in Language Teaching",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "COPVU-CT",
    "description": "Certificate of Proficiency for Victoria University",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "COPOU-CT",
    "description": "Certificate of Proficiency for Otago University",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MNURS-DG",
    "description": "Master of Nursing",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MPSTD-DG",
    "description": "Master of Professional Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPTPR-DP",
    "description": "Diploma in Teaching (Primary)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPCNT-DP",
    "description": "Diploma in Counselling Theory",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPLST-DP",
    "description": "Diploma in Labour Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPSTA-DP",
    "description": "Diploma in Statistics",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "TT",
    "description": "Bachelor of Law/Law Professionals",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "TZ",
    "description": "Master of Music/ Bachelor of Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "UC",
    "description": "Master of Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "UK",
    "description": "Master of Consumer and Applied Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "UN",
    "description": "Master of Literature",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "UR",
    "description": "Master of Physical Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "WL",
    "description": "Foreign Language Reading Examination",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MLS-DG",
    "description": "Master of Legal Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "XA",
    "description": ":",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "TH",
    "description": "Law Professionals/Cert of Proficiency",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "TG",
    "description": "Law Professional/Bachelor of Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "WN",
    "description": "Health Science 1",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMHR",
    "description": "Master of Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMHT",
    "description": "Master of Maori and Pacific Development",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMIM",
    "description": "Master of Speech and Language Therapy",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ECIW",
    "description": "MIT Certificate in Pre-Degree Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MCLPH-DG",
    "description": "Master of Clinical Pharmacy",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EHBM",
    "description": "Bachelor of Information Sciences (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EBCH",
    "description": "Bachelor of Sport",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EGDFF",
    "description": "Graduate Diploma in Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMHW",
    "description": "Master of Midwifery",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMIO",
    "description": "Master of Surveying",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ENCIZ",
    "description": "National Certificate in Automotive Reglazing",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ENCJJ",
    "description": "National Certificate in Tourism",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ENCJO",
    "description": "National Diploma in Construction Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ENCJQ",
    "description": "National Diploma in Fitness",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ENCJW",
    "description": "National Diploma in Social Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ENCJY",
    "description": "National Diploma in Reo Mäori",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ENCJZ",
    "description": "National Diploma in Teaching",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EBAJ",
    "description": "Bachelor of Aviation Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EHAL",
    "description": "Bachelor of Biomedical Science (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EBCC",
    "description": "Bachelor of Physiotherapy",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ECCZ",
    "description": "Certificate in Science and Technology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EDDG",
    "description": "Diploma (any other diploma)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EDDJ",
    "description": "Diploma in Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EDDP",
    "description": "Diploma in Information Technology Support",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "FV",
    "description": "Diploma in Fine Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "LE",
    "description": "Certificate in Social Work",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BPEHN-DG",
    "description": "Bachelor Physical Education (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EGDFQ",
    "description": "Graduate Diploma in Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EGDFZ",
    "description": "Graduate Diploma in Viticulture and Oenology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMGV",
    "description": "Master of Defense Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDEM-DP",
    "description": "Postgraduate Diploma of Education (Music)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDMG-DP",
    "description": "Postgraduate Diploma in Educational Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GDFA-DP",
    "description": "Graduate Diploma in Fine Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGCCE-CT",
    "description": "Postgraduate Certificate in Clinical Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BHUMS-DG",
    "description": "Bachelor of Human Services",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGCPS-CT",
    "description": "Postgraduate Certificate in Professional Supervision",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPDLV",
    "description": "Postgraduate Diploma in Health Sciences",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPDLX",
    "description": "Postgraduate Diploma in Horticultural Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPDLY",
    "description": "Postgraduate Diploma in Industrial Health",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDMT-DP",
    "description": "Postgraduate Diploma in Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPDML",
    "description": "Postgraduate Diploma in Psychology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPCKZ",
    "description": "Postgraduate Certificate in Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ENDKB",
    "description": "National Diploma in Youth Work",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EHOKG",
    "description": "Overseas Bachelors Degree (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPDNE",
    "description": "Postgraduate Diploma in Human Resources",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPDLN",
    "description": "Postgraduate Diploma in Computer and Infor Sciences",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MPACC-DG",
    "description": "Master of Professional Accounting",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPDMP",
    "description": "Postgraduate Diploma in Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPDMT",
    "description": "Postgraduate Diploma in Specialist Teaching",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0233",
    "description": "National Certificate in Stonemasonry",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "C00002",
    "description": "Certificate in Animal Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "C00005",
    "description": "Certificate in Automotive and Mechanical Engineering",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "C00007",
    "description": "Certificate in Business Administration and Computing",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "C00011",
    "description": "Certificate in Electrical and Electronic Engineering",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "C00014",
    "description": "Certificate in English",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "C00020",
    "description": "Certificate in Multiskill Building Construction",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0006",
    "description": "National Certificate in Adventure Tourism",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0011",
    "description": "National Certificate in Agrichemical",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0028",
    "description": "National Certificate in Biosecurity",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0033",
    "description": "National Certificate in Brick and Block Laying",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0044",
    "description": "National Certificate in Civil Construction",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0046",
    "description": "National Certificate in Civil Engineering - Asset Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0072",
    "description": "National Certificate in Demolition",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0076",
    "description": "National Certificate in Diving",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0085",
    "description": "National Certificate in Electrity Supply",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0105",
    "description": "National Certificate in Floor and Wall tiling",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0140",
    "description": "National Certificate in Irrigation",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0149",
    "description": "National Certificate in Maori Business and Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0173",
    "description": "National Certificate in Paint Manufacture",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0176",
    "description": "National Certificate in Pavement",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0182",
    "description": "National Certificate in Police Forensic Mapping",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0184",
    "description": "National Certificate in Poultry Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0188",
    "description": "National Certificate in Printing",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0193",
    "description": "National Certificate in Quality Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0196",
    "description": "National Certificate in Radio",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0199",
    "description": "National Certificate in Recreation and Sport",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0202",
    "description": "National Certificate in Renewable Energy Systems",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0206",
    "description": "National Certificate in Rigging",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EGCEU",
    "description": "Graduate Certificate in English",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EGCEW",
    "description": "Graduate Certificate in Law",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EGDFN",
    "description": "Graduate Diploma in Higher Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EGDFR",
    "description": "Graduate Diploma in Music",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EGDFT",
    "description": "Graduate Diploma in Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0238",
    "description": "National Certificate in Telecommunications",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0248",
    "description": "National Certificate in Water Treatment",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0251",
    "description": "National Certificate in Whanau and Foster Care",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0012",
    "description": "National Diploma in Building Control Surveying",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0013",
    "description": "National Diploma in Business Administration",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0020",
    "description": "National Diploma in Composting",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0025",
    "description": "National Diploma in Drinking-Water",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0026",
    "description": "National Diploma in Electrical Engineering",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0031",
    "description": "National Diploma in Extractive Industries",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0051",
    "description": "National Diploma in Plastics Processing Tech",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0057",
    "description": "National Diploma in Public Sector Compliance Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0058",
    "description": "National Diploma in Public Sector Maori",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0075",
    "description": "National Diploma in Tourism Maori",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ECDD",
    "description": "Certificate of Attainment in Foundation Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MUP-DG",
    "description": "Master of Urban Planning",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MIT-DG",
    "description": "Master of Information Technology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MSCL-DG",
    "description": "Master of Social and Community Leadership",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MNUP-DG",
    "description": "Master of Nursing Practice",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MTES-DG",
    "description": "Master of Teaching (Secondary)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BSHPE-DG",
    "description": "Bachelor of Sport, Health and Physical Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGCEQ-CT",
    "description": "Postgraduate Certificate in Earthquake Engineering",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0130",
    "description": "National Certificate in Horticulture",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGCHI-CT",
    "description": "Postgraduate Certificate in Higher Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BASH-DG",
    "description": "Bachelor of Advanced Science (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGCER-CT",
    "description": "Postgraduate Certificate in Energy",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0074",
    "description": "National Diploma in Tourism Conventions and Incentives",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPTRE-DP",
    "description": "Diploma of Teaching (Early Childhood Education)",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "SFSC",
    "description": "Foundation Science",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "FSCS",
    "description": "Foundation Science",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "TFSC",
    "description": "Foundation Science",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "WCAD",
    "description": "Whitecliffe Foundation Cert of Arts and Design",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "MITC",
    "description": "MIT Cert in Fdn Education Tert Pathways",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "FSCC",
    "description": "Foundation Science",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "FSCF",
    "description": "Foundation Science",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "UOCF",
    "description": "Certificate in Foundation Studies",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "CAN2",
    "description": "Senior Secondary Graduation Dip (British Columbia)",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "SSG",
    "description": "Secondary School Graduation",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "DIP",
    "description": "Diplome de Fin d'Etudes",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "ESLC",
    "description": "Ethiopean School Leaving Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "SSC",
    "description": "Secondary School Leaving Certificate/Matura",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "NCE2",
    "description": "NCEA Level 2",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "SCGH",
    "description": "Singapore/Cambridge GCE H Levels",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "CANB",
    "description": "High School Diploma (Saskatchewan)",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "CSE",
    "description": "Certificate of Secondary Education",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "MAW3",
    "description": "Maw 3",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "HD",
    "description": "Diploma in Horticultural Fruit Production",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "HR",
    "description": "Diploma in Maoritanga",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "HU",
    "description": "Diploma in Musculoskeletal Medicine",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "IU",
    "description": "Diploma in Police Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "IV",
    "description": "Diploma in Parks and Recreation Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "JJ",
    "description": "Diploma in Visual Communication Design",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "KF",
    "description": "Certificate in Maoritanga",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "KG",
    "description": "Certificate in Early Childhood Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "KU",
    "description": "Diploma in Travel and Tourism",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MB",
    "description": "Bachelor of Social Work",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ML",
    "description": "Bachelor of Regional Planning",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NE",
    "description": "Bachelor of Laws",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NL",
    "description": "Diploma in Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NR",
    "description": "Diploma in Landscape Architecture",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NV",
    "description": "Diploma in Educational Psychology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "OL",
    "description": "Law Intermediate",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "QA",
    "description": "Bachelor of Commerce (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "QE",
    "description": "Bachelor of Pharmacy Honours",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "RP",
    "description": "Bachelor of Mineral Technology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "SO",
    "description": "Diploma in Commerce",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "TJ",
    "description": "Bachelor of Law Honours/Bach. of Commerce",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DARTA-DP",
    "description": "Diploma in Arts Administration",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPMGT-DP",
    "description": "Diploma in Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "FA",
    "description": "Diploma in Accounting",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "FK",
    "description": "Higher Technical Teacher's Certificate",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "FM",
    "description": "Diploma in Banking Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "FX",
    "description": "Diploma in Food Quality Assurance",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "FY",
    "description": "Diploma in Financial Mathematics",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GD",
    "description": "Bachelor of Arts/Bachelor of Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GG",
    "description": "Bachelor of Arts/Diploma in Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ZE",
    "description": "Certificate of Proficiency in English Language",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "CCEF",
    "description": "Certificado de Conclusão de Ensino Fundamental",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "CECS",
    "description": "Caribbean Exam Council Secondary Education Cert",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "UEC",
    "description": "Unified Examination Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "HSGD",
    "description": "High School Graduation Diploma",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "SUG",
    "description": "Swiadectwo Ukonczenia Gimnazjum",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "GHO",
    "description": "Getuigschrift van Hoger Secundair Onderwijs",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "ZAH",
    "description": "Zeugnis der Allgemeinen Hochschulreife / Abitur",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "USC",
    "description": "Upper Secondary School Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "HSDI",
    "description": "High School Diploma -International American School",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "IHSC",
    "description": "Intermediate Cert. /Higher Secondary Cert.",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "DDB",
    "description": "Diploma de Bachiller Científico-Humanístico",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "KDFN",
    "description": "Kolej Danabsara Utama, Petaling Jaya - Foundation",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "GCSE",
    "description": "General Certificate of Secondary Education",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "HSII",
    "description": "High School Diploma (Iran)",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "ISC",
    "description": "Indian School Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "RSSC",
    "description": "Religious Secondary School Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "SCG",
    "description": "Diploma di Scuola Cultura Generale",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "BAC",
    "description": "Bachillerato",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "SSA",
    "description": "Secondary School Leaving Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "HASC",
    "description": "Hong Kong Advanced Supplementary Level Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "JSC",
    "description": "Junior Secondary Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "SISC",
    "description": "Solomon Islands School Certificate",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "LLBH-DG",
    "description": "Bachelor of Laws (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MA-DG",
    "description": "Master of Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MCOM-DG",
    "description": "Master of Commerce",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ME-DG",
    "description": "Master of Engineering",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BED-DG",
    "description": "Bachelor of Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BMUSP-DG",
    "description": "Bachelor of Music (Performance)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BMUSH-DG",
    "description": "Bachelor of Music (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MPA-DG",
    "description": "Master of Property Administration",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MFA-DG",
    "description": "Master of Fine Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MTHEO-DG",
    "description": "Master of Theology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MJUR-DG",
    "description": "Master of Jurisprudence",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MAUD-DG",
    "description": "Master of Audiology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DSC-DG",
    "description": "Doctor of Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "LC",
    "description": "NZ Certificate in Architectural Draughting",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "LK",
    "description": "NZ Certificate in Data Processing",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "LL",
    "description": "NZ Certificate in Engineering",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "LQ",
    "description": "NZ Certificate in Quantity Surveying",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BFACJ-DG",
    "description": "Bachelor of Fine Arts (Conjoint)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "CTRAD-CT",
    "description": "Certificate in Radiochemistry",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "CJ",
    "description": "Victoria Certificate of Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DMA-DG",
    "description": "Doctor of Musical Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGCGT-CT",
    "description": "Postgraduate Certificate in Geothermal Energy Technology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "CTEDS-CT",
    "description": "Certificate in Educational Support (Disability Studies)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPTEP-DP",
    "description": "Diploma of Teaching (Early Childhood - Pacific Islands)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GDED-DP",
    "description": "Graduate Diploma in Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "UNITC-CT",
    "description": "Unitech Certificate",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "COPWU-CT",
    "description": "Certificate of Proficiency for Waikato University",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BBIM-DG",
    "description": "Bachelor of Business and Information Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GCBUS-CT",
    "description": "Graduate Certificate in Business",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDGM-DP",
    "description": "Postgraduate Diploma in Geriatric Medicine",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPADE-DP",
    "description": "Diploma in Adult Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPCOM-DP",
    "description": "Diploma in Commerce",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPOMD-DP",
    "description": "Diploma in Occupational Medicine",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPPPT-DP",
    "description": "Diploma in Pulp and Paper Technology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPMUA-DP",
    "description": "Diploma in Music (Advanced)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPPAE-DP",
    "description": "Diploma in Paediatrics",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPENT-DP",
    "description": "Diploma in Engineering (Transportation)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "CTOCP-CT",
    "description": "Certificate in Ocular Pharmacology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "TU",
    "description": "Bachelor of Med & Surgery/Bachelor of Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "TV",
    "description": "Bachelor of Med & Surgery/Bachelor of Sci",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "UD",
    "description": "Master of Applied Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "UO",
    "description": "Master of Music",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "VA",
    "description": "Accounting Professionals",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GCLAW-CT",
    "description": "Graduate Certificate in Law",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GDETE-DP",
    "description": "Graduate Diploma in Teaching (Early Childhood Education)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GDTES-DP",
    "description": "Graduate Diploma in TESSOL",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BMSCH-DG",
    "description": "Bachelor of Medical Science (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MSW-DG",
    "description": "Master of Social Work",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "TW",
    "description": "Medicine Intermediate/Architecture Int",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMHE",
    "description": "Master of Finance",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMHL",
    "description": "Master of International Development",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMIB",
    "description": "Master of Physiotherapy",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMIK",
    "description": "Master of Social and Community Work",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMIQ",
    "description": "Master of Teaching",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMIU",
    "description": "Master of Tourism Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EDEE",
    "description": "Diploma of Higher Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EDCEN",
    "description": "Doctorate (any other field)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EGCFC",
    "description": "Graduate Certificate in TESOL",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EGDFJ",
    "description": "Graduate Diploma in Computing",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMHM",
    "description": "Master of International Hospitality Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMHN",
    "description": "Master of International Law and Politics",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMIN",
    "description": "Master of Sport and Leisure Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ENCJG",
    "description": "National Certificate in Retail",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ENCJN",
    "description": "National Diploma in Computing",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EBAO",
    "description": "Bachelor of Business Administration",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EHAP",
    "description": "Bachelor of Business Administration (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EBBF",
    "description": "Bachelor of Design Innovation",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EBBQ",
    "description": "Bachelor of Liberal Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EBCA",
    "description": "Bachelor of Oral Health",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ECCY",
    "description": "Certificate in Natural Resources",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ECDB",
    "description": "Certificate in Theological Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EBCS",
    "description": "Bachelors Degree (any other degree)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EDDO",
    "description": "Diploma in English (Advanced)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "TN",
    "description": "Bach. of Law/Bach. of Arts/Cert. Profficie",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GP",
    "description": "Bachelor of Commerce/Cert of Proficiency",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NS",
    "description": "Diploma in Teaching (o'seas)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGCCR-CT",
    "description": "Postgraduate Certificate in Clinical Pharmacy",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EBBN",
    "description": "Bachelor of International Hosptality Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMGX",
    "description": "Master of Design",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MAS-DG",
    "description": "Master of Architectural Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDBS-DP",
    "description": "Postgraduate Diploma in Building Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDEJ-DP",
    "description": "Postgraduate Diploma in Education (Jointly Badged)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DTPY-DP",
    "description": "Diploma of Teaching (Primary)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MBE-DG",
    "description": "Master of Bioscience Enterprise",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "CTITS-CT",
    "description": "Certificate in Introductory Tertiary Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "HCTTE-CT",
    "description": "Higher Certificate in Technology Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPDED-DP",
    "description": "Diploma of Drama Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPTSY-DP",
    "description": "Diploma of Teaching (Secondary)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "CTTPD-CT",
    "description": "Certificate in Teaching People with Disabilities",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "CTRUR-CT",
    "description": "Certificate in Rumaki Reo",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPCKO",
    "description": "Postgraduate Certificate in Counselling",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPDLM",
    "description": "Postgraduate Diploma in Clinical Psychology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPDMF",
    "description": "Postgraduate Diploma in Music",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPDMG",
    "description": "Postgraduate Diploma in Nursing",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPDMI",
    "description": "Postgraduate Diploma in Pharmacy",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDTF-DP",
    "description": "Postgraduate Diploma in Teaching (Secondary Field-based)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPCLA",
    "description": "Postgraduate Certificate in Social Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPCKI",
    "description": "Postgraduate Certificate in International Hospitality Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ENDKA",
    "description": "National Diploma in Tourism",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPDMV",
    "description": "Postgraduate Diploma in Sport Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPDMY",
    "description": "Postgraduate Diploma in Theology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0216",
    "description": "National Certificate in Security",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0217",
    "description": "National Certificate in Seed Dressing",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0221",
    "description": "National Certificate in Social Service Work",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0230",
    "description": "National Certificate in Sport Turf Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "C00004",
    "description": "Certificate in Applied Technology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "D00001",
    "description": "Diploma in Applied Computer Systems Engineering",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0001",
    "description": "National Certficate in Distribution",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0005",
    "description": "National Certificate in Advanced Scaffolding",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0007",
    "description": "National Certificate in Aerial Agrichemical Application",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0014",
    "description": "National Certificate in Airport Operations",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0020",
    "description": "National Certificate in Aquaculture",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0024",
    "description": "National Certificate in Baking",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0027",
    "description": "National Certificate in Bedding Manufacture",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0038",
    "description": "National Certificate in Carpentry",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0042",
    "description": "National Certificate in Character Animation",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0043",
    "description": "National Certificate in Christian Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0070",
    "description": "National Certificate in Dance",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0075",
    "description": "National Certificate in Diversional Therapy",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0078",
    "description": "National Certificate in Drama",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0084",
    "description": "National Certificate in Electrical Equipment",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0087",
    "description": "National Certificate in Electronics Engineering",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0091",
    "description": "National Certificate in Energy and Chemical Plant",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0102",
    "description": "National Certificate in Fire Detection and Alarm Systems",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0104",
    "description": "National Certificate in Fixed Fire Protection Systems",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0123",
    "description": "National Certificate in Gunsmithing",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0128",
    "description": "National Certificate in Heavy Haulage",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0143",
    "description": "National Certificate in Laundry",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0145",
    "description": "National Certificate in Lifts and Escalators",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0151",
    "description": "National Certificate in Maori Performing Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0164",
    "description": "National Certificate in Motor Industry",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0170",
    "description": "National Certificate in Outdoor Recreation",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0179",
    "description": "National Certificate in Plastics Materials",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0185",
    "description": "National Certificate in Poultry Production",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0200",
    "description": "National Certificate in Refractory Installation",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0203",
    "description": "National Certificate in Reo Maori",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EDCEH",
    "description": "Doctor of Commerce",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EGCFA",
    "description": "Graduate Certificate in Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMGK",
    "description": "Master of Business",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMGO",
    "description": "Master of Computer Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPCKV",
    "description": "Postgraduate Certificate in Pharmacy",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPDLI",
    "description": "Postgraduate Diploma in Applied Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0243",
    "description": "National Certificate in Travel",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0247",
    "description": "National Certificate in Water Reticulation",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0253",
    "description": "National Certificate in Wood Manufacturing",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0004",
    "description": "National Diploma in Aeronautical Engineering",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0014",
    "description": "National Diploma in Career Practice",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0021",
    "description": "National Diploma in Contact Centre Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0036",
    "description": "National Diploma in Hazardous Waste",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0054",
    "description": "National Diploma in Professional Practice in Design and Construction",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0067",
    "description": "National Diploma in Social Service Work",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0076",
    "description": "National Diploma in Veterinary Nursing",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0082",
    "description": "National Diploma of Business Education ",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MUPHC-DG",
    "description": "Master of Urban Planning (Professional) and Heritage Conservation",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MHRM-DG",
    "description": "Master of Human Resource Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MHC-DG",
    "description": "Master of Heritage Conservation",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDEN-DP",
    "description": "Postgraduate Diploma in Engineering",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BGS-DG",
    "description": "Bachelor of Global Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MARUP-DG",
    "description": "Master of Architecture (Professional) and Urban Planning (Professional)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "COPSP-CT",
    "description": "Certificate of Proficiency Short Programme",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDIS-DG",
    "description": "Postgraduate Diploma in Indigenous Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MEQE-DG",
    "description": "Master of Earthquake Engineering",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDTD-DP",
    "description": "Postgraduate Diploma in Therapeutic Dance",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDER-DP",
    "description": "Postgraduate Diploma in Energy",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BMECJ-DG",
    "description": "BMusEd (Joint)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0155",
    "description": "National Certificate in Maritime",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "FSSV",
    "description": "Foundation Social Science",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "FSSF",
    "description": "Foundation Social Science",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "WHPT",
    "description": "Whitireia Cert in Preparation for Tert Stds",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "FSSS",
    "description": "Foundation Social Science",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "FSCN",
    "description": "Foundation Science",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "MUCP",
    "description": "Certificate of University Preparation",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "TCFS",
    "description": "Trinity College Foundation Studies",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "UOWP",
    "description": "Certificate of University Preparation",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "CISE",
    "description": "Certificate of Incomplete Secondary Education",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "NQL",
    "description": "No Formal Secondary Qualification",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "BEPC",
    "description": "Brevet d'Etudes du Premier Cycle",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "HKDE",
    "description": "Hong Kong Diploma of Secondary Education",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "HSDP",
    "description": "High School Diploma",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "LEM",
    "description": "Licencia de Educacion Media",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "ACE",
    "description": "Accelerated Christian Education",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "HALC",
    "description": "Hong Kong Advanced Level Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "CAN8",
    "description": "Ontario Secondary School Diploma",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "MATR",
    "description": "Matriculation",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "CII2",
    "description": "Cambridge International Examinations - AS Levels",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "TB",
    "description": "Titulo de Bachiller",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "MAT",
    "description": "Maturitna Skuska / Maturita",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "CEGE",
    "description": "Certificate of General Secondary Education",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "DPIE",
    "description": "Diploma of Incompleted Secondary Education",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "GER",
    "description": "Gerchilgee (School Leaving Certificate)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "HK",
    "description": "Diploma in Nursing Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "HL",
    "description": "Diploma in Health Service Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "HX",
    "description": "Diploma in Nursery Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "II",
    "description": "Diploma in Public Health",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "IJ",
    "description": "Diploma in Theraputic Radiography",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "IN",
    "description": "Diploma in Chiropody",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "IQ",
    "description": "Diploma in Personnel Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "IT",
    "description": "Diploma in Theology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "JF",
    "description": "Diploma in Graphic Design",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "JP",
    "description": "Diploma in Computer Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "KJ",
    "description": "Diploma in Safety Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "KL",
    "description": "Diploma in Sports Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "KM",
    "description": "Diploma in Teaching",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "KQ",
    "description": "Diploma in Wildlife Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MI",
    "description": "Bachelor of Forestry Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MV",
    "description": "Bachelor of Philosophy",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MX",
    "description": "Overseas Bachelors Degree (unclassified)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MY",
    "description": "Bachelor of Commerce",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NM",
    "description": "Diploma in Social Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NN",
    "description": "Diploma in Agriculture",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "OA",
    "description": "Bachelor of Technology Intermediate",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "OI",
    "description": "Optometry Intermediate",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "OR",
    "description": "Regional Planning Intermediate",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "OT",
    "description": "Surveying Intermediate",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PE",
    "description": "Master of Laws",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PH",
    "description": "Master of Medicine",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "QG",
    "description": "Bachelor of Physical Education (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "QJ",
    "description": "Bachelor of Mineral Technology with Honours",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "QP",
    "description": "Bachelor of Law/Bach. Management Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "RA",
    "description": "Bachelor of Fine Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "RC",
    "description": "Bachelor of Agriculture",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "SC",
    "description": "Diploma in Librarianship",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "SD",
    "description": "Diploma in Regional and Resource Planning",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "SM",
    "description": "Diploma of Computational Mathematics",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "SN",
    "description": "Diploma in Psychology (Community)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "SW",
    "description": "Diploma in Optometry",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "TE",
    "description": "Bachelor of Surveying/Bachelor of Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DARTM-DP",
    "description": "Diploma in Arts Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "FL",
    "description": "Diploma in Banking",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ZB",
    "description": "Certificate in Continuing Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ZO",
    "description": "Certificate in Social Services Supervision",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "BFS",
    "description": "Bevis for Studentereksamen",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "ACG",
    "description": "ACG Certificate in Foundation Studies",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "UCE",
    "description": "Uganda Certificate of Education",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "SAT",
    "description": "SAT Reasoning Test",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "ENTR",
    "description": "Equivalent National Tertiary Entrance Rank (ENTER)",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "B1P",
    "description": "Baccalauréat 1ère partie",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "UE85",
    "description": "University Entrance Examinations prior to 1986",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "NZC",
    "description": "New Zealand School Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "CAN0",
    "description": "Diplôme d'Etudes Collègiales (Quebec)",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "EBAC",
    "description": "European Baccalaureate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "CAPE",
    "description": "Caribbean Advanced Proficiency Examination",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "OP",
    "description": "Overall Position (OP)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BSCH-DG",
    "description": "Bachelor of Science (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BSCHB-DG",
    "description": "Bachelor of Science (Human Biology)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "LLB-DG",
    "description": "Bachelor of Laws",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MSC-DG",
    "description": "Master of Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BCOMC-DG",
    "description": "Bachelor of Commerce (Conjoint)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BPRPC-DG",
    "description": "Bachelor of Property (Conjoint)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BAS-DG",
    "description": "Bachelor of Architectural Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MPLPR-DG",
    "description": "Master of Planning Practice",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MHSC-DG",
    "description": "Master of Health Sciences",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MESTU-DG",
    "description": "Master of Engineering Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "LT",
    "description": "NZ Certificate in Survey Draughting",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BEDTC-DG",
    "description": "Bachelor of Education (Teaching) (Conjoint)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BHBH-DG",
    "description": "Bachelor of Human Biology (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BPHAR-DG",
    "description": "Bachelor of Pharmacy",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPHCA-DP",
    "description": "Diploma in Health (Child Adolescent Mental Health)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPPNM-DP",
    "description": "Diploma in Politics and the News Media",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGCPP-CT",
    "description": "Postgraduate Certificate in Pharmacy Practice",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "TFC",
    "description": "Tertiary Foundation Certificate",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "FCEAP-CT",
    "description": "Foundation Certificate in English for Academic Purposes",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDBU-DP",
    "description": "Postgraduate Diploma in Business",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGCBU-CT",
    "description": "Postgraduate Certifcate in Business",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GDENT-DP",
    "description": "Graduate Diploma in Engineering (Transportation)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GDTS-DP",
    "description": "Graduate Diploma of Teaching (Secondary)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GDLED-DP",
    "description": "Graduate Diploma in Literacy Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "CTSWS-CT",
    "description": "Certificate in Social Work Supervision",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPIS-DP",
    "description": "Diploma of Information Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDTH-DP",
    "description": "Postgraduate Diploma in Theology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPESS-DP",
    "description": "Diploma of Education of Students with Special Teaching Needs",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDEP-DP",
    "description": "Postgraduate Diploma in Educational Psychology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDHS-DP",
    "description": "Postgraduate Diploma in Health Sciences",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "COPOS",
    "description": "Certificate of Proficiency for Overseas",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDPP-DP",
    "description": "Postgraduate Diploma in Pharmacy Practice",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPGET-DP",
    "description": "Diploma in Geothermal Energy Technology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPELT-DP",
    "description": "Diploma in English Language Teaching",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPENG-DP",
    "description": "Diploma in Engineering",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPCOU-DP",
    "description": "Diploma in Counselling",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPED-DP",
    "description": "Diploma in Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPLGA-DP",
    "description": "Diploma in Local Government and Administration",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPPH-DP",
    "description": "Diploma in Public Health",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPSPM-DP",
    "description": "Diploma in Sports Medicine",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "UL",
    "description": "Master of Dental Surgery",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "WF",
    "description": "Postgraduate Diploma in Natural Resources",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "WK",
    "description": "Diploma in Industrial Health",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "WP",
    "description": "HTC Diploma",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BEHCJ-DG",
    "description": "Bachelor of Engineering (Honours) (Conjoint)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MPP-DG",
    "description": "Master of Public Policy",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPLAN-DP",
    "description": "Diploma in Languages",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "WH",
    "description": "Doctor of Philosophy/Bachelor of Commerce",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "WS",
    "description": "Non Degree",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "WX",
    "description": "Speech Therapy",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "FD",
    "description": "O'S Teacher's Certificate",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "JZ",
    "description": "To be assigned",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "WQ",
    "description": "Wellesley Programme",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMHG",
    "description": "Master of Forensic Information Technology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMHQ",
    "description": "Master of Landscape Architecture",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMIA",
    "description": "Master of Performance and Media Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMIH",
    "description": "Master of Public Health",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMIJ",
    "description": "Master of Resource Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMIP",
    "description": "Master of Te Reo Maori",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EDEA",
    "description": "Diploma in Sport and Fitness Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EDCEI",
    "description": "Doctor of Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EGDFG",
    "description": "Graduate Diploma in Business",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ENCIY",
    "description": "National Certificate in Agriculture",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ENCJA",
    "description": "National Certificate in Computing",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ENCJF",
    "description": "National Certificate in Pharmacy",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ENCJI",
    "description": "Certificate in Social Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ENCJK",
    "description": "National Certificate in Trades",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ENCJM",
    "description": "National Diploma in Business",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDCR-DP",
    "description": "Postgraduate Diploma in Clinical Pharmacy",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EBAQ",
    "description": "Bachelor of Business Information Systems",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EHAR",
    "description": "Bachelor of Communication Studies (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EHBJ",
    "description": "Bachelor of Fine Arts (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EHBP",
    "description": "Bachelor of Landscape Architecture (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EHBR",
    "description": "Bachelor of Management Studies (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EBBW",
    "description": "Bachelor of Medical Laboratory Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EBCM",
    "description": "Bachelor of Tourism",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ECCT",
    "description": "Certificate in Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ECDC",
    "description": "Certificate of Attainment in English Language",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EHAZ",
    "description": "Bachelor of Creative Technologies (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EDDM",
    "description": "Diploma in Contemporary Music",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "TC",
    "description": "Bachelor of Science/Cert of Proficiency",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NX",
    "description": "Diploma in Arts (P.G.)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "IK",
    "description": "Diploma in Crimnology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "IW",
    "description": "Diploma in Psychology/Cert Proficiency",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "KH",
    "description": "Diploma in Science/Master of Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "FF",
    "description": "Advanced Diploma of Teaching",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BBMED",
    "description": "Bachelor of Biomedical Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ECCV",
    "description": "Certificate in Commerce",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EDEB",
    "description": "Diploma in Sports Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMGB",
    "description": "Magister (Master)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMGC",
    "description": "Master of Agricultural Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMGR",
    "description": "Master of Counselling",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPTCJ-DP",
    "description": "Diploma of Teaching (Jointly Badged)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "CTSES-CT",
    "description": "Certificate of Secondary Education Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "CTTRM-CT",
    "description": "Certificate in Te Reo Maori",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDAC-DP",
    "description": "Postgraduate Diploma in Architecture",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GCSWS-CT",
    "description": "Graduate Certificate in Social Work Supervision",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "HCTEC-CT",
    "description": "Higher Certificate in Educational Computing and Information Technology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPTEC-DP",
    "description": "Diploma of Teaching - Early Childhood Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPDLR",
    "description": "Postgraduate Diploma in Engineering",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPDLZ",
    "description": "Postgraduate Diploma in Information Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPCKU",
    "description": "Postgraduate Certificate in Nursing",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPDMC",
    "description": "Postgraduate Diploma in Medical Laboratory Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPCKN",
    "description": "Postgraduate Certificate in Commerce",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPCLC",
    "description": "Postgraduate Certificate in Sport and Exercise",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0236",
    "description": "National Certificate in Te Ao Turoa",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "C00001",
    "description": "Certificate in Animal Care",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "C00009",
    "description": "Certificate in Community Skills",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "C00010",
    "description": "Certificate in Design and Visual Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "C00017",
    "description": "Certificate in Intensive English",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "D00007",
    "description": "Diploma in Performance Technology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0010",
    "description": "National Certificate in Agribusiness Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0018",
    "description": "National Certificate in Animal Product Examination Services",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0025",
    "description": "National Certificate in Barbering",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0030",
    "description": "National Certificate in Boatbuilding",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0036",
    "description": "National Certificate in Cable Making",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0049",
    "description": "National Certificate in Civil Works and Services ",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0050",
    "description": "National Certificate in Cleaning and Caretaking",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0055",
    "description": "National Certificate in Competitive Manufacturing",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0074",
    "description": "National Certificate in Disability Support Assessment, Planning, and Coordination",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0086",
    "description": "National Certificate in Electronic Security",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0097",
    "description": "National Certificate in Fibre Cement Linings",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0100",
    "description": "National Certificate in Financial Services",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0106",
    "description": "National Certificate in Flooring",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0107",
    "description": "National Certificate in Floristry",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0110",
    "description": "National Certificate in Forest Health Surveillance",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0117",
    "description": "National Certificate in Glass",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0120",
    "description": "National Certificate in Goods Service",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0122",
    "description": "National Certificate in Greyhound Care and Training",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0131",
    "description": "National Certificate in Hot Dip Galvanizing",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0135",
    "description": "National Certificate in Industrial Textile Fabrication",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0138",
    "description": "National Certificate in Intelligence Analysis",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0141",
    "description": "National Certificate in Iwi/Maori Social Services",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0144",
    "description": "National Certificate in Leather",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0167",
    "description": "National Certificate in Nga Mahia te Whare Pora",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0177",
    "description": "National Certificate in Performing Arts ",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0190",
    "description": "National Certificate in Property Consultation and Valuation",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0210",
    "description": "National Certificate in Rural Servicing",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0214",
    "description": "National Certificate in Seafood",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EDCEM",
    "description": "Doctor of Music",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EGDFP",
    "description": "Graduate Diploma in Landscape Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ELGA",
    "description": "Licence",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EDKD",
    "description": "NZ Diploma in Engineering",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPCLB",
    "description": "Postgraduate Certificate in Social Welfare",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0254",
    "description": "National Certificate in Wool ",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0255",
    "description": "National Certificate in Workplace Fire and Emergency Response",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0002",
    "description": "National Diploma in Adult Literacy and Numeracy Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0006",
    "description": "National Diploma in Air Traffic Services",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0027",
    "description": "National Diploma in Electricity Supply",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0032",
    "description": "National Diploma in Fire and Rescue Services",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0045",
    "description": "National Diploma in Maori Business and Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0046",
    "description": "National Diploma in Maori Environmental Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0056",
    "description": "National Diploma in Public Sector Compliance Investment",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0062",
    "description": "National Diploma in Road Transport Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0070",
    "description": "National Diploma in Stevedoring and Ports Industry",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0077",
    "description": "National Diploma in Wastewater Treatment",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0079",
    "description": "National Diploma in Whanau/Family and Foster Care",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MMMT-DG",
    "description": "Master of Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPMUH-DP",
    "description": "Diploma in Music with Honours",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BUPH-DG",
    "description": "Bachelor of Urban Planning (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MDANS-DG",
    "description": "Master of Dance Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDBM-DG",
    "description": "Postgraduate Diploma in Biomedical Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGCHL-DP",
    "description": "Postgraduate Certificate in Health Leadership",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MHL-DG",
    "description": "Master of Health Leadership",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MENS-DG",
    "description": "Master of Environmental Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MEPM-DG",
    "description": "Master of Engineering Project Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GDAS-DP",
    "description": "Graduate Diploma in Architectural Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MTRA-DG",
    "description": "Master of Translation",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BMI-DG",
    "description": "Bachelor of Medical Imaging",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BASHC-DG",
    "description": "Bachelor of Advanced Science (Honours) Conjoint",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MDATS-DG",
    "description": "Master of Data Science",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "MUFY",
    "description": "Monash University Foundation Programme",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "MITE",
    "description": "MIT Cert in Pre-degree Studies Engineering",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "MITT",
    "description": "MIT Cert in Tertiary Degree Studies",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "TSST",
    "description": "Foundation Social Science",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "TUFS",
    "description": "Taylors Foundation in Science",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "NSWF",
    "description": "University of NSW Foundation Studies Certificate",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "OPFS",
    "description": "Otago Polytechnic Certificate in Fdn Studies",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "TUFB",
    "description": "Taylors Foundation in Business",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "MITN",
    "description": "MIT Pre-degree Nursing",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "NTFS",
    "description": "Northtec Certificate in Fdn Studies",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "WFBP",
    "description": "Whitecliffe Foundation One Year Bridging Programme",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "IBIN",
    "description": "International Baccalaureate Diploma",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "DSE",
    "description": "Diploma of Completion of Secondary Education",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "AVG",
    "description": "Avgangsbetyg Upper Secondary School Leaving Cert",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "STPM",
    "description": "Sijil Tinggi Persekolahan Malaysia (STPM)",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "NSCO",
    "description": "Namibia Senior Secondary Cert. (Ordinary Level)",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "BGSE",
    "description": "Botswana General Cert. of Secondary Education",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "SCGA",
    "description": "Singapore/Cambridge GCE A levels",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "HSDK",
    "description": "High School Diploma (Korea)",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "AUS5",
    "description": "South Australian (SACE)",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "EDS",
    "description": "Esame di Stato",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "CGSE",
    "description": "Certificate of General Secondary Education",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "STN2",
    "description": "Steiner School Certificate Level 2",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "CHE",
    "description": "Certificate of Higher Secondary Education (Form 6)",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "HKCE",
    "description": "Hong Kong Certificate of Education",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "HSD",
    "description": "High School Diploma",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "SFG",
    "description": "Slutbetyg fran Grundskola",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GX",
    "description": "Bachelor of Phys Ed/Bachelor of Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "HH",
    "description": "Diploma in Urban Valuation",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "HT",
    "description": "Diploma in Meat Technology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "IC",
    "description": "Diploma in Pharmacy",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "IG",
    "description": "Diploma in Physiotherapy",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "IH",
    "description": "Diploma in School Dental Nursing",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "IS",
    "description": "Diploma in Obstetrics",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "JU",
    "description": "Certificate in Horticultural Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "KK",
    "description": "Diploma in Social Science Research",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "KT",
    "description": "Diploma in Marketing",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MH",
    "description": "Bachelor of Theology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MO",
    "description": "Bachelor of Agricultural Economics",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NG",
    "description": "Bachelor of Commerce and Administration",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NH",
    "description": "Bachelor of Health Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NI",
    "description": "Bachelor of Horticultural Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NJ",
    "description": "Bachelor of Town Planning",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NQ",
    "description": "Diploma in Horticulture",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "OG",
    "description": "Mineral Technology Intermediate",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ON",
    "description": "Medicine Intermediate",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "OQ",
    "description": "Pharmacy Intermediate",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PL",
    "description": "Bachelor of Science (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "QF",
    "description": "Bachelor of Arts (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "RE",
    "description": "Bach of Agricultural Business & Administrn",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "RL",
    "description": "Bachelor of Horticulture",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "RT",
    "description": "Bachelor of Property",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "SF",
    "description": "Diploma in Education (Postgraduate)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "SR",
    "description": "Diploma in Second Language Teaching",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "TK",
    "description": "Bachelor of Law Honours/Law Professionals",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "TM",
    "description": "Bachelor of Law/Bachelor of Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BE-DG",
    "description": "Bachelor of Engineering",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "COPUA-CT",
    "description": "Certificate of Proficiency",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GDARM-DP",
    "description": "Graduate Diploma in Arts Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MBA-DG",
    "description": "Master of Business Administration",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "FO",
    "description": "Diploma in Child Health",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "FS",
    "description": "Diploma in Computer Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "FU",
    "description": "Diploma in Educational Administration",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "FW",
    "description": "Diploma in Farm Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GC",
    "description": "Bachelor of Arts/Bachelor of Commerce",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GI",
    "description": "Bachelor of Arts/Bachelor of Law",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GO",
    "description": "Bachelor of Commerce/Bachelor of Theology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ZP",
    "description": "Certificate in Wool Handling Systems",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPCAE-DP",
    "description": "Diploma in Community Accident & Emergency Medicine",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BOPT-DG",
    "description": "Bachelor of Optometry",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "NCE3",
    "description": "NCEA Level 3",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "BJH",
    "description": "Bahamas Junior Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "CANA",
    "description": "Alberta High School Diploma",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "VWO",
    "description": "VWO Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "BESE",
    "description": "Basic Education High School Exam",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "CESC",
    "description": "Certificado de Educaion",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "CAN6",
    "description": "General High School Diploma (NW Territories)",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "AUS6",
    "description": "Tasmanian Certificate of Education (TCE)",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "YES",
    "description": "Young Enterprise Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "SLC",
    "description": "Senior High School Leaving Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "COHC",
    "description": "Cambridge Overseas Higher School Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "REA",
    "description": "Realschulabschluss/ Mittlere Reife / Mittlerer Sch",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "CSG",
    "description": "Certificat de Studii Gimnaziale",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "CCSE",
    "description": "Certificate of Complete General Secondary Educ.",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "AUS9",
    "description": "Foundation Studies Certificate (University of NSW)",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "CDP",
    "description": "Certificat de Probabtion",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "ICCH",
    "description": "ICCE Advanced (Academic) Certificate Honours",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "AU10",
    "description": "ACT Senior Secondary Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "AUS7",
    "description": "Victorian Certificate of Education (VCE)",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "LSLC",
    "description": "Lower Secondary Leaving Certificate (CSS)",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "COU",
    "description": "Titulo de Bachillerato",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "GCLC",
    "description": "Cert de Absolvire a Ciclului Inferior al Liceului",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "DDD",
    "description": "Diplom der Diplommittelschule",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPBIA-DP",
    "description": "Diploma in Business and Industrial Administration",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPETG-DP",
    "description": "Diploma in Energy Technology (Geothermal)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPOBS-DP",
    "description": "Diploma in Obstetrics",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BPROP-DG",
    "description": "Bachelor of Property",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BMUSE-DG",
    "description": "Bachelor of Music Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BCHCJ-DG",
    "description": "Bachelor of Commerce (Honours) (Conjoint)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MMEDS-DG",
    "description": "Master of Medical Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MCOML-DG",
    "description": "Master of Commercial Law",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DMUS-DG",
    "description": "Doctor of Music",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DENG-DG",
    "description": "Doctor of Engineering",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "LLA",
    "description": "NZ Certificate in Design",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPGYO-DP",
    "description": "Diploma in Gynaecology and Obstetrics",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BAHON-DG",
    "description": "Bachelor of Arts (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BNURS-DG",
    "description": "Bachelor of Nursing",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BNURH-DG",
    "description": "Bachelor of Nursing (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDOM-DP",
    "description": "Postgraduate Diploma in Obstetrics and Medical Gynaecology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MTESL-DG",
    "description": "Master of Teaching English to Speakers of Other Languages",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BEDHN-DG",
    "description": "Bachelor of Education (Teaching) (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPAED-DP",
    "description": "Diploma of Art Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GDEDC-DP",
    "description": "Graduate Diploma of Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPESH-DP",
    "description": "Diploma in Education of Students with Hearing Impairment",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPYS-DP",
    "description": "Diploma in Youth Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPRRT-DP",
    "description": "Diploma for Reading Recovery Tutors",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EDD-DG",
    "description": "Doctor of Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DOCFA-DG",
    "description": "Doctor of Fine Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GDSCI-DP",
    "description": "Graduate Diploma in Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GDTST-DP",
    "description": "Graduate Diploma in Translation Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGCAI-CT",
    "description": "Postgraduate Certificate in Advanced Interpreting",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "COPCU-CT",
    "description": "Certificate of Proficiency for Univ of Canterbury",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "COPLU-CT",
    "description": "Certificate of Proficiency for Lincoln University",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGCPH",
    "description": "Postgraduate Certificate in Public Health",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDPH-DP",
    "description": "Postgraduate Diploma in Public Health",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPSCI-DP",
    "description": "Diploma in Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPEDS-DP",
    "description": "Diploma in Educational Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPHMH-DP",
    "description": "Diploma in Health (Mental Health Nursing)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "LA",
    "description": "NZ Certificate in Advertising",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "US",
    "description": "Master of Public Policy",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "WA",
    "description": "Certificate of Proficiency",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "WE",
    "description": "Postgraduate Diploma in Commerce",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "WW",
    "description": "Science Certificate",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "XC",
    "description": "Law Intermediate and LLB",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MENRG-DG",
    "description": "Master of Energy",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ZS",
    "description": "Cultural Class",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "KD",
    "description": "NZ Chartered Institute of Secretaries",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMHO",
    "description": "Master of International Relations",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMIE",
    "description": "Master of Professional Accounting",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMIF",
    "description": "Master of Professional Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMIG",
    "description": "Master of Property Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ECIX",
    "description": "MIT Certificate in Pre-Degree Studies (Engineering)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EBBD",
    "description": "Bachelor of Design",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EBCD",
    "description": "Bachelor of Resource and Environmental Planning",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EGDFY",
    "description": "Graduate Diploma in Valuation",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMIT",
    "description": "Master of Theology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ENCJC",
    "description": "National Certificate in Engineering",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ENCJU",
    "description": "National Diploma in Real Estate",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EBAH",
    "description": "Bachelor of Architectural Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EBBG",
    "description": "Bachelor of Education (Teaching)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EBBS",
    "description": "Bachelor of Maori and Pacific Development",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EBBT",
    "description": "Bachelor of Maori Visual Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EBCF",
    "description": "Bachelor of Social Practice",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EHCL",
    "description": "Bachelor of Theology (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EDDI",
    "description": "Diploma in Applied Technology (Building)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EHAN",
    "description": "Bachelor of Business (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GA",
    "description": "Bachelor Arts(Hons)/Arch Intermediate",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MCE-DG",
    "description": "Master of Commercialisation and Entrepreneurship",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EBAS",
    "description": "Bachelor of Computer Graphic Design",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EDDY",
    "description": "Diploma in Social and Community Work",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EGDFV",
    "description": "Graduate Diploma in Teaching",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMGJ",
    "description": "Master of Building Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMGT",
    "description": "Master of Creative Writing",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMGU",
    "description": "Master of Dance Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMGY",
    "description": "Master of Development Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDMU-DP",
    "description": "Postgraduate Diploma in Music",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDOR-DP",
    "description": "Postgraduate Diploma in Operations Research",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPTLB-DP",
    "description": "Diploma of Teacher Librarianship",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPASS-DP",
    "description": "Advanced Diploma in Social Work Supervision",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPDLT",
    "description": "Postgraduate Diploma in Financial Analysis",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPDMJ",
    "description": "Postgraduate Diploma in Physical Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPDMN",
    "description": "Postgraduate Diploma in Public Policy",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPCKQ",
    "description": "Postgraduate Certificate in Health Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPDLL",
    "description": "Postgraduate Diploma in Clinical Pharmacy",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EDNG",
    "description": "Teachers College Diploma",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPCLD",
    "description": "Postgraduate Certificate in Strategic Leadership",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPDLP",
    "description": "Postgraduate Diploma in Economics",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "CTEAP-CT",
    "description": "CertAcadPrep",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPDMX",
    "description": "Postgraduate Diploma in Teaching",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPCKK",
    "description": "Postgraduate Certificate in Applied Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPDLQ",
    "description": "Postgraduate Diploma in Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPDND",
    "description": "Postgraduate Diploma in Development Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0219",
    "description": "National Certificate in Sign Making",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0220",
    "description": "National Certificate in Snow Sport",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "C00015",
    "description": "Certificate in Home Garden Design",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "C00016",
    "description": "Certificate in Information Technology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "D00002",
    "description": "Diploma in Community and Social Work",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "D00010",
    "description": "Diploma in Tourism Leadership and Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "D00011",
    "description": "Diploma in Tourism Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0009",
    "description": "National Certificate in Aeronautical Storekeeping",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0016",
    "description": "National Certificate in Amenity Turf Maintenance",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0048",
    "description": "National Certificate in Civil Plant Operation ",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0052",
    "description": "National Certificate in Commercial Road Transport ",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0056",
    "description": "National Certificate in Compliance and Regulatory Control",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0069",
    "description": "National Certificate in Dairy Manufacturing",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0083",
    "description": "National Certificate in Electrical Engineering",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0094",
    "description": "National Certificate in Farming",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0099",
    "description": "National Certificate in Fibrous Plaster",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0111",
    "description": "National Certificate in Forest Operations Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0114",
    "description": "National Certificate in Freight Forwarding",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0116",
    "description": "National Certificate in Gas",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0121",
    "description": "National Certificate in Governance of Maori Authorities",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0136",
    "description": "National Certificate in Infrastructure Civil Engineering",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0154",
    "description": "National Certificate in Maritime",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0156",
    "description": "National Certificate in Masonry",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0162",
    "description": "National Certificate in Mental Health and Addiction Support ",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0168",
    "description": "National Certificate in Occupational Health and Safety ",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0180",
    "description": "National Certificate in Plastics Processing Technology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0195",
    "description": "National Certificate in Racing Operations",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0198",
    "description": "National Certificate in Real Estate",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0205",
    "description": "National Certificate in Resource Recovery",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EDDT",
    "description": "Diploma in Language and Culture",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EDDX",
    "description": "Diploma in Professional Accountancy",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EDEF",
    "description": "Diploma of Specialist",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EGCEP",
    "description": "Graduate Certificate in Applied Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EGCER",
    "description": "Graduate Certificate in Business",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EGDFK",
    "description": "Graduate Diploma in Economics",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMGI",
    "description": "Master of Biomedical Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPDLS",
    "description": "Postgraduate Diploma in Environmental Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0246",
    "description": "National Certificate in Waste Water Treatment",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0250",
    "description": "National Certificate in Whakairo",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0008",
    "description": "National Diploma in Ambulance Practice",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0018",
    "description": "National Diploma in Community Work",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0028",
    "description": "National Diploma in Embalming",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0038",
    "description": "National Diploma in Human Services",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0052",
    "description": "National Diploma in Pork Production",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0064",
    "description": "National Diploma in Seafood Vessel Operations",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0072",
    "description": "National Diploma in Textile Dyeing and Finishing",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ECDF",
    "description": "Certificate of University Preparation",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MTEP-DG",
    "description": "Master of Teaching (Primary)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGCED-CT",
    "description": "Postgraduate Certificate in Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MHLP-DG",
    "description": "Master of Health Practice",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDIT-DG",
    "description": "Postgraduate Diploma in Information Technology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MHP-DG",
    "description": "Master of Health Psychology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDHL-DP",
    "description": "Postgraduate Diploma in Health Leadership",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDMI-DP",
    "description": "Postgraduate Diploma in Māori and Indigenous Business",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGCIT-CT",
    "description": "Postgraduate Certificate in Information Technology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MARHC-DG",
    "description": "Master of Architecture (Professional) and Heritage Conservation",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MMS-DG",
    "description": "Master of Marine Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MHIED-DG",
    "description": "Master of Higher Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGCTR-CT",
    "description": "Postgraduate Certificate in Translation",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BEDEC-DG",
    "description": "Bachelor of Education (Teaching)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BEDP-DG",
    "description": "Bachelor of Education (Teaching)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MDMT-DG",
    "description": "Master of Dance Movement Therapy",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "USFP",
    "description": "University of Sydney Foundation Programme",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "MITF",
    "description": "MIT Foundation all specialty programmes",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "TAFY",
    "description": "Auckland Foundation Year",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "HSC",
    "description": "Higher School Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "MAU",
    "description": "Maturitat",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "IB",
    "description": "International Baccalaureate (NZ)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GT",
    "description": "Bachelor of Education/Bachelor of Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "HF",
    "description": "Registered Psychiatric Nurse",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "HO",
    "description": "Diploma in Industrial Relations",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "HY",
    "description": "Diploma in Occupational Health",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "IE",
    "description": "Diploma in Occupational Therapy",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "IP",
    "description": "Diploma in Planning",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "JX",
    "description": "Certificate in Arts Humanities",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MG",
    "description": "Bachelor of Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MJ",
    "description": "Bachelor of Physical Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MK",
    "description": "Bachelor of Management Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MR",
    "description": "Bachelor of Science (Technology)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MT",
    "description": "Bachelor of Pharmacy",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NF",
    "description": "Bachelor of Engineering",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NK",
    "description": "Bachelor of Architecture",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NU",
    "description": "Diploma in Parks and Recreation",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NW",
    "description": "Diploma in Physical Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NY",
    "description": "Diploma in Speech Therapy",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "OB",
    "description": "Architecture Intermediate",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PU",
    "description": "Master of Home Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "QB",
    "description": "Bachelor of Business Studies (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "RD",
    "description": "Bachelor of Agricultural Commerce",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "RG",
    "description": "Bachelor of Consumer and Applied Sciences",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "RH",
    "description": "Bachelor of Computing and Mathematical Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "SJ",
    "description": "Diploma in Educational Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "SS",
    "description": "Diploma in Town Planning",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "SX",
    "description": "Diploma in Health Administration",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MEDMG-DG",
    "description": "Master of Educational Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MEMGT-DG",
    "description": "Master of Engineering Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GDBUS-DP",
    "description": "Graduate Diploma in Business",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "FJ",
    "description": "Diploma in Aviation",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "FP",
    "description": "Diploma in Business Data Processing",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "FQ",
    "description": "Teachers College Diploma",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GB",
    "description": "Bachelor of Arts(Hons)/Bachelor of Law",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GH",
    "description": "Bachelor of Arts/Bachelor of Law (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GJ",
    "description": "Bach Business Studies/Dip Banking Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GL",
    "description": "B. of Commerce Honours/Bachelor of Law",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GM",
    "description": "Bachelor of Commerce/Bachelor of Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ZI",
    "description": "Certificate in Japanese Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ZJ",
    "description": "Cert in Labour and Trade Union Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ZK",
    "description": "Certificate in Liberal Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BHB-DG",
    "description": "Bachelor of Human Biology",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "VEP",
    "description": "Vestibular/ENEM/PAS",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "CAN3",
    "description": "High School Graduation Diploma (Manitoba)",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "EACE",
    "description": "East African Certificate of Education",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "BAC",
    "description": "Bachiller/Bachillerato",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "GPC",
    "description": "General Prepartory Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "MAW6",
    "description": "Maw 6",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "CSEC",
    "description": "Caribbean Exam Council Secondary Education Cert",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "STUS",
    "description": "Special Training School Upper Secondary Cert",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "GECS",
    "description": "General Secondary Education Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "SCGO",
    "description": "Singapore/Cambridge GCE O levels",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "BHEC",
    "description": "Bhutan Higher Secondary Education Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "CAN5",
    "description": "High School Graduation Diploma (Newfoundland)",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "MHSC",
    "description": "Malaysian Higher School Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "SPFS",
    "description": "South Pacific Form Seven Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "HSD",
    "description": "High School Diploma (International USA school)",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "HLC",
    "description": "Higher Leaving Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "CAN7",
    "description": "High School Completion Certificate (Nova Scotia)",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "DES",
    "description": "Diplome d'Etat d'Etudes Secondaires du Cycle Long",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "HAU",
    "description": "Hauptschulabschluss",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "UMCE",
    "description": "Uttar Madhyama Certificate Examination",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "SPFS",
    "description": "South Pacific Form Seven",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "USSC",
    "description": "Upper Secondary School Leaving Certificate (KSS)",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "YLP",
    "description": "Ylioppilastutkinto/Studentexamen",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "STTB",
    "description": "Cert of Completion of Academic Secondary School",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "QSSC",
    "description": "Qatar Senior School Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "CSBM",
    "description": "NZIM Certificate in Small Business Management",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "NZHC",
    "description": "NZ Higher School Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "TBD",
    "description": "Título de Bachillerato Diversificado",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "KCSE",
    "description": "Kenya Certificate of Secondary Education",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "SSCE",
    "description": "Senior Secondary School Certificate Examination",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "SDM",
    "description": "Swiadectwo Dojrzalosci/Matura",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "CES",
    "description": "Certificado de fim de Estudo Secundarios",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPFA-DP",
    "description": "Diploma in Fine Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPGCO-DP",
    "description": "Diploma in Guidance and Counselling",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MHB-DG",
    "description": "Master of Human Biology (Physiology)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BARCH-DG",
    "description": "Bachelor of Architecture",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MBCHB-DG",
    "description": "Bachelor of Medicine and Bachelor of Surgery",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BTECH-DG",
    "description": "Bachelor of Technology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BMUSC-DG",
    "description": "Bachelor of Music (Conjoint)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MMUS-DG",
    "description": "Master of Music",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MD-DG",
    "description": "Doctor of Medicine",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "LI",
    "description": "NZ Certificate in Customs",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "LO",
    "description": "NZ Certificate in Forestry",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "LU",
    "description": "NZ Certificate in Library Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGCHC-CT",
    "description": "PGCertHealth(ChildAdolescentMentalHealth)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BPRPH-DG",
    "description": "Bachelor of Property (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPOR-DP",
    "description": "Diploma in Operations Research",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGCLM-CT",
    "description": "Postgraduate Certificate in Light Metals Reduction",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GDHSC-DP",
    "description": "Graduate Diploma in Health Sciences",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GDPFA-DP",
    "description": "Graduate Diploma in Performing Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GDTPR-DP",
    "description": "Graduate Diploma in Teaching (Primary)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MMH-DG",
    "description": "Master of Maori Health",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDMH-DP",
    "description": "Postgraduate Diploma in Maori Health",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPADT-DP",
    "description": "Advanced Diploma of Teaching",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDOE-DP",
    "description": "Postgraduate Diploma of Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDLE-DP",
    "description": "Postgraduate Diploma in Literacy Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "HCTBE-CT",
    "description": "Higher Certificate in Bilingual/Immersion Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GCPS-CT",
    "description": "Graduate Certificate in Professional Supervision",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPTCE-DP",
    "description": "Diploma in Technology Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDED-DP",
    "description": "Postgraduate Diploma in Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "CTART-CT",
    "description": "Certificate in Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDSM-DP",
    "description": "Postgraduate Diploma in Sports Medicine",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPCLP-DP",
    "description": "Diploma in Clinical Psychology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPPFA-DP",
    "description": "Diploma in Performing Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDAR-DP",
    "description": "Postgraduate Diploma in Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPHCM-DP",
    "description": "Diploma in Health (Case Management)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPGMD-DP",
    "description": "Diploma in Geriatric Medicine",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "TP",
    "description": "Bachelor of Law/Bachelor of Commerce",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "TR",
    "description": "Bachelor of Law/Bachelor of Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "UG",
    "description": "Master of Art (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "UJ",
    "description": "Master of Management Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "UP",
    "description": "Master of Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "UQ",
    "description": "Master of Pharmacy",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "UW",
    "description": "Master of Philosophy",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "CTLAN-CT",
    "description": "Certificate in Languages",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "WM",
    "description": "FLRTCert",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "WR",
    "description": "Non-Matriculated",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "CC",
    "description": "Ontario Academic Credits",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMHI",
    "description": "Master of Health Sciences",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMHJ",
    "description": "Master of Information Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMHP",
    "description": "Master of International Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMHS",
    "description": "Master of Maori and Indigenous Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMIR",
    "description": "Master of Technology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMIV",
    "description": "Masters Degree",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EBAY",
    "description": "Bachelor of Creative Technologies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EHCR",
    "description": "Bachelor of Viticulture and Oenology (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EDCEJ",
    "description": "Doctor of Juridical Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EMHU",
    "description": "Master of Medical Laboratory Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ENCJD",
    "description": "National Certificate in Hairdressing",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ENCJE",
    "description": "National Certificate in Hospitality",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ENCJL",
    "description": "National Diploma in Architectural Technology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EHAT",
    "description": "Bachelor of Computer Graphic Design (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EBAV",
    "description": "Bachelor of Construction",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EHAW",
    "description": "Bachelor of Consumer and Applied Sciences (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EBBE",
    "description": "Bachelor of Design and Visual Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EBBO",
    "description": "Bachelor of Landscape Architecture",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EHBX",
    "description": "Bachelor of Medical Science (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ECCU",
    "description": "Certificate in Business Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EBBV",
    "description": "Bachelor of Media and Creative Technologies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EBCE",
    "description": "Bachelor of Social and Community Work",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "TQ",
    "description": "Bach Law/Bach Commerce/Law Professionals",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ZD",
    "description": "Certificate in Dairy Factory Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GDIE-DP",
    "description": "Graduate Diploma in Innovation and Entrepreneurship",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "COPUF-CT",
    "description": "Certificate of Proficiency",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EBAE",
    "description": "Bachelor of Applied Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EHAF",
    "description": "Bachelor of Applied Science (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MBLDG-DG",
    "description": "Master of Building Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GDOR-DP",
    "description": "Graduate Diploma in Operations Research",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NDAET-DP",
    "description": "National Diploma in Adult Education and Training",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDIB-DP",
    "description": "Postgraduate Diploma in International Business",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDBE-DP",
    "description": "Postgraduate Diploma in Bioscience Enterprise",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPASW-DP",
    "description": "Advanced Diploma in Social Work",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NDBED-DP",
    "description": "National Diploma of Business Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DPTKK-DP",
    "description": "Diploma of Teaching (Kura Kaupapa Maori)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDCE-DP",
    "description": "Postgraduate Diploma in Clinical Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "CTGIS-CT",
    "description": "Certificate in Global Issues",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BEHON-DG",
    "description": "Bachelor of Engineering (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPDLU",
    "description": "Postgraduate Diploma in Health Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "CTACP-CT",
    "description": "Certificate in Academic Preparation",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPCLF",
    "description": "Postgraduate Certificate in Tourism Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPDMM",
    "description": "Postgraduate Diploma in Public Health",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPDMW",
    "description": "Postgraduate Diploma in Te Reo Maori",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPCKX",
    "description": "Postgraduate Certificate in Public Health",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EDKF",
    "description": "New Zealand Law Society Legal Executive Diploma",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0223",
    "description": "National Certificate in Solar Water Heating Installation",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0225",
    "description": "National Certificate in Solid Wood",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0226",
    "description": "National Certificate in Specialist Interiors",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0229",
    "description": "National Certificate in Sport Turf ",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0231",
    "description": "National Certificate in Steel Manufacturing",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "D00005",
    "description": "Diploma in Graphic Design and Animation",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "D00008",
    "description": "Diploma in Sustainable Horticulture",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "D00009",
    "description": "Diploma in Teaching (Early Childhood)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0002",
    "description": "National Certificate in Administration of Revenue Law",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0013",
    "description": "National Certificate in Aircraft Servicing",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0015",
    "description": "National Certificate in Ambulance",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0021",
    "description": "National Certificate in Arable Farming",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0022",
    "description": "National Certificate in Architectural Aluminium Joinery",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0034",
    "description": "National Certificate in Building, Construction and Allied Trades Skills",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0061",
    "description": "National Certificate in Construction",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0066",
    "description": "National Certificate in Cranes",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0068",
    "description": "National Certificate in Dairy Farming",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0080",
    "description": "National Certificate in Driving",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0093",
    "description": "National Certificate in Extractive Industries",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0101",
    "description": "National Certificate in Fire and Rescue",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0112",
    "description": "National Certificate in Forestry",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0115",
    "description": "National Certificate in Furniture",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0119",
    "description": "National Certificate in Glazing",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0134",
    "description": "National Certificate in Industrial Rope Access",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0158",
    "description": "National Certificate in Meat Processing",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0171",
    "description": "National Certificate in Pacific Island Early Childhood Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0181",
    "description": "National Certificate in Plumbing",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0183",
    "description": "National Certificate in Pork Production",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0194",
    "description": "National Certificate in Racing Broadcasting",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0197",
    "description": "National Certificate in Rail",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0204",
    "description": "National Certificate in Resouce Efficiency",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0208",
    "description": "National Certificate in Road Transport Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0209",
    "description": "National Certificate in Roofing",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0213",
    "description": "National Certificate in Scrap Metal Recycling",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EDDR",
    "description": "Diploma in Landscape Design",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EDDW",
    "description": "Diploma in Product Design Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EDCEL",
    "description": "Doctor of Medicine",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EGCEQ",
    "description": "Graduate Certificate in Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EGCEY",
    "description": "Graduate Certificate in Recreation Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EGCFB",
    "description": "Graduate Certificate in Social Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EGDFD",
    "description": "Graduate Diploma (any other field)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EGDFL",
    "description": "Graduate Diploma in Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EGDFW",
    "description": "Graduate Diploma in Theology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "EPCKM",
    "description": "Postgraduate Certificate in Business",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NC0240",
    "description": "National Certificate in Textiles",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0003",
    "description": "National Diploma in Aero Main Certification",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0005",
    "description": "National Diploma in Agribusiness Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0011",
    "description": "National Diploma in Boatbuilding",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0019",
    "description": "National Diploma in Competitive Manufacturing",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0029",
    "description": "National Diploma in Employment Support",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0039",
    "description": "National Diploma in Industrial Machine Knitting",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0044",
    "description": "National Diploma in Laboratory Animal Care",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0047",
    "description": "National Diploma in Maori Performing Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0066",
    "description": "National Diploma in Security Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0071",
    "description": "National Diploma in Te Matauranga Maori",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NGD001",
    "description": "National Graduate Diploma in Electricity Supply",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MINS-DG",
    "description": "Master of Indigenous Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PGDCS-DP",
    "description": "Postgraduate Diploma in Conflict and Terrorism Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BGSC-DG",
    "description": "Bachelor of Global Studies (Conjoint)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MCTS-DG",
    "description": "Master of Conflict and Terrorism Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MINBU-DG",
    "description": "Master of International Business",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MEP-DG",
    "description": "Master of Education Practice",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "DG-MTRA",
    "description": "Obsolete",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BEDPT-DG",
    "description": "Bachelor of Education (Teaching)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BEDTT-DG",
    "description": "Bachelor of Education (Teaching)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ND0001",
    "description": "National Diploma in Adult Education and Training",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "FSSC",
    "description": "Foundation Social Science",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "UCFS",
    "description": "Unitec Certificate of Foundation Studies",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "WMCP",
    "description": "Certificate of University Preparation (Massey Uni)",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "UCUP",
    "description": "Certificate of University Preparation (Massey Uni)",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "ANUF",
    "description": "Australian National University Foundation Studies",
    "country": "N/A"
  },
  {
    "type": "foundation",
    "code": "KFSC",
    "description": "Foundation Science",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "HSDI",
    "description": "High School Diploma (International USA School)",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "EAAC",
    "description": "East African Advanced Certificate of Education",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "HAVO",
    "description": "HAVO Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "TDB",
    "description": "Titulo de Bachiller",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "SULO",
    "description": "Swiadectwo Ukonczenia Lyceum Ogolnoksztalcacego",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "CCSG",
    "description": "Certificado de Conclusão de Segundo Grau",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "GCA",
    "description": "GCE/CIE A-Levels",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "CAN1",
    "description": "General High School Diploma (Alberta)",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "AUS3",
    "description": "Northern Territory Certificate of Education (NTCE)",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "HF",
    "description": "Bevis for Højere Forberedelseseksamen",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "TGE",
    "description": "Titulo de Graduado Escolar",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "TBC",
    "description": "Titulo de Bachilleratio Comercial",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "AU11",
    "description": "Australian Tertiary Admission Rank",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "ACSE",
    "description": "Advanced Certificate of Secondary Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GU",
    "description": "Bachelor of Education/Diploma of Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GZ",
    "description": "Bachelor of Science/Bachelor of Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "HB",
    "description": "Diploma in Business Studies/Cert of Prof",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "HJ",
    "description": "Diploma in Home Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "HN",
    "description": "Diploma in Humanities",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "HP",
    "description": "Diploma in Instructional Systems",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "HS",
    "description": "Diploma in Maths Education",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "HW",
    "description": "Diploma in Museum Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "JA",
    "description": "Technician's Certificate (Radio)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "JH",
    "description": "Diploma in Textile Design",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "JK",
    "description": "Diploma in Business Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "JW",
    "description": "Diploma in Hotel & Catering Administration",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "KC",
    "description": "Social Studies Certificate",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "KR",
    "description": "Diploma in Wool and Wool Technology",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "KS",
    "description": "Diploma in Womens Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MD",
    "description": "Bachelor of Dental Surgery",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ME",
    "description": "Bachelor of Surveying",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MF",
    "description": "Bachelor of Business Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "MZ",
    "description": "Bachelor of Management Studies with Honours",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NA",
    "description": "Bachelor of Agricultural Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "NT",
    "description": "Diploma in Natural Resources",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "OD",
    "description": "Dentistry Intermediate",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "OE",
    "description": "Home Science Intermediate",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "OK",
    "description": "Property Administration Intermediate",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "OU",
    "description": "Technology Intermediate",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PA",
    "description": "Master of Agricultural Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PB",
    "description": "Master of Agriculture Business and Administration",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PF",
    "description": "Master of Engineering",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "PI",
    "description": "Master of Divinity",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "QC",
    "description": "Bachelor of Education (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "QD",
    "description": "Bachelor of Music (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "QI",
    "description": "Bachelor of Social Sciences with Honours",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "QL",
    "description": "Bachelor of Agricultural Science (Honours)",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "QN",
    "description": "Bachelor of Commerce and Administration with Honours",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "RI",
    "description": "Bachelor of Commerce and Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "SE",
    "description": "Diploma in Economics",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "SH",
    "description": "Diploma in Educ Guidance and Counselling",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "SU",
    "description": "Diploma in Social Work",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "SV",
    "description": "Diploma in Horticultural Management",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "TD",
    "description": "Bachelor of Science/Bachelor of Music",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "TL",
    "description": "Bachelor of Law/Bachelor of Arts Honours",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BCOM-DG",
    "description": "Bachelor of Commerce",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "BA-DG",
    "description": "Bachelor of Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "FC",
    "description": "NZ Trained Teacher's Certificate",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "FG",
    "description": "Diploma in Amenity Horticulture",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "FI",
    "description": "Diploma in Applied Science",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "FZ",
    "description": "Diploma for Graduates",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GQ",
    "description": "Bach. Commerce/Engineering Intermediate",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GR",
    "description": "Bachelor of Commerce/Bachelor of Law",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "GS",
    "description": "Bachelor of Divinity/Bachelor of Arts",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ZA",
    "description": "Certificate of Attainment in Law Related Educatn",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ZM",
    "description": "Certificate in Rehabilitation Studies",
    "country": "N/A"
  },
  {
    "type": "tertiary",
    "code": "ZN",
    "description": "Certificate in Social and Community Work",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "IGCS",
    "description": "Intl. General Cert. of Secondary Education",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "NZ6C",
    "description": "New Zealand  Sixth Form Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "TAW",
    "description": "Tawjihuyah (General Secondary Education Cert)",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "SSC",
    "description": "Sudan School Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "KHSD",
    "description": "Kosovo High School Diploma",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "HSI4",
    "description": "High School Diploma (Pre 1992)",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "BUP",
    "description": "Graduado en Educacion Secundaria",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "HHLC",
    "description": "Hong Kong Higher Level Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "NZBS",
    "description": "New Zealand Bursary Examinations",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "DCCE",
    "description": "Diploma/Cert. of Completion of Compulsory Educ.",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "SGE",
    "description": "Cert. of Secondary (Complete) General Education",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "SGEA",
    "description": "Sri Lankan General Cert of Education (Advanced)",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "CESS",
    "description": "Certificat d'Enseignement Secondaire Supérieur",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "NZSL",
    "description": "New Zealand Scholarship Examination",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "AVVI",
    "description": "Attestation of General Secondary Education",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "SSLC",
    "description": "Secondary School Leaving Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "SSC",
    "description": "Secondary School Leaving Certificate / Matura",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "SOM",
    "description": "Svjedodzba o Maturi (Certificate of Maturity)",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "JSC",
    "description": "Junior Secondary School Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "TAW",
    "description": "Tawjihiyya (General, Religious and Vocational)",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "CHL",
    "description": "Certificado de Habilitacoes Literarias",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "TBH",
    "description": "Titulo de Bachilleratio Humanistico - Scientifico",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "TBA",
    "description": "Titulo de Bachiller Academico",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "ACT",
    "description": "ACT Test",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "PSSC",
    "description": "Pacific Secondary School Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "AUS0",
    "description": "Australian Foundation/Pre-University study",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "STU",
    "description": "Studentsprof (from Gymnasium)",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "WASC",
    "description": "West African School Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "SSLC",
    "description": "Senior Secondary School Leaving Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "CHIN",
    "description": "Senior Secondary School Certificate",
    "country": "N/A"
  },
  {
    "type": "secondary",
    "code": "MZK",
    "description": "Maturitni Zkouska/Maturita",
    "country": "N/A"
  }
]`)

		default:
			w.WriteHeader(http.StatusNotFound)
			io.WriteString(w, fmt.Sprintf("%q NOT FOUND!", ru))
		}
		if malformatResponse {
			io.WriteString(w, "~~~THIS IS NOT SUPPOSED TO BE HERE~~~")
		}
	})
}

package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

// http.HandlerFunc(

func createMockHandler(t *testing.T) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		url, ru := r.URL, r.URL.RequestURI()
		switch {
		case ru == "/ping":
			w.WriteHeader(http.StatusNoContent)
		case ru == "/oauth/token":
			io.WriteString(w, `{"access_token": "7jsxDZceygy2xNbK2M23sD5eyHimtx", "expires_in": 86400, "token_type": "Bearer", "scope": ""}`)
		case ru == "/api/v1/tasks?type=AFFILIATION" || ru == "/api/v1/tasks?type=AFFILIATION&staus=INACTIVE":
			io.WriteString(w, `[
	{"created-at":"2019-07-24T08:47:09","filename":"UOA-OH-INTEGRATION-TASK-pv51ql.json","id":781,"records":[],"status":"ACTIVE","task-type":"AFFILIATION","updated-at":"2019-07-24T09:29:24"},
	{"created-at":"2019-07-25T00:34:08","filename":"UOA-OH-INTEGRATION-TASK-pv69kw.json","id":787,"records":[],"task-type":"AFFILIATION","updated-at":"2019-07-25T01:32:36"}`)
			if withAnIncomleteTask {
				io.WriteString(w, `,{"created-at":"2099-07-25T00:34:08","filename":"UOA-OH-INTEGRATION-TASK-pv69kZ.json","id":888,"records":[],"task-type":"AFFILIATION"}`)
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
			} else if id == "rcir178@auckland.ac.nz" {
				io.WriteString(w, `[


]`)
			} else {
				w.WriteHeader(http.StatusNotFound)
				io.WriteString(w, `{"error": "User with specified identifier 'rcir178ABC@auckland.ac.nz' not found."}`)
			}
		case strings.HasPrefix(ru, "/api/v1/tasks/"):
			if r.Method == "POST" {
				filename := url.Query()["filename"][0]
				io.WriteString(w, `{
				"id":99999,
				"created-at":"2032-08-25T02:07:28",
				"filename":"`+filename+`",
				"task-type":"AFFILIATION",
				"updated-at":"2032-07-25T02:23:32"
			}`)
			} else {
				io.WriteString(w, `{
				"created-at":"2032-08-25T02:07:28",
				"filename":"UOA-OH-INTEGRATION-TASK-pv6xi6.json",
				"id":`+strings.TrimPrefix(ru, "/api/v1/tasks/")+`,
				"status":"ACTIVE",
				"task-type":"AFFILIATION",
				"updated-at":"2032-07-25T02:23:32"
			}`)
			}
		case strings.HasPrefix(ru, "/api/v1/tokens/"):
			io.WriteString(w, `[
				{
					"access_token": "ecf16b31-ad54-4ba2-ae55-e97fb90e211a",
					"expires_in": 631138518,
					"refresh_token": "a6c9da20-31be-442a-9faa-73f1d92fac45",
					"scopes": "/read-limited,/activities/update"
				}
			]`)
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
		case strings.HasPrefix(ru, "/service/identity/integrations/v3/identity/"):
			var uid = strings.TrimPrefix(ru, "/service/identity/integrations/v3/identity/")
			switch uid {
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
		default:
			w.WriteHeader(http.StatusNotFound)
			io.WriteString(w, fmt.Sprintf("%q NOT FOUND!", ru))
		}
	})
}

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/stretchr/testify/assert"
)

var server *httptest.Server

func SetupTest(t *testing.T, withAnIncomleteTask bool) {
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		ru := r.URL.RequestURI()
		switch {
		case ru == "/oauth/token":
			io.WriteString(w, `{"access_token": "7jsxDZceygy2xNbK2M23sD5eyHimtx", "expires_in": 86400, "token_type": "Bearer", "scope": ""}`)
		case ru == "/api/v1/tasks?type=AFFILIATION":
			io.WriteString(w, `[
	{"created-at":"2019-07-24T08:47:09","filename":"UOA-OH-INTEGRATION-TASK-pv51ql.json","id":781,"records":[],"status":"ACTIVE","task-type":"AFFILIATION","updated-at":"2019-07-24T09:29:24"},
	{"created-at":"2019-07-25T00:34:08","filename":"UOA-OH-INTEGRATION-TASK-pv69kw.json","id":787,"records":[],"task-type":"AFFILIATION","updated-at":"2019-07-25T01:32:36"}`)
			if withAnIncomleteTask {
				io.WriteString(w, `,{"created-at":"2099-07-25T00:34:08","filename":"UOA-OH-INTEGRATION-TASK-pv69kZ.json","id":888,"records":[],"task-type":"AFFILIATION"}`)
			}
			io.WriteString(w, "]")
		case strings.HasPrefix(ru, "/api/v1/tokens/"):
			var id = strings.TrimPrefix(ru, "/api/v1/tokens/")
			if id == "rad42@mailinator.com" || id == "0000-0001-8228-7153" {
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
				io.WriteString(w, `[]`)
			} else {
				w.WriteHeader(http.StatusNotFound)
				io.WriteString(w, `{"error": "User with specified identifier 'rcir178ABC@auckland.ac.nz' not found."}`)
			}
		case strings.HasPrefix(ru, "/api/v1/tasks/"):
			io.WriteString(w, `{
				"created-at":"2032-08-25T02:07:28",
				"filename":"UOA-OH-INTEGRATION-TASK-pv6dwg.json",
				"id":`+strings.TrimPrefix(ru, "/api/v1/tasks/")+`,
				"status":"ACTIVE",
				"task-type":"AFFILIATION",
				"updated-at":"2032-07-25T02:23:32"
			}`)
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
		case strings.HasPrefix(ru, "/service/identity/integrations/v3/identity/"):
			var uid = strings.TrimPrefix(ru, "/service/identity/integrations/v3/identity/")
			switch uid {
			case "rpaw053":
				io.WriteString(w, `{
   "deceased":{
      "dead":false
   },
   "disabilityInfo":{
      "disabilities":[

      ],
      "isDisabled":false
   },
   "displayName":"Roshan Prakash ABC",
   "dob":"1790-05-08",
   "emailAddress":"roshan_pawarasjdfkasdjfajs_@auckland.ac.nz",
   "emails":[
      {
         "email":"skajdfkljsadlf@auckland.ac.nz",
         "lastUpdated":"2017-01-13T17:12:23.000+0000",
         "typeId":"Campus",
         "type":"University",
         "verified":false
      },
      {
         "email":"rpfkjds@aucklanduni.ac.nz",
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
   "firstName":"Roshan Prakash",
   "id":208013283,
   "idPhotoExists":true,
   "lastName":"Pawar",
   "names":[
      {
         "first":"Roshan Prakash",
         "last":"Pawar",
         "lastUpdated":"2015-03-23T21:39:18.000+0000",
         "title":"Mr",
         "type":"Preferred"
      },
      {
         "first":"Roshan Prakash",
         "last":"Pawar",
         "lastUpdated":"2015-03-31T21:02:51.000+0000",
         "title":"Mr",
         "type":"Primary"
      }
   ],
   "primaryIdentity":true,
   "residency":"I",
   "upi":"rpaw053",
   "whenUpdated":"2019-05-14T10:51:47.597+0000",
   "resolved":true,
   "previousIds":[]}`)
			case "rcir178":
				io.WriteString(w, `{
    "displayName": "Radomirs Cirskis",
    "dob": "1896-12-28",
    "emailAddress": "kfsjdjadffsafkis@auckland.ac.nz",
    "emails": [
	{
	    "email": "kfsjdjadffsafkis@auckland.ac.nz",
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
    "firstName": "Radomirs",
    "gender": "MALE",
    "id": 484378182,
    "lastName": "Cirskis",
    "mobile": "+64221221442",
    "names": [
	{
	    "first": "Radomirs",
	    "last": "Cirskis",
	    "lastUpdated": "2017-01-19T20:53:57.000+0000",
	    "type": "Primary"
	},
	{
	    "first": "Radomirs",
	    "last": "Cirskis",
	    "lastUpdated": "2019-05-13T00:34:04.000+0000",
	    "title": "Mr",
	    "type": "Preferred"
	}
    ],
    "previousIds": [],
    "primaryIdentity": true,
    "residency": "NZ-PR",
    "resolved": true,
    "upi": "rcir178",
    "whenUpdated": "2019-05-13T00:34:04.698+0000"
}`)
			case "rad42":
			case "non-existing-upi-error":
				w.WriteHeader(http.StatusNotFound)
				io.WriteString(w, `{"timestamp":"2019-07-25T06:34:50.211+0000","status":404,"error":"Not Found","message":"Identity not found","path":"/identity/`+uid+`"}`)
			default:
				w.WriteHeader(http.StatusBadRequest)
				io.WriteString(w, `{"timestamp":"2019-07-25T02:23:32.668+0000","status":400,"error":"Bad Request","message":"Incorrect or not supported id","path":"/identity/`+uid+`"}`)
			}
		case strings.HasPrefix(ru, "/service/employment/integrations/v1/employee/"):
			var empID = strings.TrimPrefix(ru, "/service/employment/integrations/v1/employee/")
			if empID == "477579437" {
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
			}

		default:
			w.WriteHeader(http.StatusNotFound)
			io.WriteString(w, fmt.Sprintf("%q NOT FOUND!", ru))
		}
	}))
	api.BaseURL = server.URL + "/service"
	oh.BaseURL = server.URL
}

func TeardownTest(t *testing.T) {
	if server != nil {
		server.Close()
	}
}

func TestHandler(t *testing.T) {
	SetupTest(t, true)
	defer TeardownTest(t)

	res, err := HandleRequest(
		lambdacontext.NewContext(context.Background(), &lambdacontext.LambdaContext{}),
		Event{Subject: 1234})

	assert.IsType(t, nil, err)
	assert.Equal(
		t,
		`Recieved: main.Event{EPPN:"", Email:"", ORCID:"", Subject:1234, Type:"", URL:"", Records:[]events.SQSMessage(nil)}, Counter: 1`,
		res)
}

func TestIdentityAPICient(t *testing.T) {
	var c Client
	c.ApiKey = os.Getenv("API_KEY")
	c.BaseURL = "https://api.dev.auckland.ac.nz/service/identity/integrations/v3/identity"
	var id Identity
	c.Get("rcir178", &id)
	t.Logf("IDENTITY: %#v", id)
	assert.NotEqual(t, 0, id.ID)
	err := c.Get("rad42", &id)
	if err != nil {
		t.Error(err)
	}
	output, _ := json.MarshalIndent(id, "", "    ")
	t.Logf("ID: %s", string(output))
}

func TestEmploymentAPICient(t *testing.T) {
	var c Client
	c.ApiKey = os.Getenv("API_KEY")
	c.BaseURL = "https://api.dev.auckland.ac.nz/service/employment/integrations/v1/employee"
	var emp Employment
	// c.Get("rcir178", &id)
	// t.Logf("IDENTITY: %#v", id)
	err := c.Get("rcir178", &emp)
	if err != nil {
		t.Error(err)
	}
	output, _ := json.MarshalIndent(emp, "", "    ")
	t.Logf("EMPLOYMENT: %s", string(output))
}

func TestAccessToken(t *testing.T) {
	var c Client
	c.ClientID = os.Getenv("CLIENT_ID")
	c.ClientSecret = os.Getenv("CLIENT_SECRET")
	c.BaseURL = "http://127.0.0.1:5000"
	err := c.GetAccessToken("oauth/token")
	assert.Nil(t, err)
	assert.NotEmpty(t, c.AccessToken)
}

func TestGetOrcidToken(t *testing.T) {
	var c Client
	c.ClientID = os.Getenv("CLIENT_ID")
	c.ClientSecret = os.Getenv("CLIENT_SECRET")
	// c.BaseURL = "http://127.0.0.1:5000"
	c.GetAccessToken("oauth/token")
	var tokens []struct {
		AccessToken  string `json:"access_token"`
		ExpiresIn    int64  `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		Scopes       string `json:"scopes"`
	}
	err := c.Get("api/v1/tokens/rad42%40mailinator.com", &tokens)
	assert.Nil(t, err)
	assert.NotEmpty(t, tokens)
}

func TestProcessRegistration(t *testing.T) {
	var (
		e      Event
		err    error
		output string
	)

	e = Event{EPPN: "rpaw053@auckland.ac.nz", ORCID: "0000-0003-1255-9023"}
	output, err = e.process()
	assert.NotEmpty(t, output)
	assert.Nil(t, err)

	e = Event{EPPN: "rcir178@auckland.ac.nz", ORCID: "0000-0001-8228-7153"}
	output, err = e.process()
	assert.NotEmpty(t, output)
	assert.Nil(t, err)

	e.EPPN = "non-existing-upi-error@error.edu"
	output, err = e.process()
	assert.Empty(t, output)
	assert.NotNil(t, err)
}

func TestHealthCheck(t *testing.T) {
	var e = Event{Type: "PING"}
	output, err := e.process()
	assert.NotEmpty(t, output)
	assert.Equal(t, "GNIP", output)
	assert.Nil(t, err)

	e = Event{Type: "ABCD1234"}
	output, err = e.process()
	assert.Empty(t, output)
	assert.NotNil(t, err)
}

func TestIdentityGetORCID(t *testing.T) {
	var id Identity
	json.Unmarshal([]byte(`{
   "emailAddress":"rosh1234@auckland.ac.nz",
   "emails":[
      {
         "email":"rosh83458349@auckland.ac.nz",
         "lastUpdated":"2017-01-13T17:12:23.000+0000",
         "typeId":"Campus",
         "type":"University",
         "verified":false
      },
      {
         "email":"rpfkjds@aucklanduni.ac.nz",
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
         "id":"http://orcid.org/1234-1234-1234-ABCD",
         "type":"ORCID"
      },
      {
         "id":"2490528",
         "type":"UID"
      }
   ],
   "firstName":"Roshan Prakash",
   "id":208013283
   }`), &id)
	assert.Equal(t, "1234-1234-1234-ABCD", id.GetORCID())
}

func TestIdentityGetOrcidAccessToken(t *testing.T) {
	SetupTest(t, true)
	defer TeardownTest(t)

	oh.ClientID = os.Getenv("CLIENT_ID")
	oh.ClientSecret = os.Getenv("CLIENT_SECRET")
	// oh.BaseURL = "http://127.0.0.1:5000"
	err := oh.GetAccessToken("oauth/token")
	if err != nil {
		t.Error(err)
	}
	var id Identity
	json.Unmarshal([]byte(`{
		"emailAddress":"rcir178NOWAY@auckland.ac.nz",
		"emails":[
			{
				"email":"rad42ABC@mailinator.com",
				"lastUpdated":"2017-01-13T17:12:23.000+0000",
				"typeId":"Campus",
				"type":"University",
				"verified":false
			}
		],
		"extIds":[
			{
				"id":"http://orcid.org/0000-0001-8228-7153",
				"type":"*ORCID*"
			},
			{
				"id":"2490528",
				"type":"UID"
			}
		],
		"upi":"rcir178ABC"
   }`), &id)
	token, ok := id.GetOrcidAccessToken()
	assert.False(t, ok)
	_ = token

	id.Emails[0].Email = "rad42@mailinator.com"
	token, ok = id.GetOrcidAccessToken()
	assert.True(t, ok)
	assert.Equal(t, "ecf16b31-ad54-4ba2-ae55-e97fb90e211a", token.AccessToken)

	id.EmailAddress = "rcir178@auckland.ac.nz"
	token, ok = id.GetOrcidAccessToken()
	assert.True(t, ok)
	assert.Equal(t, "ecf16b31-ad54-4ba2-ae55-e97fb90e211a", token.AccessToken)

	id.Upi = "rcir178"
	token, ok = id.GetOrcidAccessToken()
	assert.True(t, ok)
	assert.Equal(t, "ecf16b31-ad54-4ba2-ae55-e97fb90e211a", token.AccessToken)

	id.ExtIds[0].Type = "ORCID"
	token, ok = id.GetOrcidAccessToken()
	assert.True(t, ok)
	assert.Equal(t, "ecf16b31-ad54-4ba2-ae55-e97fb90e211a", token.AccessToken)
}

// ,
func TestEmpUpdate(t *testing.T) {
	var err error

	_, err = HandleRequest(
		lambdacontext.NewContext(context.Background(), &lambdacontext.LambdaContext{}),
		Event{Subject: 208013283})
	// t.Log(err)
	assert.NotNil(t, err)

	_, err = HandleRequest(
		lambdacontext.NewContext(context.Background(), &lambdacontext.LambdaContext{}),
		Event{Subject: 484378182})
	assert.Nil(t, err)

	_, err = HandleRequest(
		lambdacontext.NewContext(context.Background(), &lambdacontext.LambdaContext{}),
		Event{
			Records: []events.SQSMessage{
				events.SQSMessage{Body: `{"subject":484378182}`},
				events.SQSMessage{Body: `{"subject":208013283}`},
			},
		})
	assert.NotNil(t, err)
}

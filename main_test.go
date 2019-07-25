package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

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
		case ru == "/api/v1/tokens/rad42%40mailinator.com":
			io.WriteString(w, `[{
				"access_token": "kdfjsb31-ad54-4ba2-ae55-e97fb90e211a", 
				"expires_in": 631138518, 
				"issue_time": "2019-07-18T03:13:35", 
				"refresh_token": "kdfjsa20-31be-442a-9faa-73f1d92fac45", 
				"scopes": "/read-limited,/activities/update"
			}]`)
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
			if uid == "1234" {
				io.WriteString(w, `{"timestamp":"2019-07-25T02:23:32.668+0000","status":400,"error":"Bad Request","message":"Incorrect or not supported id","path":"/identity/`+uid+`"}`)
			} else if uid == "rpaw053" {
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
			t.Log("!!!!!! Defaulting on :", ru)
			io.WriteString(w, "Status OK")
		}
	}))
	api.BaseURL = server.URL + "/service"
	oh.BaseURL = server.URL
}

func teardown(t *testing.T) {
	if server != nil {
		server.Close()
	}
}

func TestHandler(t *testing.T) {
	// res, err := HandleRequest(context.Background(), Event{Subject: 1234})
	SetupTest(t, true)
	defer teardown(t)

	res, err := HandleRequest(
		lambdacontext.NewContext(context.Background(), &lambdacontext.LambdaContext{}),
		Event{Subject: 1234})

	assert.IsType(t, nil, err)
	t.Logf("****%#v\n", err)
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
	c.ClientID = os.Getenv("ORCIDHUB_CLIENT_ID")
	c.ClientSecret = os.Getenv("ORCIDHUB_CLIENT_SECRET")
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
	var e = Event{EPPN: "rpaw053@auckland.ac.nz"}
	output, err := e.process()
	assert.NotEmpty(t, output)
	assert.Nil(t, err)

	e.EPPN = "non-existing-upi-error@error.edu"
	output, err = e.process()
	assert.Empty(t, output)
	assert.NotNil(t, err)
}

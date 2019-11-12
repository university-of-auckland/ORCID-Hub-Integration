//+build test

package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	server                                            *httptest.Server
	malformatResponse, withTasks, withAnIncomleteTask bool
	live                                              bool
)

func TestMain(m *testing.M) {
	flag.BoolVar(&verbose, "verbose", false, "Print out the received responses.")
	flag.BoolVar(&live, "live", false, "Run with the DEV/SANDBOX APIs.")
	flag.Parse()

	taskRetentionMin = 1
	batchSize = 1

	os.Exit(m.Run())
}

func getenv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func isValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func setupTests(t *testing.T) {

	server = httptest.NewServer(createMockHandler(t))

	APIBaseURL = server.URL + "/service"
	OHBaseURL = server.URL
	api.baseURL = server.URL + "/service"
	oh.baseURL = server.URL

	// warm-up the server
	for {
		time.Sleep(time.Millisecond * 10)
		resp, err := http.Get(server.URL + "/ping")
		if err == nil && resp.StatusCode == http.StatusNoContent {
			break
		}
	}
}

func teardownTests(t *testing.T) {
	if server != nil {
		server.Close()
	}
}

func TestCore(t *testing.T) {

	withAnIncomleteTask = true

	if !live {
		setupTests(t)
		defer teardownTests(t)
	} else {
		setupAPIClients()
	}

	t.Run("TaskControl", testTaskControl)
	t.Run("Handler", testHandler)
	t.Run("GetOrcidToken", testGetOrcidToken)
	t.Run("IdentityGetOrcidAccessToken", testIdentityGetOrcidAccessToken)
	t.Run("AccessToken", testAccessToken)
	t.Run("IdentityAPICient", testIdentityAPICient)
	t.Run("StudentAPICient", testStudentAPICient)
	t.Run("EmploymentAPICient", testEmploymentAPICient)
	t.Run("ProcessRegistration", testProcessRegistration)
	t.Run("ProcessEmpUpdate", testProcessEmpUpdate)
	t.Run("ProcessMixed", testProcessMixed)
	t.Run("HealthCheck", testHealthCheck)
	t.Run("MalformatedPayload", testMalformatedPayload)
}

func testTaskControl(t *testing.T) {

	malformatResponse = false
	counter = 0
	(&Event{Type: "PING"}).handle()

	taskRecordCount = 999
	taskCreatedAt.Add(time.Hour)
	(&Event{Type: "PING"}).handle()

	taskCreatedAt.Add(-2 * time.Hour)
	(&Event{Type: "PING"}).handle()

	assert.Equal(t, 3, counter)

	for _, o := range []struct {
		v1 bool
		v2 bool
	}{
		{true, true},
		{true, false},
		{false, true},
		{false, false},
	} {
		taskID = 0
		withTasks = o.v1
		withAnIncomleteTask = o.v2

		(&Event{Type: "PING"}).handle()
		assert.NotEqual(t, 0, taskID)
	}
	assert.Equal(t, 7, counter)
}

func testMalformatedPayload(t *testing.T) {
	if live {
		t.Skip()
	}

	var fatalCallCount int
	logFatal = func(args ...interface{}) { fatalCallCount++; t.Log("*** FATAL: ", args) }
	malformatResponse = true

	(&Task{ID: 123456}).activate()
	newTask()

	malformatResponse = false
	logFatal = log.Fatal
	assert.Equal(t, 1, fatalCallCount)
}

func testHandler(t *testing.T) {

	malformatResponse = false

	_, err := (&Event{Subject: 1233}).handle()
	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "failed to retrieve the identity record")

	_, err = (&Event{Subject: 8524255}).handle()
	if !live {
		require.NotNil(t, err)
		assert.Contains(t, err.Error(), "hasn't granted access to the profile")
	}

	_, err = (&Event{Subject: 123}).handle()
	if !live {
		require.NotNil(t, err)
		assert.Contains(t, err.Error(), "failed to retrieve the identity record")
	}

	_, err = (&Event{Subject: 1234567890123}).handle()
	if !live {
		require.NotNil(t, err)
		assert.Contains(t, err.Error(), "failed to retrieve the identity record")
	}

	_, err = (&Event{Type: "ERROR"}).handle()
	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "unhandled")
}

func testIdentityAPICient(t *testing.T) {
	var c Client
	// c.ApiKey = os.Getenv("APIKEY")
	// c.BaseURL = "https://api.dev.auckland.ac.nz"

	c.baseURL = APIBaseURL
	var id Identity
	c.get("identity/integrations/v3/identity/rcir178", &id)
	assert.NotEqual(t, -1, id.ID)

	err := c.post("identity/integrations/v3/identity/rcir178", `{"test": 1234}`, &id)
	assert.Nil(t, err)

	err = c.post("identity/integrations/v3/identity/rcir178", t.Log, &id)
	assert.NotNil(t, err)

	var idNotFound Identity
	err = c.get("identity/integrations/v3/identity/rad42", &idNotFound)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 0, idNotFound.ID)

	malformatResponse = true
	id.ID = 0
	c.get("identity/integrations/v3/identity/rcir178", &id)
	assert.Equal(t, 0, id.ID)
	malformatResponse = false

	err = c.do("POST", "identity/integrations/v3/identity/rcir178", nil, &id)
	assert.Nil(t, err)
	malformatResponse = false

}

func testStudentAPICient(t *testing.T) {
	malformatResponse = false
	var c Client
	// c.ApiKey = os.Getenv("APIKEY")
	// c.BaseURL = "https://api.dev.auckland.ac.nz"

	c.baseURL = APIBaseURL
	var degrees Degrees

	c.get("student/integrations/v1/student/208013283/degree/", &degrees)
	assert.Equal(t, 2, len(degrees))

	c.get("student/integrations/v1/student/477579437/degree/", &degrees)
	assert.Equal(t, 1, len(degrees))

	c.get("student/integrations/v1/student/8524255/degree/", &degrees)
	assert.Equal(t, 2, len(degrees))

	c.get("student/integrations/v1/student/484378182/degree/", &degrees)
	assert.Equal(t, 0, len(degrees))

	c.get("student/integrations/v1/student/9999999/degree/", &degrees)
	assert.Equal(t, 0, len(degrees))

	// malformated message:
	c.get("student/integrations/v1/student/208013283/degree/", &degrees)
	malformatResponse = true
	_, err := degrees.propagateToHub("rpaw058@auckland.ac.nz", "0000-0003-1255-9023")
	assert.NotNil(t, err)
	malformatResponse = false

}

func testEmploymentAPICient(t *testing.T) {
	var c Client

	malformatResponse = false
	c.apiKey = os.Getenv("APIKEY")
	c.baseURL = APIBaseURL

	var emp Employment
	err := c.get("employment/integrations/v1/employee/rcir178", &emp)
	if err != nil {
		t.Error(err)
	}

	count, err := emp.propagateToHub("rcir178@auckland.ac.nz", "0000-0001-8228-7153")
	assert.NotZero(t, count)
	assert.Nil(t, err)

	// malformated message:
	malformatResponse = true
	count, err = emp.propagateToHub("rcir178@auckland.ac.nz", "0000-0001-8228-7153")
	assert.Equal(t, 1, count)
	assert.NotNil(t, err)
	malformatResponse = false

	// no jobs
	emp.Job = nil
	count, err = emp.propagateToHub("rcir178@auckland.ac.nz", "0000-0001-8228-7153")
	assert.Zero(t, count)
	assert.NotNil(t, err)
}

func testAccessToken(t *testing.T) {

	var c Client
	malformatResponse = false
	c.clientID = os.Getenv("CLIENT_ID")
	c.clientSecret = os.Getenv("CLIENT_SECRET")
	c.baseURL = OHBaseURL
	err := c.getAccessToken("oauth/token")
	assert.Nil(t, err)
	assert.NotEmpty(t, c.accessToken)

	// malformated message
	logFatal = func(args ...interface{}) {}
	malformatResponse = true

	c.accessToken = ""
	err = c.getAccessToken("oauth/token")
	assert.NotNil(t, err)
	assert.Empty(t, c.accessToken)

	at := oh.accessToken
	oh.accessToken = ""
	setupAPIClients()
	oh.accessToken = at

	malformatResponse = false
	logFatal = log.Fatal
}

func testGetOrcidToken(t *testing.T) {

	var c Client
	malformatResponse = false
	c.clientID = os.Getenv("CLIENT_ID")
	c.clientSecret = os.Getenv("CLIENT_SECRET")
	c.baseURL = OHBaseURL
	c.getAccessToken("oauth/token")
	var tokens []struct {
		AccessToken  string `json:"access_token"`
		ExpiresIn    int64  `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		Scopes       string `json:"scopes"`
	}
	err := c.get("api/v1/tokens/rad42@mailinator.com", &tokens)
	assert.Nil(t, err)
	assert.NotEmpty(t, tokens)

	malformatResponse = true
	tokens = nil
	err = c.get("api/v1/tokens/rad42@mailinator.com", &tokens)
	assert.NotNil(t, err)
	assert.Empty(t, tokens)
	malformatResponse = false
}

func testProcessRegistration(t *testing.T) {
	var (
		e      Event
		err    error
		output string
	)

	taskID = 0
	taskRecordCount = 0
	malformatResponse = false

	setupAPIClients()
	if live {
		// Remove the existing ORCID iDs
		for _, upi := range []string{"rpaw053", "rcir178", "djim087"} {
			api.do("DELETE", "identity/integrations/v3/identity/"+upi+"/identifier/ORCID", nil, nil)
		}
	}

	withAnIncomleteTask = true

	e = Event{Type: "CREATED", EPPN: "rpaw053@auckland.ac.nz", ORCID: "0000-0003-1255-9023"}
	output, err = e.handle()
	assert.NotEmpty(t, output)
	assert.Nil(t, err)

	e = Event{Type: "CREATED", EPPN: "rcir178@auckland.ac.nz", ORCID: "0000-0001-8228-7153"}
	output, err = e.handle()
	assert.NotEmpty(t, output)
	assert.Nil(t, err)

	e = Event{Type: "CREATED", EPPN: "djim087@auckland.ac.nz", ORCID: "0000-0002-3008-0422"}
	output, err = e.handle()
	assert.NotEmpty(t, output)
	assert.Nil(t, err)

	withAnIncomleteTask = false
	taskID = 0

	e.EPPN = "non-existing-upi-error@error.edu"
	output, err = e.handle()
	assert.Empty(t, output)
	assert.NotNil(t, err)

	// malformatted messages:
	wg.Wait()
	logFatal = func(args ...interface{}) {}
	malformatResponse = true

	e = Event{Type: "CREATED", EPPN: "djim087@auckland.ac.nz", ORCID: "0000-0002-3008-0422"}
	output, err = e.handle()
	assert.Empty(t, output)
	assert.NotNil(t, err)

	wg.Wait()
	malformatResponse = false
	// logFatal = log.Fatal
}

func testHealthCheck(t *testing.T) {

	malformatResponse = false
	withAnIncomleteTask = false
	taskID = 0

	var e = Event{Type: "PING"}
	output, err := e.handle()
	assert.NotEmpty(t, output)
	assert.Equal(t, "GNIP", output)
	assert.Nil(t, err)

	e = Event{Type: "ABCD1234"}
	output, err = e.handle()
	assert.Empty(t, output)
	assert.NotNil(t, err)
}

func TestIdentityGetORCID(t *testing.T) {
	var id Identity
	json.Unmarshal([]byte(`{
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
   ]}`), &id)
	assert.Equal(t, "1234-1234-1234-ABCD", id.GetORCID())

	json.Unmarshal([]byte(`{
   "extIds":[
      {
         "id":"2121820801328312",
         "type":"IDCard"
      }
   ]}`), &id)
	assert.Equal(t, "", id.GetORCID())

	json.Unmarshal([]byte(`{}`), &id)
	assert.Equal(t, "", id.GetORCID())
}

func testIdentityGetOrcidAccessToken(t *testing.T) {

	malformatResponse = false
	err := oh.getAccessToken("oauth/token")
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
		"id":123443,
		"upi":"rcir178ABC"
   }`), &id)
	token, ok := id.GetOrcidAccessToken()
	assert.False(t, ok)
	_ = token

	id.Emails[0].Email = "rad42@mailinator.com"
	token, ok = id.GetOrcidAccessToken()
	assert.True(t, ok)
	assert.True(t, isValidUUID(token.AccessToken))
	if !live {
		assert.Equal(t, "ecf16b31-ad54-4ba2-ae55-e97fb90e211a", token.AccessToken)
	}

	id.EmailAddress = "rcir178@auckland.ac.nz"
	token, ok = id.GetOrcidAccessToken()
	assert.True(t, ok)
	assert.True(t, isValidUUID(token.AccessToken))
	if !live {
		assert.Equal(t, "ecf16b31-ad54-4ba2-ae55-e97fb90e211a", token.AccessToken)
	}

	id.Upi = "rcir178"
	token, ok = id.GetOrcidAccessToken()
	assert.True(t, ok)
	assert.True(t, isValidUUID(token.AccessToken))
	if !live {
		assert.Equal(t, "ecf16b31-ad54-4ba2-ae55-e97fb90e211a", token.AccessToken)
	}

	id.ExtIds[0].Type = "ORCID"
	token, ok = id.GetOrcidAccessToken()
	assert.True(t, ok)
	assert.True(t, isValidUUID(token.AccessToken))
	if !live {
		assert.Equal(t, "ecf16b31-ad54-4ba2-ae55-e97fb90e211a", token.AccessToken)
	}

	// no update scope
	id.Upi = "dthn666"
	id.ExtIds = nil
	token, ok = id.GetOrcidAccessToken()
	assert.False(t, ok)

	// malformated message
	malformatResponse = true
	token, ok = id.GetOrcidAccessToken()
	assert.False(t, ok)
	malformatResponse = false
}

func testProcessEmpUpdate(t *testing.T) {

	var err error
	taskRecordCount = 0
	taskID = 0
	malformatResponse = false
	withAnIncomleteTask = true

	(&Event{Subject: 208013283}).handle()
	assert.Nil(t, err)
	if !live {
		assert.Equal(t, 6, taskRecordCount)
	}

	_, err = (&Event{Subject: 484378182}).handle()
	assert.Nil(t, err)

	taskRecordCount = 0
	taskID = 0
	_, err = (&Event{
		Records: []events.SQSMessage{
			{Body: `{"subject":"484378182"}`},
			{Body: `{"subject":"477579437"}`},
			{Body: `{"subject":"208013283"}`},
			{Body: `{"subject":"987654321"}`},
			{Body: `{"subject":"8524255"}`},
			{Body: `{"subject":"350622514"}`},
			{Body: `{"subject":"4306445"}`},
		},
	}).handle()
	assert.True(t, taskRecordCount > 0, "The number of records should be > 0.")
	t.Log(err)
	assert.NotNil(t, err)

	// Malformatted

	logFatal = func(args ...interface{}) {}
	malformatResponse = true
	_, err = (&Event{
		Records: []events.SQSMessage{
			{Body: `{"subject":"484378182"}`},
			{Body: `{"subject":"208013283"}`},
			{Body: `{"subject":"4306445"}`},
		},
	}).handle()
	malformatResponse = false
	logFatal = log.Fatal

}

func testProcessMixed(t *testing.T) {

	var err error

	taskRecordCount = 0
	taskID = 0
	withAnIncomleteTask = true
	malformatResponse = false

	_, err = (&Event{
		Records: []events.SQSMessage{
			{Body: `{"subject":"484378182"}`},
			{Body: `{"subject":"477579437"}`},
			{Body: `{
				"orcid": "0000-0001-8228-7153", 
				"url": "https://sandbox.orcid.org/0000-0001-8228-7153", 
				"type": "CREATED", "updated-at": "2019-07-25T02:05:32", 
				"email": "rad42@mailinator.com", 
				"eppn": "rcir178@auckland.ac.nz"
			}`},
			{Body: `{"subject":"208013283"}`},
			{Body: `{"subject":"66666666"}`},
			{Body: `{"subject":"77777777"}`},
			{Body: `{
				"orcid": "0000-0001-6666-7153", 
				"url": "https://sandbox.orcid.org/0000-0001-6666-7153", 
				"type": "CREATED", 
				"updated-at": "2019-07-25T02:05:32", 
				"email": "dthn666@mailinator.com", 
				"eppn": "dthn666@auckland.ac.nz"
			}`},
			{Body: `{
				"orcid": "0000-0001-7777-7153", 
				"url": "https://sandbox.orcid.org/0000-0001-7777-7153", 
				"type": "CREATED", 
				"updated-at": "2019-07-25T02:05:32", 
				"email": "dthn7777mailinator.com", 
				"eppn": "dthn7777auckland.ac.nz"
			}`},
			{Body: `{"subject":"987654321"}`},
			{Body: `{"subject":"8524255"}`},
			{Body: `{
				"orcid": "0000-0001-8228-7153", 
				"url": "https://sandbox.orcid.org/0000-0001-8228-7153", 
				"type": "UPDATED", 
				"updated-at": "2019-07-25T02:05:32", 
				"email": "rad42@mailinator.com", 
				"eppn": "rcir178@auckland.ac.nz"
			}`},
			{Body: `{"subject":"350622514"}`},
			{Body: `{"subject":"4306445"}`},
		},
	}).handle()

	if !live {
		assert.True(t, taskRecordCount == 8, "The number of records should be 8, got: %d.", taskRecordCount)
	}
	assert.NotNil(t, err)

	counter = 0
	_, err = (&Event{
		Records: []events.SQSMessage{
			{Body: `{"subject":"484378182"}`},
			{Body: `{"subject":"477579437"}`},
		},
	}).handle()
	assert.NotNil(t, err)
	assert.Equal(t, 3, counter)

	logFatal = func(args ...interface{}) {}
	malformatResponse = true
	counter = 0
	assert.NotNil(t, err)
	_, err = (&Event{
		Records: []events.SQSMessage{
			{Body: `{"unknown":ABC484378182}`},
			{Body: `{"unknown":ABC477579437}`},
		},
	}).handle()
	assert.Nil(t, err)
	assert.Equal(t, 1, counter)
	malformatResponse = false
	logFatal = log.Fatal
}

func TestIsValidUPIAndID(t *testing.T) {
	assert.True(t, isValidUPI("rcir178"))
	assert.True(t, isValidUPI("rpaw053"))
	assert.False(t, isValidUPI("123456"))
	assert.False(t, isValidUPI("ABC123456"))
	assert.False(t, isValidUPI("abc1234"))
	assert.False(t, isValidUPI("abcdd34"))
	assert.False(t, isValidUPI("abcd23x"))
}

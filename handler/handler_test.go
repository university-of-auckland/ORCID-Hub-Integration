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
	server                         *httptest.Server
	withTasks, withAnIncomleteTask bool
	live                           bool
)

func init() {
	flag.BoolVar(&verbose, "verbose", false, "Print out the received responses.")
	flag.BoolVar(&live, "live", false, "Run with the DEV/SANDBOX APIs.")
	flag.Parse()
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
		defer func() {
			gotAccessTokenWG.Wait()
			taskSetUpWG.Wait()
			wg.Wait()

			teardownTests(t)
		}()
	} else {
		setupAPIClients()
	}

	t.Run("TaskControl", testTaskControl)
	t.Run("Handler", testHandler)
	t.Run("GetOrcidToken", testGetOrcidToken)
	t.Run("IdentityGetOrcidAccessToken", testIdentityGetOrcidAccessToken)
	t.Run("AccessToken", testAccessToken)
	t.Run("IdentityAPICient", testIdentityAPICient)
	t.Run("EmploymentAPICient", testEmploymentAPICient)
	t.Run("ProcessRegistration", testProcessRegistration)
	t.Run("ProcessEmpUpdate", testProcessEmpUpdate)
	t.Run("ProcessMixed", testProcessMixed)
	t.Run("HealthCheck", testHealthCheck)
}

func testTaskControl(t *testing.T) {

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
		{false, false},
		{false, true},
		{true, false},
		{true, true},
	} {
		taskSetUpWG.Wait()
		taskID = 0
		withTasks = o.v1
		withAnIncomleteTask = o.v2

		(&Event{Type: "PING"}).handle()
		assert.NotEqual(t, 0, taskID)
	}
	assert.Equal(t, 7, counter)

}

func testHandler(t *testing.T) {

	_, err := (&Event{Subject: 1234}).handle()
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
	// c.ApiKey = os.Getenv("API_KEY")
	// c.BaseURL = "https://api.dev.auckland.ac.nz"

	c.baseURL = APIBaseURL
	var id Identity
	c.get("identity/integrations/v3/identity/rcir178", &id)
	assert.NotEqual(t, 0, id.ID)

	var idNotFound Identity
	err := c.get("identity/integrations/v3/identity/rad42", &idNotFound)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 0, idNotFound.ID)
}

func testEmploymentAPICient(t *testing.T) {
	var c Client

	c.apiKey = os.Getenv("API_KEY")
	c.baseURL = APIBaseURL

	var emp Employment
	err := c.get("rcir178", &emp)
	if err != nil {
		t.Error(err)
	}
}

func testAccessToken(t *testing.T) {

	var c Client
	c.clientID = os.Getenv("CLIENT_ID")
	c.clientSecret = os.Getenv("CLIENT_SECRET")
	c.baseURL = OHBaseURL
	err := c.getAccessToken("oauth/token")
	assert.Nil(t, err)
	assert.NotEmpty(t, c.accessToken)
}

func testGetOrcidToken(t *testing.T) {

	var c Client
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
}

func testProcessRegistration(t *testing.T) {
	var (
		e      Event
		err    error
		output string
	)

	setupAPIClients()
	gotAccessTokenWG.Wait()
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
}

func testHealthCheck(t *testing.T) {

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
}

func testIdentityGetOrcidAccessToken(t *testing.T) {

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
}

func testProcessEmpUpdate(t *testing.T) {

	var err error

	taskRecordCount = 0
	_, err = (&Event{Subject: 208013283}).handle()
	// t.Log(err)
	assert.NotNil(t, err)
	assert.Equal(t, 0, taskRecordCount)

	_, err = (&Event{Subject: 484378182}).handle()
	assert.Nil(t, err)

	taskRecordCount = 0
	_, err = (&Event{
		Records: []events.SQSMessage{
			{Body: `{"subject":484378182}`},
			{Body: `{"subject":477579437}`},
			{Body: `{"subject":208013283}`},
			{Body: `{"subject":987654321}`},
			{Body: `{"subject":8524255}`},
			{Body: `{"subject":350622514}`},
			{Body: `{"subject":4306445}`},
		},
	}).handle()
	assert.True(t, taskRecordCount > 0, "The number of records should be > 0.")
	t.Log(err)
	assert.NotNil(t, err)
}

func testProcessMixed(t *testing.T) {

	var err error

	taskRecordCount = 0
	_, err = (&Event{
		Records: []events.SQSMessage{
			{Body: `{"subject":484378182}`},
			{Body: `{"subject":477579437}`},
			{Body: `{
				"orcid": "0000-0001-8228-7153", 
				"url": "https://sandbox.orcid.org/0000-0001-8228-7153", 
				"type": "CREATED", "updated-at": "2019-07-25T02:05:32", 
				"email": "rad42@mailinator.com", 
				"eppn": "rcir178@auckland.ac.nz"
			}`},
			{Body: `{"subject":208013283}`},
			{Body: `{"subject":66666666}`},
			{Body: `{"subject":77777777}`},
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
			{Body: `{"subject":987654321}`},
			{Body: `{"subject":8524255}`},
			{Body: `{
				"orcid": "0000-0001-8228-7153", 
				"url": "https://sandbox.orcid.org/0000-0001-8228-7153", 
				"type": "UPDATED", 
				"updated-at": "2019-07-25T02:05:32", 
				"email": "rad42@mailinator.com", 
				"eppn": "rcir178@auckland.ac.nz"
			}`},
			{Body: `{"subject":350622514}`},
			{Body: `{"subject":4306445}`},
		},
	}).handle()
	assert.True(t, taskRecordCount == 3, "The number of records should be 3, got: %d.", taskRecordCount)
	t.Log(err)
	assert.NotNil(t, err)
}

func TestIsValidUPIAndID(t *testing.T) {
	assert.True(t, isValidUPI("rcir178"))
	assert.True(t, isValidUPI("rpaw053"))
	assert.False(t, isValidUPI("123456"))
	assert.False(t, isValidUPI("ABC123456"))
	assert.False(t, isValidUPI("abc1234"))
	assert.False(t, isValidUPI("abcdd34"))
	assert.False(t, isValidUPI("abcd23x"))

	assert.False(t, isValidID("123456789A"))
	assert.False(t, isValidID("123"))
	assert.False(t, isValidID("1234567890123445"))
	assert.True(t, isValidID("12345678"))
	assert.True(t, isValidID("123456789"))
	assert.True(t, isValidID("1234567890"))
}

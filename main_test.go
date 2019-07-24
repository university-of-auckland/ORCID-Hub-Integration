package main

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	// res, err := HandleRequest(context.Background(), Event{Subject: 1234})
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
	c.BaseURL = "http://127.0.0.1:5000"
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

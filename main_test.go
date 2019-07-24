package main

import (
	"context"
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
	c.BaseURL = "https://api.dev.auckland.ac.nz/service/identity/integrations/v2/identity"
	var id Identity
	// c.Get("rcir178", &id)
	// t.Logf("IDENTITY: %#v", id)
	err := c.Get("rad42@mailinator.com", &id)
	if err != nil {
		t.Error(err)
	}
	t.Logf("IDENTITY: %#v", id)
}

func TestEmploymentAPICient(t *testing.T) {
	var c Client
	c.ApiKey = os.Getenv("API_KEY")
	c.BaseURL = "https://api.dev.auckland.ac.nz/service/identity/integrations/v2/identity"
	var id Identity
	// c.Get("rcir178", &id)
	// t.Logf("IDENTITY: %#v", id)
	err := c.Get("rad42@mailinator.com", &id)
	if err != nil {
		t.Error(err)
	}
	t.Logf("IDENTITY: %#v", id)
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

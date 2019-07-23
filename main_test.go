package main

import (
	"context"
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
		`Recieved: main.Event{EPPN:"", Email:"", ORCID:"", Subject:1234, Type:"", URL:"", Records:[]events.SQSMessage(nil)}`,
		res)
}

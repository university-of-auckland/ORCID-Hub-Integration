package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

type Platform interface {
}

// Generic message suitable for both EMP Update event and
// ORCIDHub Webhook propagagted event (with wrapped in SQS message batch):
type Event struct {
	EPPN    string `json:"eppn"`
	Email   string `json:"email"`
	ORCID   string `json:"orcid"`
	Subject int    `json:"subject"`
	Type    string `json:"type"`
	URL     string `json:"url"`
	// SQS Message if used SQS
	Records []events.SQSMessage
}

var (
	api     Client
	counter int
)

func processEmpUpdate(e Event) (string, error) {
	return "", nil
}

func processUserRegistration(e Event) (string, error) {
	if api.ApiKey == "" {
		api.ApiKey = os.Getenv("API_KEY")
		api.BaseURL = "https://api.dev.auckland.ac.nz/service/identity/integrations/v2/identity"
	}
	var id Identity
	parts := strings.Split(e.EPPN, "@")
	log.Println("UID: ", parts[0])
	err := api.Get(parts[0], &id)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%#v", id), nil
}

func HandleRequest(ctx context.Context, e Event) (string, error) {

	log.Printf("Cotext: %#v, counter: %d", ctx, counter)
	counter += 1
	lc, _ := lambdacontext.FromContext(ctx)
	log.Printf("Lambda Context: %#v", lc)
	// log.Print(lc.Identity.CognitoIdentityPoolID)
	// isLambda := os.Getenv("_LAMBDA_SERVER_PORT") != ""
	// for _, pair := range os.Environ() {
	// 	log.Println(pair)
	// }
	log.Printf("Recieved: %#v", e)
	if e.Records != nil {
		for _, r := range e.Records {
			var e Event
			json.Unmarshal([]byte(r.Body), &e)
			log.Printf("Event: %#v", e)
		}
	} else if e.EPPN != "" {
		return processUserRegistration(e)
	}

	return fmt.Sprintf("Recieved: %#v, Counter: %d", e, counter), nil

	// var e Event
	// err := json.Unmarshal(message, &e)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// if e.Subject != 0 {
	// 	processEmpUpdate(e)
	// } else if e.Type == "CREATED" {
	// 	processUserRegistration(e)
	// } else {
	// 	log.Printf("The event %#v discarded.", e)
	// }
	// return fmt.Sprintf("Recieved: %#v", e), nil
}

func main() {
	log.SetPrefix("OHI")
	lambda.Start(HandleRequest)
}

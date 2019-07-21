package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

type Platform interface {
}

// Generic message suitable for both EMP Update event and
// ORCIDHub Webhook propagagted event:
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

func processEmpUpdate(e Event) (string, error) {
	return "", nil
}

func processUserRegistration(e Event) (string, error) {
	return "", nil
}

func HandleRequest(ctx context.Context, e Event) (string, error) {
	log.Printf("Cotext: %#v", ctx)
	lc, _ := lambdacontext.FromContext(ctx)
	log.Printf("Lambda Context: %#v", lc)
	log.Print(lc.Identity.CognitoIdentityPoolID)
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
	}

	return fmt.Sprintf("Recieved: %#v", e), nil

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

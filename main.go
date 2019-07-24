package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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
	oh      Client
	counter int
)

func (e *Event) process() (string, error) {
	if api.ApiKey == "" {
		api.ApiKey = os.Getenv("API_KEY")
		api.BaseURL = "https://api.dev.auckland.ac.nz/service"
	}
	if e.EPPN != "" {
		return e.processUserRegistration()
	} else if e.Subject != 0 {
		return e.processEmpUpdate()
	}
	return "", fmt.Errorf("Unhandled event: %#v", e)
}

func (e *Event) processEmpUpdate() (string, error) {
	if oh.ClientID == "" {
		oh.ClientID = os.Getenv("CLIENT_ID")
		oh.ClientSecret = os.Getenv("CLIENT_SECRET")
		oh.BaseURL = "https://api.dev.auckland.ac.nz/service"
	}

	return "", nil
}

func (e *Event) processUserRegistration() (string, error) {
	var id Identity
	parts := strings.Split(e.EPPN, "@")
	log.Println("UID: ", parts[0])
	err := api.Get("identity/integrations/v3/identity/"+parts[0], &id)
	if err != nil {
		return "", err
	}
	if id.ID == 0 {
		return "", fmt.Errorf("No Identity for %q (%s)", e.EPPN, e.Email)
	}
	var emp Employment
	err = api.Get("employment/integrations/v1/employee/"+strconv.Itoa(id.ID), &emp)
	if err != nil {
		return "", err
	}
	if len(emp.Job) > 0 {
		for _, j := range emp.Job {
			log.Print("JOB: ", j)
		}
	}

	return fmt.Sprintf("%#v", id), nil

}

func HandleRequest(ctx context.Context, e Event) (string, error) {

	log.Printf("Cotext: %#v, counter: %d", ctx, counter)
	counter += 1
	// lc, _ := lambdacontext.FromContext(ctx)
	// log.Print(lc.Identity.CognitoIdentityPoolID)
	// isLambda := os.Getenv("_LAMBDA_SERVER_PORT") != ""
	// for _, pair := range os.Environ() {
	// 	log.Println(pair)
	// }
	if e.Records != nil {
		var resp []string
		for _, r := range e.Records {
			var e Event
			json.Unmarshal([]byte(r.Body), &e)
			r, err := e.process()
			if err != nil {
				return "", err
			}
			resp = append(resp, r)
		}
		return strings.Join(resp, "; "), nil
	}
	return e.process()

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

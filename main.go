package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

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
	api             Client
	oh              Client
	counter         int
	taskID          int
	taskCreatedAt   time.Time
	taskRecordCount int
)

const taskFilenamePrefix = "UOA-OH-INTEGRATION-TASK-"

var (
	APIBaseURL = "https://api.dev.auckland.ac.nz/service"
	OHBaseURL  = "https://dev.orcidhub.org.nz"
)

func setup() {
	if api.ApiKey == "" {
		api.ApiKey = os.Getenv("API_KEY")
		api.BaseURL = APIBaseURL
	}
	if oh.AccessToken == "" {
		oh.ClientID = os.Getenv("CLIENT_ID")
		oh.ClientSecret = os.Getenv("CLIENT_SECRET")
		oh.BaseURL = OHBaseURL
		// oh.BaseURL = "http://127.0.0.1:5000"
		err := oh.GetAccessToken("oauth/token")
		if err != nil {
			panic(err)
		}
	}
	go setupTask()
}

func (e *Event) process() (string, error) {

	if e.EPPN != "" {
		return e.processUserRegistration()
	} else if e.Subject != 0 {
		return e.processEmpUpdate()
	} else if e.Type == "PING" { // Heartbeat Check
		return "GNIP", nil
	}
	return "", fmt.Errorf("Unhandled event: %#v", e)
}

func (e *Event) processEmpUpdate() (string, error) {
	var employeeID = strconv.Itoa(e.Subject)
	var id Identity
	err := api.Get("identity/integrations/v3/identity/"+employeeID, &id)
	if err != nil {
		return "", err
	}
	// TODO: check if the user has ORCID and linked account on the HUB
	// TODO: update ORCID
	// TODO: update employment records

	return "", nil
}

func (e *Event) processUserRegistration() (string, error) {
	var id Identity
	parts := strings.Split(e.EPPN, "@")
	log.Println("UID: ", parts[0])

	setup()
	err := api.Get("identity/integrations/v3/identity/"+parts[0], &id)
	if err != nil {
		return "", err
	}
	if id.ID == 0 {
		return "", fmt.Errorf("No Identity for %q (%s)", e.EPPN, e.Email)
	}
	var employeeID = strconv.Itoa(id.ID)
	for _, eid := range id.ExtIds {
		if eid.Type == "ORCID" {
			parts := strings.Split(eid.ID, "/")
			orcid := parts[len(parts)-1]
			if e.ORCID == "" {
				e.ORCID = orcid
			} // else orcid != e.ORCID { // TODO
			goto HAS_ORCID
		}
	}
	{
		// Add ORCID ID if the user doesn't have one
		var resp struct {
			StatusCode string `json:"statusCode"`
		}
		err = api.Put("identity/integrations/v3/identity/"+employeeID+"/identifier/ORCID", map[string]string{"identifier": e.ORCID}, &resp)
		if err != nil {
			return "", err
		}
	}
HAS_ORCID:
	var emp Employment
	err = api.Get("employment/integrations/v1/employee/"+employeeID, &emp)
	if err != nil {
		return "", err
	}
	// TODO: register employment entries
	if len(emp.Job) > 0 {
		records := make([]Record, len(emp.Job))
		for i, job := range emp.Job {
			records[i] = Record{
				AffiliationType: "employment",
				Department:      job.DepartmentDescription,
				EndDate:         job.JobEndDate,
				ExternalID:      job.PositionNumber,
				Email:           id.EmailAddress,
				Orcid:           e.ORCID,
				Role:            job.PositionDescription,
				StartDate:       job.JobStartDate,
			}
			log.Print("JOB: ", job)
		}
		var task Task
		err := oh.Patch("api/v1/affiliations/"+strconv.Itoa(taskID), Task{ID: taskID, Records: records}, &task)
		if err != nil {
			panic(err)
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

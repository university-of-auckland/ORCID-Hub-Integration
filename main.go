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

	"github.com/aws/aws-lambda-go/lambda"
)

var (
	api             Client
	oh              Client
	counter         int
	gotAccessToken  chan bool
	taskSetUp       chan bool
	taskID          int
	taskCreatedAt   time.Time
	taskRecordCount int
	verbose         bool
)

const taskFilenamePrefix = "UOA-OH-INTEGRATION-TASK-"

var (
	APIBaseURL = "https://api.dev.auckland.ac.nz/service"
	OHBaseURL  = "https://dev.orcidhub.org.nz"
)

func init() {
	gotAccessToken = make(chan bool, 1)
	taskSetUp = make(chan bool, 1)
}

func setup() {
	if api.ApiKey == "" {
		api.ApiKey = os.Getenv("API_KEY")
		api.BaseURL = APIBaseURL
	}
	go func() {
		if oh.AccessToken == "" {
			oh.ClientID = os.Getenv("CLIENT_ID")
			oh.ClientSecret = os.Getenv("CLIENT_SECRET")
			oh.BaseURL = OHBaseURL
			// oh.BaseURL = "http://127.0.0.1:5000"
			err := oh.GetAccessToken("oauth/token")
			if err != nil {
				log.Panic(err)
			}
		}
		gotAccessToken <- true
	}()
	go setupTask()
}

func (e *Event) process() (string, error) {

	if e.EPPN != "" || e.Subject != 0 || e.Type == "PING" {
		setup()
		if e.EPPN != "" {
			return e.processUserRegistration()
		} else if e.Subject != 0 {
			return e.processEmpUpdate()
		} else if e.Type == "PING" { // Heartbeat Check
			return "GNIP", nil
		}
	}
	return "", fmt.Errorf("Unhandled event: %#v", e)
}

func (e *Event) processEmpUpdate() (string, error) {

	var employeeID = strconv.Itoa(e.Subject)
	identities := make(chan Identity, 1)
	employments := make(chan Employment, 1)

	go getIdentidy(identities, employeeID)
	id := <-identities
	token, ok := id.GetOrcidAccessToken()

	if !ok {
		return "", fmt.Errorf("the user (ID: %s) hasn't granted access to the profile", employeeID)
	}

	go getEmp(employments, employeeID)
	emp := <-employments

	count, err := emp.propagateToHub(token.Email, token.ORCID)
	if err != nil {
		return "", err
	}
	taskRecordCount += count

	return "", nil
}

func getIdentidy(output chan<- Identity, upiOrID string) {
	var id Identity
	err := api.Get("identity/integrations/v3/identity/"+upiOrID, &id)
	if err != nil {
		log.Fatalln("Failed to retrieve the identity record: ", err)
	}
	output <- id
}

func getEmp(output chan<- Employment, upiOrID string) {
	var emp Employment
	err := api.Get("employment/integrations/v1/employee/"+upiOrID, &emp)
	if err != nil {
		log.Fatalln("Failed to get employment record: ", err)
	}
	output <- emp
}

func (id *Identity) updateOrcid(done chan<- bool, ORCID string) {
	defer func() {
		done <- true
	}()

	currentORCID := id.GetORCID()
	if currentORCID != "" {
		if ORCID != currentORCID {
			// TODO
		}
		return
	}
	// Add ORCID ID if the user doesn't have one
	var resp struct {
		StatusCode string `json:"statusCode"`
	}
	err := api.Put(fmt.Sprintf("identity/integrations/v3/identity/%d/identifier/ORCID", id.ID), map[string]string{"identifier": ORCID}, &resp)
	if err != nil {
		log.Println("ERROR: Failed to update or add ORCID: ", err)
	}
}

func (e *Event) processUserRegistration() (string, error) {
	var (
		id  Identity
		emp Employment
	)
	identities := make(chan Identity, 1)
	employments := make(chan Employment, 1)

	parts := strings.Split(e.EPPN, "@")
	upi := parts[0]
	log.Println("UPI: ", upi)

	go getIdentidy(identities, upi)
	go getEmp(employments, upi)

	id = <-identities
	if id.ID != 0 {
		idUpdateDone := make(chan bool, 1)
		go id.updateOrcid(idUpdateDone, e.ORCID)
		defer func() {
			<-idUpdateDone
		}()
	}
	emp = <-employments

	if id.ID == 0 || emp.Job == nil {
		return "", fmt.Errorf("no Identity for %q (%s, %s) or employment records", e.EPPN, e.Email, e.ORCID)
	}

	count, err := emp.propagateToHub(id.EmailAddress, e.ORCID)
	if err != nil {
		return "", err
	}
	taskRecordCount += count
	return fmt.Sprintf("%#v", id), nil
}

type errorList []error

func (el errorList) Error() string {
	var sb strings.Builder
	for i, e := range el {
		if i > 0 {
			sb.WriteString("; ")
		}
		sb.WriteString(e.Error())
	}
	return sb.String()
}

func HandleRequest(ctx context.Context, e Event) (string, error) {

	log.Printf("Cotext: %#v, counter: %d", ctx, counter)
	counter += 1

	if e.Records != nil {
		var (
			resp   []string
			errors errorList
		)
		for _, r := range e.Records {
			var e Event
			json.Unmarshal([]byte(r.Body), &e)
			rv, err := e.process()
			if err != nil {
				errors = append(errors, err)
			}
			resp = append(resp, rv)
		}
		return strings.Join(resp, "; "), errors
	}
	return e.process()
}

func main() {
	log.SetPrefix("OHI")
	lambda.Start(HandleRequest)
}

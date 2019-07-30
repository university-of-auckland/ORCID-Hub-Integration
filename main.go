package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/aws/aws-lambda-go/lambda"
)

var (
	api                  Client
	oh                   Client
	counter              int
	taskSetUpWG          sync.WaitGroup
	gotAccessTokenWG     sync.WaitGroup
	taskID               int
	taskCreatedAt        time.Time
	taskRecordCount      int
	taskRecordCountMutex sync.Mutex
	verbose              bool
)

const taskFilenamePrefix = "UOA-OH-INTEGRATION-TASK-"

var (
	// APIBaseURL is the UoA API base URL
	APIBaseURL = "https://api.dev.auckland.ac.nz/service"
	// OHBaseURL is the ORCID Hub API base URL
	OHBaseURL = "https://dev.orcidhub.org.nz"
)

func setup() {
	if api.ApiKey == "" {
		api.ApiKey = os.Getenv("API_KEY")
		api.BaseURL = APIBaseURL
	}
	if oh.AccessToken == "" {
		gotAccessTokenWG.Add(1)
		go func() {
			oh.ClientID = os.Getenv("CLIENT_ID")
			oh.ClientSecret = os.Getenv("CLIENT_SECRET")
			oh.BaseURL = OHBaseURL
			// oh.BaseURL = "http://127.0.0.1:5000"
			err := oh.GetAccessToken("oauth/token")
			if err != nil {
				log.Panic(err)
			}
			gotAccessTokenWG.Done()
		}()
	}
	go setupTask()
}

// process performs the incoming message routing.
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

// processEmpUpdate handles the employer update event.
func (e *Event) processEmpUpdate() (string, error) {

	var employeeID = strconv.Itoa(e.Subject)
	identities := make(chan Identity)
	employments := make(chan Employment)

	// TODO: this can be doen sychroniously
	go getIdentidy(identities, employeeID)
	id := <-identities
	token, ok := id.GetOrcidAccessToken()

	if !ok {
		return "", fmt.Errorf("the user (ID: %s) hasn't granted access to the profile", employeeID)
	}

	// TODO: this can be doen sychroniously
	go getEmp(employments, employeeID)
	emp := <-employments

	go emp.propagateToHub(token.Email, token.ORCID)

	return "", nil
}

// getIdentidy retrieves the user identity records.
func getIdentidy(output chan<- Identity, upiOrID string) {
	var id Identity
	err := api.Get("identity/integrations/v3/identity/"+upiOrID, &id)
	if err != nil {
		log.Fatalln("Failed to retrieve the identity record: ", err)
	}
	output <- id
}

// getEmp retrieves the user employment records.
func getEmp(output chan<- Employment, upiOrID string) {
	var emp Employment
	err := api.Get("employment/integrations/v1/employee/"+upiOrID, &emp)
	if err != nil {
		log.Fatalln("Failed to get employment record: ", err)
	}
	output <- emp
}

// updateOrcid updates the user ORCID iD.
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

// isValidUPI validates UPI
func isValidUPI(upi string) bool {
	if len(upi) != 7 {
		return false
	}
	for _, r := range upi[0:4] {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	for _, r := range upi[4:] {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// processUserRegistration handles the user registration/ORCID account linking on the Hub.
func (e *Event) processUserRegistration() (string, error) {
	var (
		id  Identity
		emp Employment
	)
	identities := make(chan Identity)
	employments := make(chan Employment)

	parts := strings.Split(e.EPPN, "@")

	upi := parts[0]
	if !isValidUPI(upi) {
		return "Invalid UPI", fmt.Errorf("Invalid UPI: %q", upi)
	}
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

// HandleRequest handle "AWS lambda" request with a single event message or
// a batch of event messages.
func HandleRequest(ctx context.Context, e Event) (string, error) {

	log.Printf("Cotext: %#v, counter: %d", ctx, counter)
	counter++

	if e.Records != nil {
		var (
			resp   []string
			errors errorList
		)

		type restponse struct {
			message string
			err     error
		}

		output := make(chan restponse, len(e.Records))
		for _, r := range e.Records {
			var e Event
			json.Unmarshal([]byte(r.Body), &e)

			go func(e Event, o chan<- restponse) {
				resp, err := e.process()
				o <- restponse{resp, err}
			}(e, output)
		}
		for range e.Records {
			r := <-output
			if r.err != nil {
				errors = append(errors, r.err)
			}
			resp = append(resp, r.message)
		}
		return strings.Join(resp, "; "), errors
	}
	return e.process()
}

func main() {
	log.SetPrefix("OHI: ")
	lambda.Start(HandleRequest)
	log.Println("=================== DONE =============================================")
}

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"

	"github.com/dougEfresh/lambdazap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	api                  Client
	counter              int
	gotAccessTokenWG     sync.WaitGroup
	oh                   Client
	taskCreatedAt        time.Time
	taskID               int
	taskRecordCount      int
	taskRecordCountMutex sync.Mutex
	updateOrcidWG        sync.WaitGroup
	verbose              bool
	wg                   sync.WaitGroup

	loggingLevel zapcore.Level
	lambdazapper *lambdazap.LambdaLogContext
	logger       *zap.Logger
	log          *zap.SugaredLogger
)

const taskFilenamePrefix = "UOA-OH-INTEGRATION-TASK-"

var (
	// APIBaseURL is the UoA API base URL
	APIBaseURL = "https://api.dev.auckland.ac.nz/service"
	// OHBaseURL is the ORCID Hub API base URL
	OHBaseURL = "https://dev.orcidhub.org.nz"
)

func init() {
	godotenv.Load()

	isDevelopment := strings.Contains(os.Getenv("ENV"), "dev")

	if verbose || isDevelopment {
		loggingLevel = zap.DebugLevel
	} else {
		loggingLevel = zap.InfoLevel
	}
	logger, _ = zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
		Development: isDevelopment,
		Encoding:    "console",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "T",
			LevelKey:       "L",
			NameKey:        "N",
			CallerKey:      "C",
			MessageKey:     "M",
			StacktraceKey:  "S",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}.Build()
	if os.Getenv("_LAMBDA_SERVER_PORT") != "" {
		lambdazapper = lambdazap.New().With(lambdazap.AwsRequestID)
		logger.With(lambdazapper.NonContextValues()...)
	}
	log = logger.Sugar()
}

func setup() {
	setupAPIClients()
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
	err := api.get("identity/integrations/v3/identity/"+upiOrID, &id)
	if err != nil {
		log.Fatal("failed to retrieve the identity record", zap.Error(err))
	}
	output <- id
}

// getEmp retrieves the user employment records.
func getEmp(output chan<- Employment, upiOrID string) {
	var emp Employment
	err := api.get("employment/integrations/v1/employee/"+upiOrID, &emp)
	if err != nil {
		log.Fatal("failed to get employment record", zap.Error(err))
	}
	output <- emp
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

// isValidID validates employment/student ID
func isValidID(uid string) bool {
	if l := len(uid); l < 8 || l > 10 {
		return false
	}
	for _, r := range uid {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// processUserRegistration handles the user registration/ORCID account linking on the Hub.
func (e *Event) processUserRegistration() (string, error) {

	parts := strings.Split(e.EPPN, "@")
	upi := parts[0]
	if !isValidUPI(upi) {
		return "", fmt.Errorf("Invalid UPI: %q", upi)
	}
	log.Info("UPI: ", upi)

	var (
		id  Identity
		emp Employment
	)
	identities := make(chan Identity)
	employments := make(chan Employment)

	go getIdentidy(identities, upi)
	go getEmp(employments, upi)

	id = <-identities
	if id.ID != 0 {
		go id.updateOrcid(e.ORCID)
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

	defer func() {
		wg.Wait()
		logger.Sync()
	}()
	counter++
	log.Infof("Context: %+v, counter: %d", ctx, counter)
	log.Infof("Event message: %#v", e)

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

	lambda.Start(HandleRequest)
	log.Info("=================== DONE =============================================")
}

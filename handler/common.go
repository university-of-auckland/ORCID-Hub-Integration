package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/joho/godotenv"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	taskFilenamePrefix      = "UOA-OH-INTEGRATION-TASK-"
	defaultTaskRetentionMin = 5.0
	defaultBatchSize        = 12
)

var (
	api                  Client
	batchSize            = defaultBatchSize
	counter              int
	gotAccessTokenWG     sync.WaitGroup
	log                  *zap.SugaredLogger
	logger               *zap.Logger
	loggingLevel         zapcore.Level
	oh                   Client
	taskCreatedAt        time.Time
	taskID               int
	taskRecordCount      int
	taskRecordCountMutex sync.Mutex
	taskRetentionMin     = defaultTaskRetentionMin
	updateOrcidWG        sync.WaitGroup
	verbose              bool
	wg                   sync.WaitGroup

	// for testing/mocking
	logFatal func(args ...interface{})
)

var (
	// APIBaseURL is the UoA API base URL
	APIBaseURL = "https://api.dev.auckland.ac.nz/service"
	// OHBaseURL is the ORCID Hub API base URL
	OHBaseURL = "https://dev.orcidhub.org.nz"
)

func getenvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func init() {
	godotenv.Load()

	isDevelopment := strings.Contains(os.Getenv("ENV"), "dev")

	if verbose || isDevelopment {
		loggingLevel = zap.DebugLevel
	} else {
		loggingLevel = zap.InfoLevel
	}
	logger, _ = zap.Config{
		Level:       zap.NewAtomicLevelAt(loggingLevel),
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
	log = logger.Sugar()
	logFatal = log.Fatal
}

func setup(wg *sync.WaitGroup) {
	setupAPIClients()

	go setupTask(wg)
}

// handle performs the incoming message routing.
func (e *Event) handle() (string, error) {

	counter++
	log.Infof("Event message #%d: %+v", counter, e)

	if e.Records != nil {
		var (
			resp   []string
			errors errorList
			events []Event
		)

		type restponse struct {
			message string
			err     error
		}

		for _, r := range e.Records {
			var e Event
			json.Unmarshal([]byte(r.Body), &e)

			if e.Subject != 0 || (e.EPPN != "" && e.Type == "CREATED") {
				events = append(events, e)
			}
		}

		if events == nil {
			return "", nil
		}

		output := make(chan restponse, len(events))
		for _, e := range events {
			go func(e Event, o chan<- restponse) {
				resp, err := e.handle()
				o <- restponse{resp, err}
			}(e, output)
		}
		for range events {
			rr := <-output
			if rr.err != nil {
				errors = append(errors, rr.err)
			}
			resp = append(resp, rr.message)
		}
		return strings.Join(resp, "; "), errors
	}

	var setUpWG sync.WaitGroup

	if (e.EPPN != "" && e.Type == "CREATED") || e.Subject != 0 || e.Type == "PING" {
		setUpWG.Add(1)
		go setup(&setUpWG)

		if e.EPPN != "" {
			return e.processUserRegistration(&setUpWG)
		} else if e.Subject != 0 {
			return e.processEmpUpdate(&setUpWG)
		} else if e.Type == "PING" { // Heartbeat Check
			setUpWG.Wait()
			return "GNIP", nil
		}
	}
	return "", fmt.Errorf("unhandled event: %#v", e)
}

// processEmpUpdate handles the employer update event.
func (e *Event) processEmpUpdate(wg *sync.WaitGroup) (string, error) {

	var employeeID = strconv.Itoa(e.Subject)

	var id Identity
	err := api.get("identity/integrations/v3/identity/"+employeeID, &id)
	if err != nil {
		log.Fatal("failed to retrieve the identity record", err)
	}
	if id.Upi == "" {
		return "", errors.New("failed to retrieve the identity record")
	}

	token, ok := id.GetOrcidAccessToken()
	if !ok {
		return "", fmt.Errorf("the user (ID: %s) hasn't granted access to the profile", employeeID)
	}

	var emp Employment
	err = api.get("employment/integrations/v1/employee/"+employeeID, &emp)
	if err != nil {
		log.Fatal("failed to get employment record", zap.Error(err))
	}

	emp.propagateToHub(token.Email, token.ORCID, wg)

	return "", nil
}

// getIdentidy retrieves the user identity records.
func getIdentidy(output chan<- Identity, upiOrID string) {
	var id Identity
	err := api.get("identity/integrations/v3/identity/"+upiOrID, &id)
	if err != nil {
		logFatal("failed to retrieve the identity record", err)
	}
	output <- id
}

// getEmp retrieves the user employment records.
func getEmp(output chan<- Employment, upiOrID string) {
	var emp Employment
	err := api.get("employment/integrations/v1/employee/"+upiOrID, &emp)
	if err != nil {
		logFatal("failed to get employment record", err)
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
func (e *Event) processUserRegistration(wg *sync.WaitGroup) (string, error) {

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

	count, err := emp.propagateToHub(id.EmailAddress, e.ORCID, wg)
	taskRecordCountMutex.Lock()
	taskRecordCount += count
	taskRecordCountMutex.Unlock()
	return fmt.Sprintf("%#v", id), err
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

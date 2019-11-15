//+build !test,!heroku,!container,!standalone

package main

import (
	"context"
	"encoding/base64"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/ssm"

	"github.com/dougEfresh/lambdazap"
)

var (
	kmsClient    *kms.KMS
	ssmClient    *ssm.SSM
	lambdazapper *lambdazap.LambdaLogContext
	isLambda     bool
)

// awsPsPrefix - AWS Paramter Store parameter name prefix
const awsPsPrefix = "/ORCIDHUB-INTEGRATION-"

// HandleRequest handle "AWS lambda" request with a single event message or
// a batch of event messages.
func HandleRequest(ctx context.Context, e Event) (string, error) {

	defer func() {
		logger.Sync()
	}()

	return e.handle()
}

func main() {

	if isLambda {
		lambdazapper = lambdazap.New().With(lambdazap.AwsRequestID)
		logger.With(lambdazapper.NonContextValues()...)
		log = logger.Sugar()
		logFatal = log.Fatal
	}

	lambda.Start(HandleRequest)
}

func init() {

	env = os.Getenv("ENV")
	isLambda = os.Getenv("_LAMBDA_SERVER_PORT") != ""
	if isLambda {
		s, err := session.NewSession()
		if err != nil {
			log.Fatal(err)
		}
		kmsClient = kms.New(s)
		ssmClient = ssm.New(s)
	}
	go func() {
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGPIPE, syscall.SIGTERM)

	TASK_HANDLING:
		for {
			select {
			// every 10 min check if the current task can be submitted for processing
			case <-time.Tick(time.Minute * 10):
				if taskID != 0 && taskRecordCount > batchSize && time.Since(taskCreatedAt).Minutes() > taskRetentionMin {
					(&Task{ID: taskID}).activate()
					newTask()
				}
			case <-sc:
				// activate the current task (if it might be activated) at the shutdown
				if taskID != 0 && taskRecordCount > batchSize && time.Since(taskCreatedAt).Minutes() > taskRetentionMin {
					(&Task{ID: taskID}).activate()
				}
				log.Info("service terminated")
				break TASK_HANDLING
			}
		}
		close(sc)
		logger.Sync()
	}()
}

// getenv returns enviroment variable value if it's defined
// or the default. If the value is encrypted, it will depcrypt it first.
func getenv(key, defaultValue string) (value string) {
	if isLambda && (key == "APIKEY" || key == "CLIENT_ID" || key == "CLIENT_SECRET") {
		keyname := awsPsPrefix + key
		if env != "" {
			keyname = "/" + env + keyname
		}
		log.Debugf("Reading parameter %q", keyname)
		withDecryption := true
		param, err := ssmClient.GetParameter(
			&ssm.GetParameterInput{
				Name:           &keyname,
				WithDecryption: &withDecryption,
			})
		if err != nil {
			log.Errorf("Failed to retrieve parameter %q: %v", keyname, err)
		} else {
			value = *param.Parameter.Value
		}
	}
	// attempt to use the environment variable
	if value == "" {
		value = os.Getenv(key)
	}

	if value != "" {
		log.Debug("KEY: ", key, ", VALUE: ", value)
		// unecrypted or looks unencrypted
		if !isLambda || len(value) < 40 || !strings.Contains(value, "+") {
			return value
		}
		// Decrypt secrets:
		decodedBytes, err := base64.StdEncoding.DecodeString(value)
		if err != nil {
			log.Fatal(err)
		}
		input := &kms.DecryptInput{
			CiphertextBlob: decodedBytes,
		}
		response, err := kmsClient.Decrypt(input)
		if err != nil {
			log.Fatal(err)
		}
		// Plaintext is a byte array, so convert to string
		value = string(response.Plaintext[:])
		return value
	}
	return defaultValue
}

//+build !test,!heroku,!container

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

// HandleRequest handle "AWS lambda" request with a single event message or
// a batch of event messages.
func HandleRequest(ctx context.Context, e Event) (string, error) {

	defer func() {
		logger.Sync()
		wg.Wait()
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

	isLambda = os.Getenv("_LAMBDA_SERVER_PORT") != ""
	if isLambda {
		kmsClient = kms.New(session.New())
		ssmClient = ssm.New(session.New())
	}
	go func() {
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGPIPE, syscall.SIGKILL, syscall.SIGTERM)

	TASK_HANDLING:
		for {
			select {
			case <-time.Tick(time.Minute * 10):
				if taskID != 0 && taskRecordCount > batchSize && time.Now().Sub(taskCreatedAt).Minutes() > taskRetentionMin {
					(&Task{ID: taskID}).activate()
					newTask()
				}
			case <-sc:
				if taskID != 0 && taskRecordCount > batchSize && time.Now().Sub(taskCreatedAt).Minutes() > taskRetentionMin {
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
func getenv(key, defaultValue string) string {
	var value string
	if isLambda {
		keyname := "ORCIDHUB_INTEGRATION_LAMBDA_" + key
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
			panic(err)
		}
		input := &kms.DecryptInput{
			CiphertextBlob: decodedBytes,
		}
		response, err := kmsClient.Decrypt(input)
		if err != nil {
			panic(err)
		}
		// Plaintext is a byte array, so convert to string
		value = string(response.Plaintext[:])
		return value
	}
	return defaultValue
}

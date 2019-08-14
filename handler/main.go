//+build !test,!heroku,!container

package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/dougEfresh/lambdazap"
)

var lambdazapper *lambdazap.LambdaLogContext

// Response - lambda response
type Response struct {
	Message string `json:"message,omitempty"`
	Retry   bool   `json:"retry"`
}

// HandleRequest handle "AWS lambda" request with a single event message or
// a batch of event messages.
func HandleRequest(ctx context.Context, e Event) (Response, error) {

	defer func() {
		logger.Sync()
		wg.Wait()
		logger.Sync()
	}()

	message, err := e.handle()
	if err != nil {
		message += ": " + err.Error()
	}
	return Response{Message: message, Retry: err != nil}, err

}

func main() {

	if os.Getenv("_LAMBDA_SERVER_PORT") != "" {
		lambdazapper = lambdazap.New().With(lambdazap.AwsRequestID)
		logger.With(lambdazapper.NonContextValues()...)
		log = logger.Sugar()
		logFatal = log.Fatal
	}

	lambda.Start(HandleRequest)
}

func init() {
	go func() {
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGPIPE, syscall.SIGKILL, syscall.SIGTERM)

	TASK_HANDLING:
		for {
			select {
			case <-time.Tick(time.Minute * 10):
				if taskID != 0 && taskRecordCount > 0 && time.Now().Sub(taskCreatedAt).Minutes() > taskRetentionMin {
					taskSetUpWG.Add(1)
					go (&Task{ID: taskID}).activateTask()
					taskSetUpWG.Add(1)
					go newTask()
					taskSetUpWG.Done()
				}
			case <-sc:
				if taskID != 0 && taskRecordCount > 0 && time.Now().Sub(taskCreatedAt).Minutes() > taskRetentionMin {
					taskSetUpWG.Add(1)
					go (&Task{ID: taskID}).activateTask()
					taskSetUpWG.Done()
				}
				log.Info("service terminated")
				break TASK_HANDLING
			}
		}
		close(sc)
		logger.Sync()
	}()
}

//+build !test

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

// HandleRequest handle "AWS lambda" request with a single event message or
// a batch of event messages.
func HandleRequest(ctx context.Context, e Event) {

	defer func() {
		logger.Sync()
		wg.Wait()
		logger.Sync()
	}()

	e.handle()
	return

}

func main() {

	if os.Getenv("_LAMBDA_SERVER_PORT") != "" {
		lambdazapper = lambdazap.New().With(lambdazap.AwsRequestID)
		logger.With(lambdazapper.NonContextValues()...)
		log = logger.Sugar()
	}

	lambda.Start(HandleRequest)
}

func init() {
	go func() {
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGPIPE, syscall.SIGKILL, syscall.SIGTERM)

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
				break
			}
		}
	}()
}

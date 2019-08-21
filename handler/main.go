//+build !test,!heroku,!container

package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/dougEfresh/lambdazap"
)

var lambdazapper *lambdazap.LambdaLogContext

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

		var wg sync.WaitGroup
	TASK_HANDLING:
		for {
			select {
			case <-time.Tick(time.Minute * 10):
				if taskID != 0 && taskRecordCount > batchSize && time.Now().Sub(taskCreatedAt).Minutes() > taskRetentionMin {
					wg.Add(2)
					go (&Task{ID: taskID}).activate(&wg)
					go newTask(&wg)
				}
			case <-sc:
				if taskID != 0 && taskRecordCount > batchSize && time.Now().Sub(taskCreatedAt).Minutes() > taskRetentionMin {
					wg.Add(1)
					go (&Task{ID: taskID}).activate(&wg)
				}
				log.Info("service terminated")
				break TASK_HANDLING
			}
			wg.Wait()
		}
		close(sc)
		logger.Sync()
	}()
}

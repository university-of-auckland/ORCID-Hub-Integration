//+build !test

package main

import (
	"context"
	"os"

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

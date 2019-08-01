//+build !test

package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/dougEfresh/lambdazap"
)

var lambdazapper *lambdazap.LambdaLogContext

func main() {

	if os.Getenv("_LAMBDA_SERVER_PORT") != "" {
		lambdazapper = lambdazap.New().With(lambdazap.AwsRequestID)
		logger.With(lambdazapper.NonContextValues()...)
		log = logger.Sugar()
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGPIPE)

	go func() {
		<-sc
		if taskID != 0 {
			log.Info("====================== SIGPIPE ======================================================")
			log.Infof("task (ID: %d) activated", taskID)
			log.Info("====================== SIGPIPE ======================================================")
			logger.Sync()
		}
	}()

	lambda.Start(HandleRequest)
}

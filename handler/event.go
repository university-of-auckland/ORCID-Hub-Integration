package main

import "github.com/aws/aws-lambda-go/events"

// Event - a generic message suitable for both EMP Update event and
// ORCIDHub Webhook propagagted event (with wrapped in SQS message batch):
type Event struct {
	EPPN    string `json:"eppn"`
	Email   string `json:"email"`
	ORCID   string `json:"orcid"`
	Subject int    `json:"subject,string"`
	Type    string `json:"type"`
	URL     string `json:"url"`
	// SQS Message if used SQS
	Records []events.SQSMessage
}

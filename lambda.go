package main

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

// AWS SDK session and Lambda service objects
var sess session.Session
var svc *lambda.Lambda

// AWS session and service objects are created during startup and
// reused across all Lambda invocations. If for any reason they turn unusable,
// there is no attempt to recreate them.
func init() {

	sess, err := session.NewSession()
	if err != nil {
		log.WithError(err).Panic("unable to acquire AWS session")
	}

	svc = lambda.New(sess)

}

// Function lambdaCall - call Lambda using ApplicId of MQTM as function name
// and the entire MQTM as payload
func lambdaCall(lambdaName string, tm MQTM) error {

	payload, err := json.Marshal(tm)

	if err != nil {
		log.WithError(err).Error("json.Marshal failed")
		return err
	}

	log.WithFields(log.Fields{
		"NAME":    lambdaName,
		"PAYLOAD": string(payload),
	}).Debug("about to invoke Lambda")

	params := &lambda.InvokeInput{
		FunctionName:   aws.String(lambdaName),
		InvocationType: aws.String("Event"),
		Payload:        payload,
	}

	resp, err := svc.Invoke(params)

	// There is no check for intermittent failures, nor is there any
	// retry attempt. It is assumed that MQ will eventually retrigger the process
	if err != nil {
		log.WithError(err).Error("Lambda invocation failed")
		return err
	}

	log.WithFields(log.Fields{
		"RESPONSE": resp,
	}).Debug("successful Lambda invocation")

	return nil
}

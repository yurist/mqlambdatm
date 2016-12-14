package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

var sess session.Session

func lambdaCall(lambdaName string, tm MQTM) error {

	sess, err := session.NewSession()
	if err != nil {
		return err
	}

	svc := lambda.New(sess)

	payload, err := json.Marshal(tm)

	fmt.Println(string(payload))

	if err != nil {
		fmt.Println("json error ", err)
		return err
	}

	params := &lambda.InvokeInput{
		FunctionName:   aws.String(lambdaName),
		InvocationType: aws.String("Event"),
		Payload:        payload,
	}

	resp, err := svc.Invoke(params)

	fmt.Println(resp, err)

	return err
}

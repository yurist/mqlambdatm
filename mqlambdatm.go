package main

/*

  A simple MQ trigger monitor for AWS Lambda.
  Follows pretty closely a usual trigger monitor logic, like one found in
  /opt/mqm/samp/amqstrg0.c

  Can be invoked quite similarly to runmqtrm, either from command line or
  using MQSC DEFINE SERVICE:

  mqlambdatm -q <initiation-queue> [-m <queue-manager>] [--log-level <log-level>]

  -q parameter is mandatory (nobody I know of uses SYSTEM.DEFAULT.INITIATION.QUEUE anyway)

  --log-level defaults to INFO (case is immaterial,) set it to DEBUG to get trace information
    about every trigger message

  The name of Lambda function must be encoded in ApplicId of the triggered process, case-sensitive.
  The function is invoked asynchronously (Event invocation) and passed a JSON-encoded trigger
  message. The JSON encoding uses the MQTM field names from
  cmqc.h *with the first letter in lower case* to permit Java bean deserialization
  on the Lambda side.

  AWS credential chain and AWS region must be externally configured. See AWS documentation for possible
  ways of configuring the credential chain. As for the region, the simplest way to make it work
  is to set AWS_REGION environment variable. (Warning - AWS_DEFAULT_REGION may require additional
  tweaking, see AWS Go SDK documentation for details if you want to use it.)

  Invalid trigger messages are skipped with an error message (no dead letter queue etc.)

  Any MQI error on MQGET will terminate the program. No retries.

*/

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/ibm-messaging/mq-golang/ibmmq"
	"strings"
)

// Initiation queue to serve
var initQ string

// Queue manager to connect
var qMgrName string

func init() {
	flag.StringVar(&initQ, "q", "", "initiation queue to serve")
	flag.StringVar(&qMgrName, "m", "", "queue manager to connect, default queue manager if not given)")
	sLogLevel := flag.String("log-level", "info", "log level (DEBUG, INFO, WARN, ERROR, FATAL, PANIC)")

	flag.Parse()

	logLevel, err := log.ParseLevel(*sLogLevel)

	if err != nil {
		log.WithField("LOG-LEVEL", sLogLevel).Error("invalid log level, INFO assumed")
		logLevel = log.InfoLevel
	}

	log.SetLevel(logLevel)
}

func main() {

	log.Infoln("MQ trigger monitor for AWS Lambda")

	if initQ == "" {
		log.Fatalln("-q parameter missing")
	}

	qMgr, mqreturn, err := ibmmq.Conn(qMgrName)
	if err != nil {
		log.WithFields(log.Fields{
			"MQRC": mqreturn.MQRC,
		}).Fatal("error connecting to queue manager")
	}

	defer qMgr.Disc()

	mqod := ibmmq.NewMQOD()

	mqod.ObjectType = ibmmq.MQOT_Q
	mqod.ObjectName = initQ

	var openOpts int32 = ibmmq.MQOO_INPUT_AS_Q_DEF + ibmmq.MQOO_FAIL_IF_QUIESCING

	qObj, mqreturn, err := qMgr.Open(mqod, openOpts)

	if err != nil {
		log.WithFields(log.Fields{
			"INITQ": initQ,
			"MQRC":  mqreturn.MQRC,
		}).Fatal("error opening initiation queue")
	}

	defer qObj.Close(ibmmq.MQCO_NONE)

	md := ibmmq.NewMQMD()
	gmo := ibmmq.NewMQGMO()

	gmo.Version = ibmmq.MQGMO_VERSION_2
	gmo.MatchOptions = ibmmq.MQGMO_NONE

	gmo.Options = ibmmq.MQGMO_WAIT +
		ibmmq.MQGMO_FAIL_IF_QUIESCING +
		ibmmq.MQGMO_ACCEPT_TRUNCATED_MSG +
		ibmmq.MQGMO_NO_SYNCPOINT

	gmo.WaitInterval = ibmmq.MQWI_UNLIMITED

	// ibmmq package currently doesn't provide correct MQTM structure length,
	// so we are using a constant

	const mqtm_length = 684 // should be ibmmq.MQTM_CURRENT_LENGTH

	msg := make([]byte, mqtm_length)

	for {
		datalen, mqreturn, err := qObj.Get(md, gmo, msg)

		if err != nil {
			log.WithFields(log.Fields{
				"INITQ": initQ,
				"MQRC":  mqreturn.MQRC,
			}).Fatal("error getting a message")
		}

		if datalen != mqtm_length {
			log.Error("invalid message received, skipping (wrong length)")
			continue
		}

		tm := TMfromC(msg)

		if tm.StrucId != "TM  " {
			log.Error("invalid message received, skipping (wrong StrucId)")
			continue
		}

		log.WithFields(log.Fields{
			"TM": tm,
		}).Debug("trigger message received")

		// Ignoring any error, lambdaCall is logging its errors
		lambdaCall(strings.TrimSpace(tm.ApplId), tm)
	}
}

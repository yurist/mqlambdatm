package main

import (
	"flag"
	"fmt"
	"github.com/ibm-messaging/mq-golang/ibmmq"
	"os"
	"strings"
)

var initQ string
var qMgrName string

func init() {
	flag.StringVar(&initQ, "q", "", "initiation queue to serve")
	flag.StringVar(&qMgrName, "m", "", "queue manager to connect, default queue manager if not given)")

	flag.Parse()
}

func main() {

	fmt.Println("MQ trigger monitor for AWS Lambda")

	if initQ == "" {
		fmt.Println("-q parameter missing")
		os.Exit(1)
	}

	qMgr, mqreturn, err := ibmmq.Conn(qMgrName)
	if err != nil {
		fmt.Printf("Error connecting to queue manager %v", mqreturn.MQRC)
		os.Exit(1)
	}

	defer qMgr.Disc()

	mqod := ibmmq.NewMQOD()

	mqod.ObjectType = ibmmq.MQOT_Q
	mqod.ObjectName = initQ

	var openOpts int32 = ibmmq.MQOO_INPUT_AS_Q_DEF + ibmmq.MQOO_FAIL_IF_QUIESCING

	qObj, mqreturn, err := qMgr.Open(mqod, openOpts)

	if err != nil {
		fmt.Printf("Error opening initiation queue %v %v", initQ, mqreturn.MQRC)
		os.Exit(1)
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

	msg := make([]byte, 684)
	// msg := make([]byte, ibmmq.MQTM_CURRENT_LENGTH)

	for {
		datalen, mqreturn, err := qObj.Get(md, gmo, msg)

		if err != nil {
			fmt.Printf("Error getting a message off queue %v %v\n", initQ, mqreturn.MQRC)
			os.Exit(1)
		}

		if datalen != 684 {
			fmt.Println("Invalid message received, skipping (wrong length)")
			continue
		}

		tm := TMfromC(msg)

		if tm.StrucId != "TM  " {
			fmt.Println("Invalid message received, skipping (wrong StrucId)")
		}

		fmt.Printf("msg %v\n", tm)

		lambdaCall(strings.TrimSpace(tm.ApplId), tm)
	}

}

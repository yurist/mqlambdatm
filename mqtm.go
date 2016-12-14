package main

/*
#cgo CFLAGS: -I/opt/mqm/inc -D_REENTRANT

#include <stdlib.h>
#include <string.h>
#include <cmqc.h>

*/
import "C"

import (
	// "bytes"
	"unsafe"
)

type MQTM struct {
	StrucId     string `json:"strucId"`
	Version     int32  `json:"version"`
	QName       string `json:"qName"`
	ProcessName string `json:"processName"`
	TriggerData string `json:"triggerData"`
	ApplType    int32  `json:"applType"`
	ApplId      string `json:"applId"`
	EnvData     string `json:"envData"`
	UserData    string `json:"userData"`
}

func TMfromC(buffer []byte) (gotm MQTM) {

	mqtm := (C.PMQTM)(unsafe.Pointer(&buffer[0]))

	gotm.StrucId = C.GoStringN((*C.char)(&mqtm.StrucId[0]), 4)
	gotm.Version = int32(mqtm.Version)
	gotm.QName = C.GoStringN((*C.char)(&mqtm.QName[0]), 48)
	gotm.ProcessName = C.GoStringN((*C.char)(&mqtm.ProcessName[0]), 48)
	gotm.TriggerData = C.GoStringN((*C.char)(&mqtm.TriggerData[0]), 64)
	gotm.ApplType = int32(mqtm.ApplType)
	gotm.ApplId = C.GoStringN((*C.char)(&mqtm.ApplId[0]), 256)
	gotm.EnvData = C.GoStringN((*C.char)(&mqtm.EnvData[0]), 128)
	gotm.UserData = C.GoStringN((*C.char)(&mqtm.UserData[0]), 128)

	return
}

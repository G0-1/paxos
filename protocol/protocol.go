package protocol

import (
	"encoding/json"
	"log"
)

type cmdPack struct {
	Cmd string
	ID  string
}

/* cmd: cre,inc,get; id */
func ParseCmd(cmdJson []byte) *cmdPack {
	cmdData := new(cmdPack)
	err := json.Unmarshal(cmdJson, cmdData)
	if err == nil {
		return cmdData
	} else {
		log.Println("parse cmdJson error", err)
		return nil
	}
}

func AsmCmd(cmd, id string) []byte {
	cmdData := &cmdPack{Cmd: cmd, ID: id}
	cmdJson, err := json.Marshal(cmdData)
	if err != nil {
		log.Println("asm cmdJson error: ", err)
		return nil
	}
	cmdJson = append(cmdJson, '\n')
	return cmdJson
}

/*
{"cmd": "cre","id": "1111111111111111111"}
{"cmd": "get","id": "1111111111111111111"}
{"cmd": "inc","id": "1111111111111111111"}
*/

type respPack struct {
	Code uint8
	Seq  uint64
	Msg  string
}

func ParseResp(respJson []byte) *respPack {
	respData := new(respPack)
	err := json.Unmarshal(respJson, respData)
	if err == nil {
		return respData
	} else {
		log.Println("parse respJson error", err)
		return nil
	}
}

/* code, seq, msg */
func AsmResp(code uint8, seq uint64, msg string) []byte {
	respData := &respPack{Code: code, Seq: seq, Msg: msg}
	respJson, err := json.Marshal(respData)
	if err != nil {
		log.Println("asm respJson error: ", err)
		return nil
	}

	respJson = append(respJson, '\n')
	return respJson
}

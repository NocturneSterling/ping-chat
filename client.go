package main

import (
	"encoding/json"
	"time"
)

type ChatMessage struct {
	Message string `json:"msg"`
	User    string `json:"user"`
}

var lastTimestamp int

func runClientSender(msg string) {
	chanNum := 0
	msgJson := ChatMessage{Message: msg, User: currentUser}
	pass := channelPass(chanNum)
	jsonBytes, _ := json.Marshal(msgJson)
	hash := passHash(pass)
	sendBytes(append(hash, encryptToBytes(jsonBytes, []byte(currentPass))...), *ip)
}

func runClientListener(numChans int) {
	name := currentChannel(numChans)
	pass := channelPass(numChans)
	lastTs := 0
	for {
		passwordHashBytes := passHash(pass)
		responseBytes := sendBytes(passwordHashBytes, *ip)
		responseStr := decryptFromBytes(responseBytes, passwordHashBytes)
		var response MsgRecord
		err := json.Unmarshal(responseStr, &response)
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		var incomingMsgJson ChatMessage
		msgTextStr := decryptUsingPass(response.LastMsgEncrypted, pass)
		json.Unmarshal([]byte(msgTextStr), &incomingMsgJson)
		if response.LastMsgTimestamp != lastTimestamp {
			if incomingMsgJson.Message == "" && incomingMsgJson.User == "" {
				tuiPrint(name, "Chat begins here")
			} else {
				tuiPrint(name, incomingMsgJson.User + ": " + incomingMsgJson.Message)
			}
			lastTs = response.LastMsgTimestamp
		}
		time.Sleep(1 * time.Second)
	}
}

func listenClient(){
	for i := 0; i < listNum; i++{
	go runClientListener(i)
	}
}

func runClient() {
	initTUI(runClientSender)
	//go runClientListener()
	runTUI()
}

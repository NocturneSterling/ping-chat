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
	msgJson := ChatMessage{Message: msg, User: currentUser}
	jsonBytes, _ := json.Marshal(msgJson)
	hash := passHash(currentPass)
	sendBytes(append(hash, encryptToBytes(jsonBytes, []byte(*pass))...), *ip)
}

func runClientListener() {
	for {
		passwordHashBytes := passHash(currentPass)
		responseBytes := sendBytes(passwordHashBytes, *ip)
		responseStr := decryptFromBytes(responseBytes, passwordHashBytes)
		var response MsgRecord
		err := json.Unmarshal(responseStr, &response)
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		var incomingMsgJson ChatMessage
		msgTextStr := decryptUsingPass(response.LastMsgEncrypted, currentPass)
		json.Unmarshal([]byte(msgTextStr), &incomingMsgJson)
		if response.LastMsgTimestamp != lastTimestamp {
			if incomingMsgJson.Message == "" && incomingMsgJson.User == "" {
				tuiPrint("Chat begins here")
			} else {
				tuiPrint(incomingMsgJson.User + ": " + incomingMsgJson.Message)
			}
			lastTimestamp = response.LastMsgTimestamp
		}
		time.Sleep(1 * time.Second)
	}
}

func listenClient(){
	go runClientListener()
}

func runClient() {
	initTUI(runClientSender)
	//go runClientListener()
	runTUI()
}

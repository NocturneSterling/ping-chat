package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type ChatMessage struct {
	Message string `json:"msg"`
	User    string `json:"user"`
}

var lastTimestamp = map[string]int{}
var initalizedChannels = map[string]bool{}

func runClientSender(msg string) {
	chanNum := 0
	fmt.Sscanf(activeChannel,"channel %d",&chanNum)//parce activeChannel number to chanNum
	pass := channelPass(chanNum)
	msgJson := ChatMessage{Message: msg, User: currentUser}
	jsonBytes, _ := json.Marshal(msgJson)
	hash := passHash(pass)
	sendBytes(append(hash, encryptToBytes(jsonBytes, []byte(currentPass))...), *ip)
}

func runClientListener(chanNum int) {
	channel := currentChannel(chanNum)
	name := currentChannel(chanNum)
	pass := channelPass(chanNum)
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
		if response.LastMsgTimestamp != lastTimestamp[name] {
			if incomingMsgJson.Message == "" && incomingMsgJson.User == "" && !initalizedChannels[channel]{
				tuiPrint(channel,name,"Chat begins here")
				initalizedChannels[channel] = true
			} else {
				tuiPrint(channel,name,incomingMsgJson.User + ": " + incomingMsgJson.Message)
			}
			if response.LastMsgTimestamp != lastTimestamp[name]{//make match global
			lastTimestamp[name] = response.LastMsgTimestamp
			}
		}
		time.Sleep(300 * time.Millisecond)
	}
}

func listenClient(){//listens on all channels
	for i := 0; i < listNum; i++{
	go runClientListener(i)
	}
}

func runClient() {
	initTUI(runClientSender)
	//go runClientListener()
	runTUI()
}

package main

import (
	"encoding/json"
	"time"
)

type MsgRecord struct {
	LastMsgIp        string
	LastMsgEncrypted []byte
	LastMsgTimestamp int
}

var store = map[string]MsgRecord{}

// reply to any incoming pings
// gets incoming message ip and encrypted bytes from the client, returns what to reply with
func sendReply(ip string, incomingPayload []byte) []byte {
	key := extractHash(incomingPayload)
	if len(incomingPayload) == 32 {
		record, exists := store[string(key)]
		if !exists {
			store[string(key)] = MsgRecord{LastMsgIp: ip, LastMsgTimestamp: int(time.Now().Unix())}
		} else {
			_ = record
		}
		record = store[string(key)]
		recordJson, _ := json.Marshal(record)
		return encryptToBytes(recordJson, key)
	}

	encrypted := incomingPayload[32:]
	store[string(key)] = MsgRecord{
		LastMsgIp:        ip,
		LastMsgEncrypted: encrypted,
		LastMsgTimestamp: int(time.Now().Unix()),
	}
	record := store[string(key)]
	recordJson, _ := json.Marshal(record)
	return encryptToBytes(recordJson, key)
}

func runServer() {
	listenForPackets()
}

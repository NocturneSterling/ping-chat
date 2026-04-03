package main

import (
	"flag"
	"fmt"
)

var server = flag.Bool("server", false, "run as a server")

// reply to any incoming pings
// gets incoming message ip and text, returns what to reply with
func sendReply(ip string, text string) string {
	return ip + ", i got your " + text + ". thx"
}

func main() {
	flag.Parse()
	if *server {
		enableKernelReplies(false)
		listenForPackets()
		enableKernelReplies(true)
	} else {
		x := ""
		fmt.Scanln(&x)
		fmt.Println(sendString(x, "91.98.131.91"))
	}
}

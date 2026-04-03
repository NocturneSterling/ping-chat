package main

import (
	"flag"
)

var (
	server       = flag.Bool("server", false, "run as a server")
	reauthKernel = flag.Bool("reauthKernel", false, "give authorization back to kernel to handle pings")
)

func main() { //takes input and redirects to run either client or server
	mode := flag.String("mode", "client", "inputs server or client")
	flag.Parse()

	if *mode == "server"{
		runServer()
	} else {
		runClient()
	}
}

package main

import (
	"fmt"
	"github.com/Liam-E2/mainevent/eventsource"
	"log"
	"os"
)

func main() {
	// Get addr from env
	host, specified := os.LookupEnv("EVENTSOURCEHOST")
	if !specified {
		host = "localhost"
	}
	port, specified := os.LookupEnv("EVENTSOURCEPORT")
	if !specified {
		port = "9019"
	}
	addr := fmt.Sprintf("%s:%s", host, port)

	// Run server
	log.Fatal(eventsource.RunServer(addr))
}

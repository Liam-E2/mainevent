package main

import (
	"fmt"
	"github.com/Liam-E2/mainevent/eventsource"
	"log"
	"os"
)

// Here we don't parameterize the demo poller with anything
// But one could use Struct fields to parameterize custom pollers
type DemoPoller struct {
}

func (p DemoPoller) Poll(opts eventsource.PollerConfig) ([]byte, error){
	// Returns JSON-encoded bytes
	data := []byte("{\"static_example\": \"of polling data...\"}")
	return data, nil
}


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

	// Demo Poller Setup
	doneChan := make(chan bool)
	pollConfs := eventsource.PollerConfig{
		PollSeconds: 3, 
		DoneChan: doneChan, 
		EventName: "stream", 
		EventServerAddr: "http://"+addr}

	poller := DemoPoller{}

	go eventsource.Run(poller, pollConfs) // Run poller in background
	defer func(){doneChan <- true}() // defer turning off the poller
	
	// Run server
	log.Fatal(eventsource.RunServer(addr))
}

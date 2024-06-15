package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/Liam-E2/mainevent/eventsource"
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

	// Demo Poller Setup
	doneChan := make(chan bool)
	pollConfs := eventsource.PollerConfig{
		PollSeconds:     3,
		DoneChan:        doneChan,
		EventName:       "stream",
		EventServerAddr: "http://" + addr}

	poller := eventsource.HTTPPoller{
		Url:    "https://googlechromelabs.github.io/chrome-for-testing/last-known-good-versions-with-downloads.json",
		Header: map[string]string{"Content-Type": "application/json"},
		Method: "GET",
		Body:   bytes.NewReader(make([]byte, 0)),
	}

	go eventsource.Run(poller, pollConfs) // Run poller in background
	defer func() { doneChan <- true }()   // defer turning off the poller

	// Run server
	log.Fatal(eventsource.RunServer(addr))
}

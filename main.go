package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Liam-E2/mainevent/eventsource"
	"github.com/gin-gonic/gin"
)

func HeadersMiddleware() gin.HandlerFunc {
	// Requires headers to be set for SSE, including CORS in case needed in the future

	return func(c *gin.Context) {
		// SSE - content-type in method
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Transfer-Encoding", "chunked")

		// CORS
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Event-Name")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		// Pre-Flight request
		if c.Request.Method == "OPTIONS" {
			c.Status(http.StatusOK)
			return
		}
		c.Next()
	}
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
	defer func() { doneChan <- true }()

	// Run server
	log.Fatal(eventsource.RunServer(addr, HeadersMiddleware()))
}

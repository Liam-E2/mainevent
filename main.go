package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"bytes"
	"github.com/gin-gonic/gin"
)

// It keeps a list of clients those are currently attached
// and broadcasting events to those clients.
type Event struct {
	// Events are pushed to this channel by the main events-gathering routine
	Message chan string

	// New client connections
	NewClients chan chan string

	// Closed client connections
	ClosedClients chan chan string

	// Total client connections
	TotalClients map[chan string]bool
}

// New event messages are broadcast to all registered client connection channels
type ClientChan chan string

func main() {
	router := gin.Default()

	// Initialize new streaming server
	stream := NewServer()

	subscribe := router.Group("/events/subscribe")

	subscribe.GET("/stream", HeadersMiddleware(), stream.serveHTTP(), func(c *gin.Context) {
		v, ok := c.Get("clientChan")
		if !ok {
			return
		}
		clientChan, ok := v.(ClientChan)
		if !ok {
			return
		}
		c.Stream(func(w io.Writer) bool {
			// Stream message to client from message channel
			if msg, ok := <-clientChan; ok {
				c.SSEvent("message", msg)
				return true
			}
			return false
		})
	})

	publish := router.Group("/events/publish")
	publish.POST("/", func(c *gin.Context){
		// Pass body data directly as an SSE
		json_bytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			errResp := map[string]string{"error": "Could not read Body"}
			c.AbortWithStatusJSON(http.StatusBadRequest, errResp)
			return 
		}
		var json_data Question
		decoder := json.NewDecoder(bytes.NewBuffer(json_bytes))
		decoder.DisallowUnknownFields()
		err = decoder.Decode(&json_data) 
		if err != nil {
			errResp := map[string]string{"error": "Misformatted Question"}
			c.AbortWithStatusJSON(http.StatusBadRequest, errResp)
			return 
		}

		stream.Message <- string(json_bytes)
		c.String(http.StatusOK, "Gadzukes, there are %v connections", len(stream.TotalClients))
	})

	// Parse Static files
	router.StaticFile("/", "./index.html")

	router.Run(":8001")
}

// Initialize event and Start processing requests
func NewServer() (event *Event) {
	event = &Event{
		Message:       make(chan string),
		NewClients:    make(chan chan string),
		ClosedClients: make(chan chan string),
		TotalClients:  make(map[chan string]bool),
	}

	go event.listen()

	return
}

// It Listens all incoming requests from clients.
// Handles addition and removal of clients and broadcast messages to clients.
func (stream *Event) listen() {
	for {
		select {
		// Add new available client
		case client := <-stream.NewClients:
			stream.TotalClients[client] = true
			log.Printf("Client added. %d registered clients", len(stream.TotalClients))

		// Remove closed client
		case client := <-stream.ClosedClients:
			delete(stream.TotalClients, client)
			close(client)
			log.Printf("Removed client. %d registered clients", len(stream.TotalClients))

		// Broadcast message to client
		case eventMsg := <-stream.Message:
			for clientMessageChan := range stream.TotalClients {
				clientMessageChan <- eventMsg
				log.Printf("Pushed data to client. %d registered clients", len(stream.TotalClients))
			}
		}
	}
}

func (stream *Event) serveHTTP() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Initialize client channel
		clientChan := make(ClientChan)

		// Send new connection to event server
		stream.NewClients <- clientChan

		defer func() {
			// Send closed connection to event server
			stream.ClosedClients <- clientChan
		}()

		c.Set("clientChan", clientChan)

		c.Next()
	}
}

func HeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// SSE
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Transfer-Encoding", "chunked")

		// CORS
        c.Writer.Header().Set("Access-Control-Allow-Origin", "localhost")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {


        c.Next()
    }
}
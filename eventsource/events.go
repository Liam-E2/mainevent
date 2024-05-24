package eventsource

import (
	"log"
	"github.com/gin-gonic/gin"
)

// Struct to hold SSE event data
// Event should have a channel name/topic 
// and string data 
type SSEPub struct{
	Name string `json:"name"`
	Data string
}


// Wrapper for a chan string attached to some topic Name
type ClientChan struct {
	Chan chan string
	Name string
}


// Keeps a list of clients those are currently attached
// and broadcasting events to those clients.
type Event struct {
	// Events are pushed to this channel by the main events-gathering routine
	Message chan SSEPub

	// New client connections
	NewClients chan ClientChan

	// Closed client connections
	ClosedClients chan ClientChan

	// SSE name -> *ClientChan
	NamedClients map[string]map[ClientChan]bool
}

// Initialize event and Start processing requests
func NewServer() (event *Event) {
	event = &Event{
		Message:       make(chan SSEPub),
		NewClients:    make(chan ClientChan),
		ClosedClients: make(chan ClientChan),
		NamedClients:  make(map[string]map[ClientChan]bool), // map[channelName][channel][dummyBool]
	}

	go event.listen()

	return
}

// Listens to incoming events
// Handles addition and removal of clients and broadcast messages to clients.
func (stream *Event) listen() {
	for {
		select {
		// Add new available client
		case client := <-stream.NewClients:
			if stream.NamedClients[client.Name] == nil {
				stream.NamedClients[client.Name] = make(map[ClientChan]bool)
			}
			stream.NamedClients[client.Name][client] = true
			log.Printf("Client added. %d registered clients", len(stream.NamedClients[client.Name]))

		// Remove closed client from NamedClients map
		case client := <-stream.ClosedClients:
			delete(stream.NamedClients[client.Name], client)
			close(client.Chan)
			log.Printf("Removed client. %d registered clients", len(stream.NamedClients[client.Name]))

		// Broadcast message to client
		case eventMsg := <-stream.Message:
			for clientMessageChan := range stream.NamedClients[eventMsg.Name] { 
				clientMessageChan.Chan <- eventMsg.Data
				log.Printf("Pushed data to client. %d registered clients", len(stream.NamedClients[clientMessageChan.Name]))
				log.Printf("%v", eventMsg)
			}
		}
	}
}


func (stream *Event) serveHTTP() gin.HandlerFunc {
	// On receiving a subscription request -
	// Create new channel for client, pass on to event server, defer closing the client, continue on
	return func(c *gin.Context) {
		// Initialize client channel
		clientChan := make(chan string)
		name := c.Param("name")
		log.Printf("Received subscription request from %s", name)
		client := ClientChan{clientChan, name}

		// Send new connection to event server
		stream.NewClients <- client

		defer func() {
			// Send closed connection to event server
			stream.ClosedClients <- client
		}()

		c.Set("clientChan", client) // Client chan attached to context here

		c.Next()
	}
}

package eventsource

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type PollerError struct {
	msg string
}

func (p PollerError) Error() string {
	return p.msg
}

type PollerConfig struct {
	PollSeconds     int
	DoneChan        chan bool
	EventName       string
	EventServerAddr string
}

type Poller interface {
	Poll(opts PollerConfig) ([]byte, error)
}

func Run(p Poller, opts PollerConfig) error {

	if opts.PollSeconds == 0 {
		return PollerError{"Must set int PollTime > 0"}
	}

	go func(poller Poller, options PollerConfig) error {
		// Setup for publishing poll data - run outside loop
		client := http.Client{}
		postUrl := fmt.Sprintf("%s/events/publish", options.EventServerAddr)

		// Either close or poll
		for {
			select {
			case close := <-options.DoneChan:
				log.Printf("Closing polling loop %v", close)
				return nil
			default:
				// Run Poll method
				pollData, err := poller.Poll(options)
				if err != nil {
					log.Printf("\nError in poller: %s\n", err)
					time.Sleep(time.Duration(opts.PollSeconds) * time.Second)
					continue
				}

				// Build Request to event source server
				req, err := http.NewRequest("POST", postUrl, bytes.NewReader(pollData))
				if err != nil {
					log.Printf("Error building new request: %s", err)
					time.Sleep(time.Duration(opts.PollSeconds) * time.Second)
					continue
				}
				req.Header.Add("X-Event-Name", options.EventName)

				// Send Request
				resp, err := client.Do(req)
				if err != nil {
					log.Printf("Error sending request to event server: %s", err)
					time.Sleep(time.Duration(opts.PollSeconds) * time.Second)
					continue
				}

				// Read req body to fulfill net.http contract/log
				bodybytes, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Printf("Error reading response Body: %s", err)
					time.Sleep(time.Duration(opts.PollSeconds) * time.Second)
					continue
				}
				log.Printf("Event source response to publish: %s\n", string(bodybytes))
				resp.Body.Close()

				time.Sleep(time.Duration(opts.PollSeconds) * time.Second)
			}
		}
	}(p, opts)

	return nil
}

type HTTPPoller struct {
	// Simple implementation of poller interface to
	// Make an HTTP request every n seconds with specified
	// method, url, header map, and request body.
	// Returns bytes of response body, error
	Url    string
	Header map[string]string
	Method string
	Body   io.Reader
}

func (p HTTPPoller) Poll(opts PollerConfig) ([]byte, error) {
	client := http.Client{}
	req, err := http.NewRequest(p.Method, p.Url, p.Body)
	if err != nil {
		return nil, err
	}
	for key := range p.Header {
		req.Header.Add(key, p.Header[key])
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	out_bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return out_bytes, nil

}

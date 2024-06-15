package eventsource_test

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/Liam-E2/mainevent/eventsource"
)

type FakePoller struct {
	t *testing.T
}

func (p FakePoller) Poll(opts eventsource.PollerConfig) ([]byte, error) {
	p.t.Log("Polling!")
	client := http.Client{}
	data := make([]byte, 10)
	resp, err := client.Post(opts.EventServerAddr+"/test", "application/json", bytes.NewReader(data))
	if err != nil {
		fmt.Printf("%s", err)
		return nil, err
	}
	fmt.Printf("%v", resp)
	return data, nil
}

func FakeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%v\n\n", *r)
}

func TestPoller(t *testing.T) {
	// Create test poller
	var p eventsource.Poller = FakePoller{t}
	doneChan := make(chan bool)
	opts := eventsource.PollerConfig{1, doneChan, "stream", "http://localhost:9019"}

	// Run poller
	eventsource.Run(p, opts)

	// Create test server
	http.HandleFunc("/test", FakeHandler)
	go http.ListenAndServe("localhost:9019", nil)

	time.Sleep(10 * time.Second)
	doneChan <- true
}

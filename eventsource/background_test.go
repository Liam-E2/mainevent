package eventsource_test

import (
	"fmt"
	"bytes"
	"net/http"
	"testing"
	"time"

	"github.com/Liam-E2/mainevent/eventsource"
)

type FakePoller struct {
	t *testing.T
}

func (p FakePoller) Poll(opts eventsource.PollerConfig) error{
	p.t.Log("Polling!")
	client := http.Client{}
	data := bytes.NewReader(make([]byte, 10))
	resp, err := client.Post(opts.EventServerAddr+"/test", "application/json", data)
	if err != nil {
		fmt.Printf("%s", err)
		return err
	}
	fmt.Printf("%v", resp)
	return nil
}

func FakeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%v\n\n", *r)
}

func TestPoller(t *testing.T){
	// Create test server
	http.HandleFunc("/test", FakeHandler)
	go http.ListenAndServe("localhost:9019", nil)

	// Create test poller
	var p eventsource.Poller = FakePoller{t}
	doneChan := make(chan bool)
	opts := eventsource.PollerConfig{1, doneChan, "stream", "http://localhost:9019"}

	// Run poller
	eventsource.Run(p, opts)
	time.Sleep(10 * time.Second)
	doneChan <- true
}
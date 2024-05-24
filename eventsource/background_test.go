package eventsource_test

import (
	"testing"
	"time"
	"github.com/Liam-E2/mainevent/eventsource"

)

type FakePoller struct {
	t *testing.T
}

func (p FakePoller) Poll(opts eventsource.PollerOpts) error{
	p.t.Log("Polling!")
	return nil
}

func TestPoller(t *testing.T){
	var p eventsource.Poller = FakePoller{t}
	doneChan := make(chan bool)
	opts := eventsource.PollerOpts{1, doneChan}
	go eventsource.Run(p, opts)
	time.Sleep(10 * time.Second)
	doneChan <- true
}
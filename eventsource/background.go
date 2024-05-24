package eventsource

import (
	"log"
	"time"
)

type PollerError struct {
	msg string
}

func (p PollerError) Error() string {
	return p.msg
}

type PollerConfig struct {
	PollSeconds int
	DoneChan chan bool
	EventName string
	EventServerAddr string
}


type Poller interface {
	Poll(opts PollerConfig) error
}

func Run(p Poller, opts PollerConfig) error {

	if opts.PollSeconds == 0 {
		return PollerError{"Must set int PollTime > 0"}
	}

	go func(poller Poller, options PollerConfig) error{
		for {
			select {
			case close := <- options.DoneChan:
				log.Printf("Closing polling loop %v", close)
				return nil
			default:
				poller.Poll(options)
				time.Sleep(time.Duration(opts.PollSeconds) * time.Second)
			}
		}
	}(p, opts)

	return nil
}

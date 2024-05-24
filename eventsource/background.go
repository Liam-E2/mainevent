package eventsource

import (
	"time"
	"log"
)

type PollerError struct {
	msg string
}

func (p PollerError) Error() string {
	return p.msg
}

type PollerOpts struct {
	PollSeconds int
	DoneChan chan bool
}

type Poller interface {
	Poll(opts PollerOpts) error
}

func Run(p Poller, opts PollerOpts) error {

	if opts.PollSeconds == 0 {
		return PollerError{"Must set int PollTime > 0"}
	}

	go func(poller Poller, options PollerOpts) error{
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

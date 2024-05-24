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
	PollTime time.Duration
	DoneChan chan bool
}

type Poller interface {
	Poll(opts PollerOpts) error
}

func Run(p Poller, opts PollerOpts) error {

	if opts.PollTime == 0 {
		return PollerError{"Must set valid pollTime"}
	}

	go func(poller Poller, options PollerOpts) error{
		for {
			select {
			case close := <- options.DoneChan:
				log.Printf("Closing polling loop %v", close)
				return nil
			default:
				poller.Poll(options)
				time.Sleep(opts.PollTime)
			}
		}
	}(p, opts)

	return nil
}

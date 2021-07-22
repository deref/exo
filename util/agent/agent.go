package agent

import (
	"context"
	"errors"
)

// Agent guards asynchronous changes of an individual location.
// Inspired by <https://clojure.org/reference/agents>.
// Note: This is not quite a Clojure agent, since there is no way to deref.
type Agent struct {
	inbox chan action
	exit  chan struct{}
	fail  chan error
}

type action struct {
	f     func() error
	reply chan error
}

func NewAgent(size int) *Agent {
	return &Agent{
		inbox: make(chan action, size),
		exit:  make(chan struct{}),
		fail:  make(chan error, 1),
	}
}

// Process messages until either an error occurs or the context is canceled.
func (a *Agent) Run(ctx context.Context) error {
	var err error
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case err = <-a.fail:
			break loop
		case action := <-a.inbox:
			action.reply <- action.f()
		}
	}
	// Drain queue.
	close(a.exit)
	close(a.inbox)
	for action := range a.inbox {
		action.reply <- context.Canceled
	}
	return err
}

// Record an asynchronous error, causing the Run loop to exit.
func (a *Agent) Fail(err error) {
	select {
	case a.fail <- err:
	default:
		// Discard secondary failure.
	}
}

var InboxFull = errors.New("inbox full")

// Queues f for execution on the agent's thread and returns f's error value.
// If the agent is running slowly, InboxFull is returned immediately.
// If the agent is shutdown, f is not run and context.Canceled is returned.
func (a *Agent) Send(f func() error) error {
	reply := make(chan error, 1)
	select {
	case <-a.exit:
		reply <- context.Canceled
	case a.inbox <- action{
		f:     f,
		reply: reply,
	}:
	default:
		reply <- InboxFull
	}
	return <-reply
}

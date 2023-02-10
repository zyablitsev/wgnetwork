package shutdown

import (
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Shutdown holds initiate and done channels and timeout variable
type Shutdown struct {
	syscallChan chan os.Signal
	signalChan  chan struct{}
	doneChan    <-chan struct{}
	timeout     time.Duration
	log         Logger
}

// New is a constructor for Shutdown object
func New(
	log Logger,
	doneChan <-chan struct{},
	timeout time.Duration,
) *Shutdown {
	s := &Shutdown{
		syscallChan: make(chan os.Signal),
		signalChan:  make(chan struct{}),
		doneChan:    doneChan,
		timeout:     timeout,
		log:         log,
	}

	return s
}

// Signal returns chan
func (s *Shutdown) Signal() <-chan struct{} {
	return s.signalChan
}

// Wait signal to shutdown.
func (s *Shutdown) Wait() {
	signal.Notify(s.syscallChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-s.syscallChan:
		s.wait()
	case <-s.doneChan:
		s.exit("shutdown")
	}
}

func (s *Shutdown) wait() {
	s.log.Info("interrupt syscall received")
	s.signalChan <- struct{}{}
	select {
	case <-s.syscallChan:
		s.exit("repeated interrupt syscall received, forced shutdown")
	case <-s.doneChan:
		s.exit("cleanup done, shutdown")
	case <-time.After(s.timeout):
		s.exit("cleanup timed out, shutdown")
	}
}

func (s *Shutdown) exit(msg string) {
	s.log.Info(msg)
	close(s.signalChan)
	close(s.syscallChan)
}

// Logger represents log object interface
type Logger interface {
	Info(a ...interface{})
}

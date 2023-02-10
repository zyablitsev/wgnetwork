package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"wgnetwork"
	"wgnetwork/pkg/logger"
	"wgnetwork/pkg/shutdown"
)

func main() {
	log, err := logger.New(os.Stdout, os.Stderr)
	if err != nil {
		panic(fmt.Errorf("can't init logger: %v", err))
	}

	// create context for service
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s, err := wgnetwork.Init(ctx, cancel, log)
	if err != nil {
		panic(fmt.Errorf("can't init service: %v", err))
	}

	go s.Run()

	exit := shutdown.New(log, ctx.Done(), 10*time.Second)
	go func() { // graceful shutdown
		<-exit.Signal()
		s.Stop()
	}()
	exit.Wait()
}

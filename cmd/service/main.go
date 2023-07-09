package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"wgnetwork"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	s, err := wgnetwork.Init(ctx)
	if err != nil {
		log.Fatal(err)
	}

	s.Run()
}

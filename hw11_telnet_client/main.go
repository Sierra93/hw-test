package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"
)

var timeout time.Duration

func init() {
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "connection timeout")
}

func main() {
	flag.Parse()

	args := os.Args[2:]
	if len(args) < 2 {
		fmt.Printf("Invalid arguments number: want - 2, got - %d", len(args))
		os.Exit(1)
	}

	address := net.JoinHostPort(args[0], args[1])
	tc := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

	if err := tc.Connect(); err != nil {
		fmt.Printf("%s: failed to connect: %v", address, err)
		os.Exit(1)
	}
	defer tc.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		tc.Send()
		stop()
	}()

	go func() {
		tc.Receive()
		stop()
	}()

	<-ctx.Done()
}

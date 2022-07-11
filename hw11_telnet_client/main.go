package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var timeout time.Duration

func init() {
	flag.DurationVar(&timeout, "timeout", time.Second*10, "connect timeout")
}

func main() {
	flag.Parse()

	args := flag.Args()

	if len(args) != 2 {
		log.Fatal("invalid arguments")
	}

	address := net.JoinHostPort(args[0], args[1])

	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	defer client.Close()

	err := client.Connect()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(2)
	// receiver func
	go func() {
		for {
			select {
			case <-ctx.Done():
				client.Close()

				wg.Done()
				return
			default:
				err := client.Receive()
				if err != nil {
					os.Exit(1)
				}
			}
		}
	}()

	// sender func
	go func() {
		for {
			select {
			case <-ctx.Done():
				client.Close()
				wg.Done()
				return
			default:
				err := client.Send()
				if err != nil {
					os.Exit(1)
				}
			}
		}
	}()

	wg.Wait()
}

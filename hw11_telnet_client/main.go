package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

var (
	ErrorIncorrectUsage  = errors.New("usage: go-telnet [--timeout=10s] host port")
	ErrorInterruptSignal = errors.New("Bye-bye")
)

func main() {
	timeout := 10 * time.Second
	var host string

	if len(os.Args) < 2 {
		fmt.Println(ErrorIncorrectUsage)
		os.Exit(1)
	}

	args := os.Args[1:]
	if after, ok := strings.CutPrefix(args[0], "--timeout="); ok {
		t, _ := time.ParseDuration(after)
		timeout = t
		args = args[1:]
	}

	if len(args) < 2 {
		fmt.Println(ErrorIncorrectUsage)
		os.Exit(1)
	}

	host = net.JoinHostPort(args[0], args[1])

	telnetClient := NewTelnetClient(host, timeout, os.Stdin, os.Stdout)

	if err := telnetClient.Connect(); err != nil {
		fmt.Println("Connection error:", err)
		os.Exit(1)
	}
	defer telnetClient.Close()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)

	errChan := make(chan error, 1)
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		errChan <- telnetClient.Send()
	}()

	go func() {
		defer wg.Done()
		errChan <- telnetClient.Receive()
	}()

	select {
	case err := <-errChan:
		if err != nil {
			fmt.Println(err)
			telnetClient.Close()
			if errors.Is(io.EOF, err) {
				os.Exit(0)
			}
			os.Exit(1)
		}
	case <-sigChan:
		fmt.Println(ErrorInterruptSignal)
		telnetClient.Close()
		os.Exit(0)
	}

	wg.Wait()
}

package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

var ErrorFewArguments = errors.New("not enougth arguments")
var ErrorBadArguments = errors.New("incorrect arguments")

func main() {
	args := os.Args
	if len(args) < 3 {
		fmt.Println(ErrorFewArguments)
		os.Exit(1)
	}
	t, _ := time.ParseDuration("10s")
	if strings.Contains(args[1], "--timeout=") {
		if len(args) < 4 {
			fmt.Println(ErrorBadArguments)
			os.Exit(1)
		}
		t, _ = time.ParseDuration(strings.ReplaceAll(args[1], "--timeout=", ""))
	}
	host := strings.Join(args[2:4], ":")
	telnetClient := NewTelnetClient(host, t, os.Stdout, os.Stdin)
	telnetClient.Connect()

	// Place your code here,
	// P.S. Do not rush to throw context down, think think if it is useful with blocking operation?
}

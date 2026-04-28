package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 3 {
		fmt.Println("Not anought arguments")
		os.Exit(1)
	}
	environment, err := ReadDir(args[1])
	if err != nil {
		fmt.Println("error while reading env directory. ", err)
		os.Exit(1)
	}
	os.Exit(RunCmd(args[1:], environment))
}

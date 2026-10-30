package main

import (
	"fmt"
	"os"
	"treenity/src/server"
)

func main() {
	argv := os.Args
	argc := len(argv)

	if argc > 1 {
		fmt.Fprintln(os.Stderr, "Error: Invalid command. No arguments needed")
		os.Exit(1)
	}

	s := server.NewServer()
	err := s.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: ", err.Error())
		os.Exit(1)
	}
}

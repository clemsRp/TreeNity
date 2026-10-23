package main

import (
	"fmt"
	"os"
	"treenity/src/client"
)

func main() {
	// Get args
	argv := os.Args
	argc := len(argv)

	if argc < 3 {
		fmt.Fprintln(os.Stderr, "Error: Invalid number of parameters")
		os.Exit(1)
	}

	entrypoint := argv[1]

	// Create client
	c := client.NewClient()

	// Connect to server
	err := c.Connect(entrypoint)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: ", err.Error())
		os.Exit(1)
	}
	defer c.Disconnect()

	// Handle command
	err = c.HandleCommand()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: ", err.Error())
		os.Exit(1)
	}
}

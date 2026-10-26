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

	// Get kong context
	ctx, kong_err := c.GetContext(os.Args[1:])
	if kong_err != nil {
		fmt.Fprintln(os.Stderr, "Error:", kong_err.Error())
		os.Exit(1)
	}

	// Connect to server
	err := c.Connect(entrypoint)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err.Error())
		os.Exit(3)
	}
	defer c.Disconnect()

	// Handle command
	cmd_err := c.HandleCommand(ctx)
	if cmd_err != nil {
		fmt.Fprintln(os.Stderr, "Error:", cmd_err.Message)
		os.Exit(cmd_err.Code)
	}
}

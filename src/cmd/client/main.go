package main

import (
	"fmt"
	"os"
	"treenity/src/client"
)

func main() {
	os.Exit(run())
}

func run() int {
	// Get args
	argv := os.Args
	argc := len(argv)

	if argc < 3 {
		fmt.Fprintln(os.Stderr, "Error: Invalid number of parameters")
		return 1
	}

	entrypoint := argv[1]

	// Create client
	c := client.NewClient()

	// Get kong context
	ctx, kong_err := c.GetContext(os.Args[1:])
	if kong_err != nil {
		fmt.Fprintln(os.Stderr, "Error:", kong_err.Error())
		return 1
	}

	// Connect to server
	err := c.Connect(entrypoint)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err.Error())
		return 1
	}
	defer c.Disconnect()

	// Handle command
	cmd_err := c.HandleCommand(ctx)
	if cmd_err != nil {
		fmt.Fprintln(os.Stderr, "Error:", cmd_err.Message)
		return cmd_err.Code
	}

	return 0
}

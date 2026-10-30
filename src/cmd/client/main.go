package main

import (
	"fmt"
	"os"
	"treenity/src/client"
)

func main() {
	argv := os.Args
	argc := len(argv)

	if argc < 3 {
		fmt.Fprintln(os.Stderr, "Error: Invalid number of parameters")
		os.Exit(1)
	}

	entrypoint := argv[1]

	c := client.NewClient()

	ctx, kong_err := c.GetContext(os.Args[1:])
	if kong_err != nil {
		fmt.Fprintln(os.Stderr, "Error:", kong_err.Error())
		os.Exit(1)
	}

	err := c.Connect(entrypoint)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err.Error())
		os.Exit(3)
	}
	defer c.Disconnect()

	cmd_err := c.HandleCommand(ctx)
	if cmd_err != nil {
		fmt.Fprintln(os.Stderr, "Error:", cmd_err.Message)
		os.Exit(cmd_err.Code)
	}
}

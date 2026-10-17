package cmd

import (
	"fmt"
	pr "treenity/src/protocol"

	"github.com/alecthomas/kong"
)

func HandleCmdList(ctx *kong.Context) error {
	list_struct := pr.CLI.IpcIdentifier.List
	fmt.Printf("IPC: '%s'\n", list_struct.IpcIdentifier)

	// TODO: Implemente behaviour

	return nil
}

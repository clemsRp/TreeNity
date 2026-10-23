package cmd

import (
	"fmt"
	pr "treenity/src/protocol"
	"treenity/src/protocol/ipc"

	"github.com/alecthomas/kong"
)

func HandleCmdList(ctx *kong.Context, writer *ipc.Writer) error {
	list_struct := pr.CLI.IpcIdentifier.List
	fmt.Printf("IPC: '%s'\n", list_struct.IpcIdentifier)

	// TODO: Implemente behaviour

	return nil
}

package cmd

import (
	"fmt"
	pr "treenity/src/protocol"
	"treenity/src/protocol/ipc"

	"github.com/alecthomas/kong"
)

func HandleCmdProduce(ctx *kong.Context, writer *ipc.Writer) error {
	pro_struct := pr.CLI.IpcIdentifier.Produce
	fmt.Printf("IPC: '%s'\n", pro_struct.IpcIdentifier)
	fmt.Printf("Topic: '%s'\n", pro_struct.TopicName)
	fmt.Printf("Raw: '%t'\n", pro_struct.Raw)

	// TODO: Implemente behaviour

	return nil
}

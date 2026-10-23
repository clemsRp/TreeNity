package cmd

import (
	"fmt"
	pr "treenity/src/protocol"
	"treenity/src/protocol/ipc"

	"github.com/alecthomas/kong"
)

func HandleCmdCreate(ctx *kong.Context, writer *ipc.Writer) error {
	cre_struct := pr.CLI.IpcIdentifier.Create
	fmt.Printf("IPC: '%s'\n", cre_struct.IpcIdentifier)
	fmt.Printf("Topic: '%s'\n", cre_struct.TopicName)

	// TODO: Implemente behaviour
	writer.WriteLine([]byte("salutsalut"))

	return nil
}

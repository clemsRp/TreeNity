package cmd

import (
	"fmt"
	pr "treenity/src/protocol"

	"github.com/alecthomas/kong"
)

func HandleCmdCreate(ctx *kong.Context) error {
	cre_struct := pr.CLI.IpcIdentifier.Create
	fmt.Printf("IPC: '%s'\n", cre_struct.IpcIdentifier)
	fmt.Printf("Topic: '%s'\n", cre_struct.TopicName)

	// TODO: Implemente behaviour

	return nil
}

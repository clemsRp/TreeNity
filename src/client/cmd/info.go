package cmd

import (
	"fmt"
	pr "treenity/src/protocol"

	"github.com/alecthomas/kong"
)

func HandleCmdInfo(ctx *kong.Context) error {
	info_struct := pr.CLI.IpcIdentifier.Info
	fmt.Printf("IPC: '%s'\n", info_struct.IpcIdentifier)
	fmt.Printf("Subscriber: '%s'\n", info_struct.SubscriberName)

	// TODO: Implemente behaviour

	return nil
}

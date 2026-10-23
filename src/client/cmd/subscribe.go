package cmd

import (
	"fmt"
	pr "treenity/src/protocol"
	"treenity/src/protocol/ipc"

	"github.com/alecthomas/kong"
)

func HandleCmdSubscribe(ctx *kong.Context, writer *ipc.Writer) error {
	sub_struct := pr.CLI.IpcIdentifier.Subscribe
	fmt.Printf("IPC: '%s'\n", sub_struct.IpcIdentifier)
	fmt.Printf("Topic: '%s'\n", sub_struct.TopicName)
	fmt.Printf("Subscriber: '%s'\n", sub_struct.SubscriberName)
	fmt.Printf("Raw: '%t'\n", sub_struct.Raw)
	fmt.Printf("Prefix: '%s'\n", *sub_struct.Prefix)
	fmt.Printf("Offset: '%d'\n", *sub_struct.Offset)

	// TODO: Implemente behaviour

	return nil
}

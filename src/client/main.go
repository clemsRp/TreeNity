package client

import (
	"errors"
	"os"
	"treenity/src/client/cmd"
	pr "treenity/src/protocol"

	"github.com/alecthomas/kong"
)

type Client struct {
	Name string
}

func NewClient() *Client {
	return &Client{
		Name: "Clément",
	}
}

func (c *Client) Connect(entrypoint string) error {
	return nil
}

func (c *Client) HandleCommand() error {
	// Get kong context
	ctx, err := c.GetContext(os.Args[1:])
	if err != nil {
		return err
	}

	// Handle valid commands
	var res error
	i_struct := &pr.CLI.IpcIdentifier
	ipc := i_struct.IpcIdentifier
	command := ctx.Selected().Name

	switch command {
	// Create
	case pr.CmdCreate:
		i_struct.Create.IpcIdentifier = ipc
		res = cmd.HandleCmdCreate(ctx)

	// List
	case pr.CmdList:
		i_struct.List.IpcIdentifier = ipc
		res = cmd.HandleCmdList(ctx)

	// Produce
	case pr.CmdProduce:
		i_struct.Produce.IpcIdentifier = ipc
		res = cmd.HandleCmdProduce(ctx)

	// Subscribe
	case pr.CmdSubscribe:
		i_struct.Subscribe.IpcIdentifier = ipc
		res = cmd.HandleCmdSubscribe(ctx)

	// Info
	case pr.CmdInfo:
		i_struct.Info.IpcIdentifier = ipc
		res = cmd.HandleCmdInfo(ctx)

	// Handle invalid command
	default:
		res = errors.New(
			"Invalid command. Need to be: create, list, produce, subscribe or info.",
		)
	}

	return res
}

func (c *Client) GetContext(args []string) (*kong.Context, error) {
	// Create kong parser
	parser, err := kong.New(&pr.CLI,
		kong.Name("client"),
		kong.Description("Treenity Message Queue System Client"),
	)

	// Handle initialization error
	if err != nil {
		return nil, err
	}

	// Parse arguments
	ctx, err := parser.Parse(args)
	if err != nil {
		return nil, err
	}

	return ctx, nil
}

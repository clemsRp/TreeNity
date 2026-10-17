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
	command := ctx.Selected().Name
	switch command {
	case pr.CmdCreate:
		res = cmd.HandleCmdCreate(ctx)
	case pr.CmdList:
		res = cmd.HandleCmdList(ctx)
	case pr.CmdProduce:
		res = cmd.HandleCmdProduce(ctx)
	case pr.CmdSubscribe:
		res = cmd.HandleCmdSubscribe(ctx)
	case pr.CmdInfo:
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

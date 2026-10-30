package client

import (
	pr "treenity/src/protocol"

	"github.com/alecthomas/kong"
)

func (c *Client) GetContext(args []string) (*kong.Context, error) {
	parser, err := kong.New(&pr.CLI,
		kong.Name("client"),
		kong.Description("Treenity Message Queue System Client"),
	)

	if err != nil {
		return nil, err
	}

	ctx, err := parser.Parse(args)
	if err != nil {
		return nil, err
	}

	return ctx, nil
}

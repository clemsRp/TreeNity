package client

import (
	"errors"
	pr "treenity/src/protocol"

	"github.com/alecthomas/kong"
)

func (c *Client) HandleCommand(ctx *kong.Context) *pr.Error {
	i_struct := &pr.CLI.IpcIdentifier
	ipc := i_struct.IpcIdentifier
	command := ctx.Selected().Name

	cmd_struct, cmd_err := c.GetCmdStruct(i_struct, ipc, command)
	if cmd_err != nil {
		return &pr.Error{
			Message: cmd_err.Error(),
			Code:    1,
		}
	}

	req_err := c.SendRequest(cmd_struct, command)
	if req_err != nil {
		return req_err
	}

	switch command {
	case pr.CmdProduce:
		cmd_product := cmd_struct.(pr.CmdProduceDatas)
		if pr.CLI.IpcIdentifier.Produce.Raw {
			return c.ExecProduceRaw(cmd_product.TopicName)
		}
		return c.ExecProduceText(cmd_product.TopicName)

	case pr.CmdSubscribe:
		return c.ExecSubscribe(cmd_struct)
	}

	return nil
}

func (c *Client) GetCmdStruct(i_struct *pr.IpcArgs, ipc, command string) (any, error) {
	var cmd_struct any

	switch command {
	case pr.CmdCreate:
		i_struct.Create.IpcIdentifier = ipc
		cmd_struct = i_struct.Create

	case pr.CmdList:
		i_struct.List.IpcIdentifier = ipc
		cmd_struct = i_struct.List

	case pr.CmdProduce:
		i_struct.Produce.IpcIdentifier = ipc
		cmd_struct = i_struct.Produce

	case pr.CmdSubscribe:
		i_struct.Subscribe.IpcIdentifier = ipc
		cmd_struct = i_struct.Subscribe

	case pr.CmdInfo:
		i_struct.Info.IpcIdentifier = ipc
		cmd_struct = i_struct.Info

	default:
		return nil, errors.New(pr.ErrInvalidCommand)
	}

	return cmd_struct, nil
}

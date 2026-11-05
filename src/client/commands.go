package client

import (
	"errors"
	pr "treenity/src/protocol"
	"treenity/src/protocol/ipc"

	"github.com/alecthomas/kong"
)

func (c *Client) HandleCommand(ctx *kong.Context) *pr.Error {
	// Basics variables
	i_struct := &pr.CLI.IpcIdentifier
	ipcPath := i_struct.IpcIdentifier
	command := ctx.Selected().Name

	cmd_struct, cmd_err := c.GetCmdStruct(i_struct, ipcPath, command)
	if cmd_err != nil {
		return &pr.Error{
			Message: cmd_err.Error(),
			Code:    1,
		}
	}

	var permReader *ipc.Reader
	if command == pr.CmdSubscribe {
		subCmd := cmd_struct.(pr.CmdSubscribeDatas)
		permanentPath := subCmd.IpcIdentifier + "." + subCmd.SubscriberName
		var err error
		permReader, err = ipc.NewReader(permanentPath)
		if err != nil {
			return &pr.Error{
				Message: err.Error(),
				Code:    1,
			}
		}
	}

	// Send command request
	req_err := c.SendRequest(cmd_struct, command)
	if req_err != nil {
		if permReader != nil {
			permReader.Close()
		}
		return req_err
	}

	// Handle Produce/Subscribe Input/Ouput
	switch command {
	case pr.CmdProduce:
		cmd_product := cmd_struct.(pr.CmdProduceDatas)
		return c.Run(cmd_product)

	case pr.CmdSubscribe:
		return c.ExecSubscribe(cmd_struct, permReader)
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

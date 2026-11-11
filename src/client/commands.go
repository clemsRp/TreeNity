package client

import (
	"errors"
	"strconv"
	"strings"
	pr "treenity/src/protocol"

	"github.com/alecthomas/kong"
)

func (c *Client) HandleCommand(ctx *kong.Context) *pr.Error {
	var err error

	// Basics variables
	i_struct := &pr.CLI.IpcIdentifier
	ipcPath := i_struct.IpcIdentifier
	command := ctx.Selected().Name

	// Get server PID
	if strings.Contains(ipcPath, pr.EntryPoint) {
		serverPid := strings.SplitN(ipcPath, pr.EntryPoint, 2)[1]
		c.ServerPid, err = strconv.ParseInt(serverPid, 10, 64)
		if err != nil {
			return &pr.Error{
				Message: err.Error(),
				Code:    1,
			}
		}
	}

	cmd_struct, cmd_err := c.GetCmdStruct(i_struct, ipcPath, command)
	if cmd_err != nil {
		return &pr.Error{
			Message: cmd_err.Error(),
			Code:    1,
		}
	}

	// Send command request
	req_err := c.SendRequest(cmd_struct, command)
	if req_err != nil {
		if c.Reader != nil {
			c.Reader.Close()
		}
		return req_err
	}

	// Handle Produce/Subscribe Input/Ouput
	switch command {
	case pr.CmdProduce:
		cmd_product := cmd_struct.(pr.CmdProduceDatas)
		return c.Run(cmd_product)

	case pr.CmdSubscribe:
		cmd_subs := cmd_struct.(pr.CmdSubscribeDatas)
		c.Name = cmd_subs.SubscriberName

		if cmd_subs.Offset != nil && *cmd_subs.Offset < 0 {
			return &pr.Error{
				Message: "offset cannot be a negative integer",
				Code:    1,
			}
		}

		return c.ExecSubscribe(cmd_subs)
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

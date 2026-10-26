package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	pr "treenity/src/protocol"

	"github.com/alecthomas/kong"
)

func (c *Client) HandleCommand(ctx *kong.Context) *pr.Error {
	// Basics variables
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

	return c.SendRequest(cmd_struct, command)
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

func (c *Client) SendRequest(cmd_struct any, command string) *pr.Error {
	// Convert cmd_struct to Raw Message
	rawDatas, marsh_err := json.Marshal(cmd_struct)
	if marsh_err != nil {
		return &pr.Error{
			Message: marsh_err.Error(),
			Code:    1,
		}
	}

	// Init request
	request := pr.ToServer{
		Command: command,
		Datas:   json.RawMessage(rawDatas),
		Path:    c.Path,
	}

	// Cast request
	payload, err := json.Marshal(request)
	if err != nil {
		return &pr.Error{
			Message: err.Error(),
			Code:    1,
		}
	}

	// Send request to Server
	err = c.Writer.WriteLine(payload)
	if err != nil {
		return &pr.Error{
			Message: err.Error(),
			Code:    3,
		}
	}

	// Wait for response
	line, read_err := c.Reader.ReadLine()
	if read_err != nil {
		if read_err == io.EOF {
			return nil
		}
		return &pr.Error{
			Message: read_err.Error(),
			Code:    3,
		}
	}

	// Cast response
	var response pr.ToClient
	json_err := json.Unmarshal(line, &response)
	if json_err != nil {
		return &pr.Error{
			Message: json_err.Error(),
			Code:    1,
		}
	}

	// Handle error
	if response.Err != nil {
		return response.Err
	}

	// Print response
	res := response.Response
	if string(response.Datas) != `""` {
		res += string(response.Datas)
	}

	fmt.Println(res)

	return nil
}

package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	pr "treenity/src/protocol"
)

func (c *Client) HandleCommand() error {
	// Get kong context
	ctx, err := c.GetContext(os.Args[1:])
	if err != nil {
		return err
	}

	// Basics variables
	i_struct := &pr.CLI.IpcIdentifier
	ipc := i_struct.IpcIdentifier
	command := ctx.Selected().Name

	cmd_struct, err := c.GetCmdStruct(i_struct, ipc, command)

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

func (c *Client) SendRequest(cmd_struct any, command string) error {
	// Convert cmd_struct to Raw Message
	rawDatas, err := json.Marshal(cmd_struct)
	if err != nil {
		return err
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
		return err
	}

	// Send request to Server
	err = c.Writer.WriteLine(payload)
	if err != nil {
		return err
	}

	// Wait for response

	line, err := c.Reader.ReadLine()
	if err != nil {
		if err == io.EOF {
			return nil
		}
		return err
	}

	// Cast response
	var response pr.ToClient
	json_err := json.Unmarshal(line, &response)
	if json_err != nil {
		return json_err
	}

	fmt.Println(response.Response, response.Err)

	return nil
}

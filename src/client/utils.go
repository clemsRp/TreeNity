package client

import (
	"encoding/json"
	"fmt"
	"io"
	pr "treenity/src/protocol"
)

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

	if len(res) > 0 {
		fmt.Println(res)
	}

	return nil
}

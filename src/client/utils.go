package client

import (
	"encoding/json"
	"fmt"
	"io"
	pr "treenity/src/protocol"
)

func (c *Client) SendRequest(cmd_struct any, command string) *pr.Error {
	rawDatas, marsh_err := json.Marshal(cmd_struct)
	if marsh_err != nil {
		return &pr.Error{
			Message: marsh_err.Error(),
			Code:    1,
		}
	}

	request := pr.ToServer{
		Command: command,
		Datas:   json.RawMessage(rawDatas),
		Path:    c.Path,
	}

	payload, err := json.Marshal(request)
	if err != nil {
		return &pr.Error{
			Message: err.Error(),
			Code:    1,
		}
	}

	err = c.Writer.WriteLine(payload)
	if err != nil {
		return &pr.Error{
			Message: err.Error(),
			Code:    3,
		}
	}

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

	var response pr.ToClient
	json_err := json.Unmarshal(line, &response)
	if json_err != nil {
		return &pr.Error{
			Message: json_err.Error(),
			Code:    1,
		}
	}

	if response.Err != nil {
		return response.Err
	}

	res := response.Response
	if string(response.Datas) != `""` {
		res += string(response.Datas)
	}

	if len(res) > 0 {
		fmt.Println(res)
	}

	return nil
}

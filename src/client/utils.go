package client

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"syscall"
	pr "treenity/src/protocol"
)

// Global utils
func (c *Client) Send(cmd_struct any, command string) *pr.Error {
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

	return nil
}

func (c *Client) SendRequest(cmd_struct any, command string) *pr.Error {
	// Send request
	send_err := c.Send(cmd_struct, command)
	if send_err != nil {
		return send_err
	}

	// Wait for response
	line, read_err := c.Reader.ReadLine()
	if read_err != nil {
		if read_err == io.EOF {
			return &pr.Error{
				Message: "server closed connection unexpectedly",
				Code:    3,
			}
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

// Subscribe utils
func PrintData(dataMsg pr.CmdMessageDatas, cmd pr.CmdSubscribeDatas) *pr.Error {
	// Handle text mode
	if !cmd.Raw {
		fmt.Printf("%s:%s\n", string(dataMsg.Key), string(dataMsg.Body))
		return nil
	}

	// Init buffer
	buffer := new(bytes.Buffer)

	// Print offset
	off_err := binary.Write(buffer, binary.LittleEndian, int32(dataMsg.Offset))
	if off_err != nil {
		return &pr.Error{
			Message: off_err.Error(),
			Code:    1,
		}
	}

	// Write datas on buffer
	for _, data := range [][]byte{dataMsg.Key, dataMsg.Body} {
		// Write data len
		if err := binary.Write(buffer, binary.LittleEndian, int32(len(data))); err != nil {
			return &pr.Error{
				Message: err.Error(),
				Code:    1,
			}
		}

		// Write data content
		buffer.Write(data)
	}

	// Write buffer on output
	_, write_err := os.Stdout.Write(buffer.Bytes())
	if write_err != nil {
		return &pr.Error{
			Message: write_err.Error(),
			Code:    1,
		}
	}

	return nil
}

func (c *Client) ReadData(lineChan chan []byte, errChan chan error) {
	go func() {
		for {
			data_read, err := c.Reader.ReadLine()
			if err != nil {
				errChan <- err
				return
			}
			lineChan <- data_read
		}
	}()
}

func ProcessExists(serverPid int64) bool {
	// Handle invalid PID
	if serverPid <= 0 {
		return false
	}

	err := syscall.Kill(int(serverPid), 0)

	return !errors.Is(err, syscall.ESRCH)
}

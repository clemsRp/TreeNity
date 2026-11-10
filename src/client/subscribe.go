package client

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	pr "treenity/src/protocol"
)

func PrintData(dataMsg pr.CmdMessageDatas, cmd pr.CmdSubscribeDatas) *pr.Error {
	// Handle text mode
	if !cmd.Raw {
		fmt.Printf("%s:%s\n", dataMsg.Key, dataMsg.Body)
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
	for _, data := range [][]byte{dataMsg.RawKey, dataMsg.RawBody} {
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

func (c *Client) ExecSubscribe(cmd_struct any) *pr.Error {
	cmd := cmd_struct.(pr.CmdSubscribeDatas)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		<-ctx.Done()
		c.Reader.Close()
	}()

	for {
		data_read, err := c.Reader.ReadLine()

		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			if errors.Is(err, io.EOF) {
				fmt.Println("Server closed the connection")
				return nil
			}
			return &pr.Error{
				Message: err.Error(),
				Code:    1,
			}
		}
		var data_convert pr.ToClient
		err = json.Unmarshal(data_read, &data_convert)

		if err != nil {
			return &pr.Error{
				Message: err.Error(),
				Code:    1,
			}
		}

		if data_convert.Quit {
			return nil
		}

		var dataMsg pr.CmdMessageDatas
		msgError := json.Unmarshal(data_convert.Datas, &dataMsg)
		if msgError != nil {
			return &pr.Error{
				Message: msgError.Error(),
				Code:    1,
			}
		}

		print_err := PrintData(dataMsg, cmd)
		if print_err != nil {
			return print_err
		}
	}
}

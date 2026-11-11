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
	"time"
	pr "treenity/src/protocol"
)

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

func (c *Client) ExecSubscribe(cmd_struct any) *pr.Error {
	cmd := cmd_struct.(pr.CmdSubscribeDatas)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	lineChan := make(chan []byte)
	errChan := make(chan error)

	c.ReadData(lineChan, errChan)

	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			// Send disconnect command to server
			cmd_disconnect := pr.CmdDisconnectDatas{SubscriberName: c.Name}
			return c.Send(cmd_disconnect, pr.CmdDisconnect)

		case <-ticker.C:
			if !ProcessExists(c.ServerPid) {
				return &pr.Error{
					Message: "server down",
					Code:    3,
				}
			}

		case err := <-errChan:
			if err != nil {
				if ctx.Err() != nil {
					return nil
				}
				if errors.Is(err, io.EOF) {
					return &pr.Error{
						Message: "server down",
						Code:    3,
					}
				}
				return &pr.Error{
					Message: err.Error(),
					Code:    1,
				}
			}
			return nil

		case data := <-lineChan:
			var data_convert pr.ToClient
			err := json.Unmarshal(data, &data_convert)

			if err != nil {
				return &pr.Error{
					Message: err.Error(),
					Code:    1,
				}
			}

			// Handle server shutdown
			if data_convert.Quit {
				return nil
			}

			// Cast message
			var dataMsg pr.CmdMessageDatas
			msgError := json.Unmarshal(data_convert.Datas, &dataMsg)
			if msgError != nil {
				return &pr.Error{
					Message: msgError.Error(),
					Code:    1,
				}
			}

			// Print message
			print_err := PrintData(dataMsg, cmd)
			if print_err != nil {
				return print_err
			}

			// Define new offset
			ack_struct := pr.CmdAckDatas{
				Offset:         dataMsg.Offset + 1,
				SubscriberName: c.Name,
			}

			// Send ACK
			ack_err := c.Send(ack_struct, pr.CmdAck)
			if ack_err != nil {
				return ack_err
			}
		}
	}
}

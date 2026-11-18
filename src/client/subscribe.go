package client

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os/signal"
	"syscall"
	"time"
	pr "treenity/src/protocol"
)

func (c *Client) ExecSubscribe(cmd_struct any) *pr.Error {
	cmd := cmd_struct.(pr.CmdSubscribeDatas)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	lineChan := make(chan []byte)
	errChan := make(chan error)

	c.ReadData(lineChan, errChan)

	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	// Handle errors, disconnections, SIGKILL...
	for {
		select {
		case <-ctx.Done():
			done := make(chan struct{})
			go func() {
				c.Send(pr.CmdDisconnectDatas{SubscriberName: c.Name}, pr.CmdDisconnect)
				close(done)
			}()
			select {
			case <-done:
			case <-time.After(pr.ShutdownTimeout):
			}
			return nil

		case <-ticker.C:
			if !ProcessExists(c.ServerPid) {
				select {
				case data := <-lineChan:
					var m pr.ToClient
					if json.Unmarshal(data, &m) == nil && m.Quit {
						return nil
					}
				default:
				}
				return &pr.Error{Message: "server down", Code: 3}
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
			print_err, should_quit := c.PrintNewMessage(cmd, data)
			if print_err != nil {
				return print_err
			} else if should_quit {
				return nil
			}
		}
	}
}

func (c *Client) PrintNewMessage(cmd pr.CmdSubscribeDatas, data []byte) (*pr.Error, bool) {
	var data_convert pr.ToClient
	err := json.Unmarshal(data, &data_convert)

	if err != nil {
		return &pr.Error{
			Message: err.Error(),
			Code:    1,
		}, false
	}

	// Handle server shutdown
	if data_convert.Quit {
		c.Send(pr.CmdDisconnectDatas{SubscriberName: c.Name}, pr.CmdDisconnect)
		return nil, true
	}

	// Cast message
	var dataMsg pr.CmdMessageDatas
	msgError := json.Unmarshal(data_convert.Datas, &dataMsg)
	if msgError != nil {
		return &pr.Error{
			Message: msgError.Error(),
			Code:    1,
		}, false
	}

	// Print message
	print_err := PrintData(dataMsg, cmd)
	if print_err != nil {
		return print_err, false
	}

	// Define new offset
	ack_struct := pr.CmdAckDatas{
		Offset:         dataMsg.Offset + 1,
		SubscriberName: c.Name,
	}

	// Send ACK
	ack_err := c.Send(ack_struct, pr.CmdAck)
	if ack_err != nil {
		return ack_err, false
	}

	return nil, false
}

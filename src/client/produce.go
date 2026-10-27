package client

import (
	"bufio"
	"os"
	"strings"
	pr "treenity/src/protocol"
)

func (c *Client) ExecProduce() *pr.Error {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		// Get message datas
		line := scanner.Text()
		key, body, err := c.GetMessageDatas(line)
		if err != nil {
			return err
		}

		// Send message
		msg_err := c.SendMessage(key, body)
		if msg_err != nil {
			return msg_err
		}
	}

	// Check scanner error
	if scan_err := scanner.Err(); scan_err != nil {
		return &pr.Error{
			Message: scan_err.Error(),
			Code:    1,
		}
	}

	return nil
}

func (c *Client) GetMessageDatas(line string) (string, string, *pr.Error) {
	// Handle Raw flag
	if pr.CLI.IpcIdentifier.Produce.Raw {
		// TODO: Convert Little Endian in string
		return "", "", nil

	} else {
		// Check valid format
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return "", "", &pr.Error{
				Message: pr.ErrInvalidFormat,
				Code:    1,
			}
		}

		return parts[0], parts[1], nil
	}
}

func (c *Client) SendMessage(key, body string) *pr.Error {
	// Send message
	msg_struct := pr.CmdMessageDatas{
		Key:  key,
		Body: body,
	}

	err := c.SendRequest(msg_struct, pr.CmdMessage)
	if err != nil {
		return err
	}

	return nil
}

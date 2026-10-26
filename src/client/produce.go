package client

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	pr "treenity/src/protocol"
)

func (c *Client) ExecProduce() *pr.Error {
	if pr.CLI.IpcIdentifier.Produce.Raw {
		// Raw Mode

	} else {
		// Text Mode

	}

	return nil
}

func (c *Client) TextMode() *pr.Error {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		// Split line
		line := scanner.Text()
		parts := strings.SplitN(line, ":", 2)

		// Check valid format
		if len(parts) != 2 {
			return &pr.Error{
				Message: pr.ErrInvalidFormat,
				Code:    1,
			}
		}
		key, body := parts[0], parts[1]

		// Send message
		err := c.SendMessage(pr.CLI.IpcIdentifier.Produce.TopicName, key, body)
		if err != nil {
			return err
		}
	}

	// Check scanner error
	if err := scanner.Err(); err != nil {
		return &pr.Error{
			Message: err.Error(),
			Code:    1,
		}
	}

	return nil
}

func (c *Client) SendMessage(topic_name, key, body string) *pr.Error {
	fmt.Println(key, body)
	return nil
}

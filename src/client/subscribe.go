package client

import (
	"bufio"
	"fmt"
	"os"
	pr "treenity/src/protocol"
)

func (c *Client) ExecSubscribe(cmd_struct any) *pr.Error {
	fifoPath := c.Path

	file, err := os.OpenFile(fifoPath, os.O_RDONLY, 0)
	if err != nil {
		return &pr.Error{
			Message: fmt.Sprintf("failed to open fifo %s: %s", fifoPath, err.Error()),
			Code:    3,
		}
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
	}

	if err := scanner.Err(); err != nil {
		return &pr.Error{
			Message: fmt.Sprintf("error reading from fifo %s: %s", fifoPath, err.Error()),
			Code:    3,
		}
	}

	return nil
}

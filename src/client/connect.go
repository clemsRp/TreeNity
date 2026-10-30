package client

import (
	"fmt"
	"os"
	pr "treenity/src/protocol"
	"treenity/src/protocol/ipc"
)

func (c *Client) Connect(path string) error {
	writer, err := ipc.NewWriter(path)
	if err != nil {
		return err
	}

	c.Writer = writer

	pid := os.Getpid()
	c.Path = fmt.Sprintf("%s%d", pr.EntryPoint, pid)

	reader, err := ipc.NewReader(c.Path)
	if err != nil {
		return err
	}

	c.Reader = reader

	return nil
}

func (c *Client) Disconnect() {
	c.Writer.Close()
	c.Reader.Close()
}

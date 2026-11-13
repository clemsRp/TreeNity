package client

import (
	"fmt"
	"os"
	"treenity/src/protocol/ipc"
)

func (c *Client) Connect(path string) error {
	// Create Writer
	writer, err := ipc.NewWriter(path)
	if err != nil {
		return err
	}

	c.Writer = writer

	// Get client PID
	pid := os.Getpid()
	c.Path = fmt.Sprintf("%s.%d", path, pid)

	// Create Reader
	reader, err := ipc.NewReader(c.Path)
	if err != nil {
		return err
	}

	c.Reader = reader

	return nil
}

func (c *Client) Disconnect() {
	// Close and Remove Writer and Reader
	c.Writer.Close()
	c.Reader.Close()
}

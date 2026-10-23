package client

import "treenity/src/protocol/ipc"

func (c *Client) Connect(path string) error {
	// Create Writer
	writer, err := ipc.NewWriter(path)
	if err != nil {
		return err
	}

	c.Writer = writer

	return nil
}

func (c *Client) Disconnect() {
	// Close and Remove Writer
	c.Writer.Close()
}

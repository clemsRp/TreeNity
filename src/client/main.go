package client

import (
	"treenity/src/protocol/ipc"
)

type Client struct {
	Name   string
	Path   string
	Writer *ipc.Writer
	Reader *ipc.Reader
}

func NewClient() *Client {
	return &Client{
		Name: "Clément",
	}
}

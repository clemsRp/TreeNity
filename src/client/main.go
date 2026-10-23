package client

import (
	"treenity/src/protocol/ipc"
)

type Client struct {
	Name   string
	Writer *ipc.Writer
}

func NewClient() *Client {
	return &Client{
		Name: "Clément",
	}
}

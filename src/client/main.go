package client

import (
	"treenity/src/protocol/ipc"
)

type Client struct {
	Path   string
	Writer *ipc.Writer
	Reader *ipc.Reader
	ServerPid int64
}

func NewClient() *Client {
	return &Client{}
}

package client

import (
	"sync"
	"treenity/src/protocol/ipc"
)

type Client struct {
	Path           string
	Writer         *ipc.Writer
	Reader         *ipc.Reader
	ServerPid      int64
	Name           string
	disconnectOnce sync.Once
}

func NewClient() *Client {
	return &Client{}
}

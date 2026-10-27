package protocol

import "treenity/src/protocol/ipc"

type Message struct {
	Key    string `json:"key"`
	Body   string `json:"body"`
	Offset uint32 `json:"offset"`
}

type Topic struct {
	Name       string     `json:"name"`
	Messages   []*Message `json:"messages"`
	NewMsgChan chan Message
}

type Client struct {
	Id             string      `json:"client"`
	Topic          string      `json:"topic"`
	Offset         uint32      `json:"offset"`
	Prefix         string      `json:"prefix"`
	PreviousPrefix string      `json:"-"`
	IpcPath        string      `json:"ipc"`
	Writer         *ipc.Writer `json:"-"`
	IsActive       bool        `json:"-"`
}

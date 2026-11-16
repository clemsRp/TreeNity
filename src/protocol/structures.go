package protocol

import (
	"encoding/json"
	"sync/atomic"
	"treenity/src/protocol/ipc"
)

type Message struct {
	Key    []byte `json:"key"`
	Body   []byte `json:"body"`
	Offset uint32 `json:"offset"`
}

type Topic struct {
	Name       string     `json:"name"`
	Messages   []*Message `json:"messages"`
	NewMsgChan chan Message
}

type Client struct {
	Id             string        `json:"client"`
	Topic          string        `json:"topic"`
	Offset         atomic.Uint32 `json:"offset"`
	Prefix         string        `json:"prefix"`
	PreviousPrefix string        `json:"-"`
	IpcPath        string        `json:"ipc"`
	IsActive       atomic.Bool   `json:"-"`
	IsWorking      bool          `json:"-"`
	Writer         *ipc.Writer   `json:"-"`
	MsgChan        chan *Message `json:"-"`
	Done           chan struct{} `json:"-"`
}

func (c *Client) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Id     string `json:"client"`
		Topic  string `json:"topic"`
		Offset uint32 `json:"offset"`
		Prefix string `json:"prefix"`
		Ipc    string `json:"ipc"`
	}{c.Id, c.Topic, c.Offset.Load(), c.Prefix, c.IpcPath})
}

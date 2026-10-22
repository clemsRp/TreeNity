package protocol

type Message struct {
	Key    string `json:"key"`
	Body   string `json:"body"`
	Offset uint32 `json:"offset"`
}

type Topic struct {
	Name     string     `json:"name"`
	Messages []*Message `json:"messages"`
}

type Client struct {
	Id             string `json:"client"`
	Topic          string `json:"topic"`
	Offset         uint32 `json:"offset"`
	Prefix         string `json:"prefix"`
	PreviousPrefix string `json:"-"`
	IpcPath        string `json:"ipc"`
	IsActive       bool   `json:"-"`
}

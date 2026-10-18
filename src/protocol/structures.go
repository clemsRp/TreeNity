package protocol

type Message struct {
	Key    string
	Body   string
	Offset uint32
}

type Topic struct {
	Name     string
	Messages []*Message
}

type Client struct {
	Id      string
	Topic   string
	Offset  uint32
	Prefix  string
	IpcPath string
}

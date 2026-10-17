package protocol

type CmdCreateDatas struct {
	IpcIdentifier string `kong:"-" json:"ipc_identifier"`
	TopicName     string `arg:"" help:"Topic Name" json:"topic_name"`
}

type CmdListDatas struct {
	IpcIdentifier string `kong:"-" json:"ipc_identifier"`
}

type CmdProduceDatas struct {
	IpcIdentifier string `kong:"-" json:"ipc_identifier"`
	TopicName     string `arg:"" help:"Topic Name" json:"topic_name"`
	Raw           bool   `help:"Binary file" json:"raw"`
}

type CmdSubscribeDatas struct {
	IpcIdentifier  string `kong:"-" json:"ipc_identifier"`
	TopicName      string `arg:"" help:"Topic Name" json:"topic_name"`
	SubscriberName string `arg:"" help:"Subscriber Name" json:"subscriber_name"`
	Prefix         string `help:"Message prefix" json:"prefix"`
	Offset         int    `help:"Message starting offset" json:"offset"`
	Raw            bool   `help:"Binary file" json:"raw"`
}

type CmdInfoDatas struct {
	IpcIdentifier  string `kong:"-" json:"ipc_identifier"`
	SubscriberName string `arg:"" help:"Subscriber Name" json:"subscriber_name"`
}

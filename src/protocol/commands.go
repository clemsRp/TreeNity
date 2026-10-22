package protocol

const (
	CmdCreate    = "create"
	CmdList      = "list"
	CmdProduce   = "produce"
	CmdSubscribe = "subscribe"
	CmdInfo      = "info"
)

// Define Client parsing
var CLI struct {
	IpcIdentifier struct {
		IpcIdentifier string            `arg:""`
		Create        CmdCreateDatas    `cmd:"" help:"Create a topic"`
		List          CmdListDatas      `cmd:"" help:"List available topics"`
		Produce       CmdProduceDatas   `cmd:"" help:"Produce messages to a topic"`
		Subscribe     CmdSubscribeDatas `cmd:"" help:"Subscribe to a topic"`
		Info          CmdInfoDatas      `cmd:"" help:"Get infos about a client"`
	} `arg:""`
}

// Define command structures
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
	Offset         *int   `help:"Message starting offset" json:"offset"`
	Raw            bool   `help:"Binary file" json:"raw"`
}

type CmdInfoDatas struct {
	IpcIdentifier  string `kong:"-" json:"ipc_identifier"`
	SubscriberName string `arg:"" help:"Subscriber Name" json:"subscriber_name"`
}

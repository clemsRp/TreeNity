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

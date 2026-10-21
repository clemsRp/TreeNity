package cmd

import (
	"strings"
	pr "treenity/src/protocol"
	"treenity/src/server/topics"
)

func HandleCmdList(command pr.CmdListDatas, topics *topics.Topics) (string, error) {
	// Get Topic names
	topic_names := topics.GetTopicNames()

	return strings.Join(topic_names, ","), nil
}

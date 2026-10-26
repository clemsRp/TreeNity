package cmd

import (
	"strings"
	pr "treenity/src/protocol"
	t "treenity/src/server/topics"
)

func HandleCmdList(command pr.CmdListDatas, topics *t.Topics) (string, any, *pr.Error) {
	// Get Topic names
	topic_names := topics.GetTopicNames()

	return strings.Join(topic_names, ","), "", nil
}

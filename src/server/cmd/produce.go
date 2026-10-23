package cmd

import (
	"fmt"
	pr "treenity/src/protocol"
	t "treenity/src/server/topics"
)

func HandleCmdProduce(command pr.CmdProduceDatas, topics *t.Topics) (string, any, error) {
	// Check topic exist
	_, exist := topics.Get(command.TopicName)
	if !exist {
		return "", "", fmt.Errorf("'%s' %s", command.TopicName, pr.ErrTopicDoesNotExist)
	}

	return "", "", nil
}

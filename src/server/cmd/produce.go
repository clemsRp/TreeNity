package cmd

import (
	"fmt"
	pr "treenity/src/protocol"
	t "treenity/src/server/topics"
)

func HandleCmdProduce(command pr.CmdProduceDatas, topics *t.Topics) (string, any, *pr.Error) {
	// Check topic exist
	_, exist := topics.Get(command.TopicName)
	if !exist {
		return "", "", &pr.Error{
			Message: fmt.Sprintf("'%s' %s", command.TopicName, pr.ErrTopicDoesNotExist),
			Code:    2,
		}
	}

	return "", "", nil
}

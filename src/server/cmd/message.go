package cmd

import (
	"fmt"
	pr "treenity/src/protocol"
	t "treenity/src/server/topics"
)

func HandleCmdMessage(command pr.CmdMessageDatas, topics *t.Topics) (string, any, *pr.Error) {
	topic, exist := topics.Get(command.TopicName)
	if !exist {
		return "", "", &pr.Error{
			Message: fmt.Sprintf("'%s' %s", command.TopicName, pr.ErrTopicDoesNotExist),
			Code:    2,
		}
	}

	topic.Topic.NewMsgChan <- pr.Message{
		Key:     command.Key,
		Body:    command.Body,
		RawKey:  command.RawKey,
		RawBody: command.RawBody,
	}

	return "", "", nil
}

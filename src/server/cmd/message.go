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

	// Get message length
	text_len := len(command.Key) + len(command.Body)

	// Check message length
	if text_len > pr.MaxMessageLength {
		return "", "", &pr.Error{
			Message: pr.ErrMessageTooLong,
			Code:    1,
		}
	}

	// Add new message
	topic.Topic.NewMsgChan <- pr.Message{
		Key:  command.Key,
		Body: command.Body,
	}

	return "", "", nil
}

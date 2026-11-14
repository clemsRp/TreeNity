package cmd

import (
	"fmt"
	pr "treenity/src/protocol"
	t "treenity/src/server/topics"
)

func HandleCmdMessage(command pr.CmdMessageDatas, topics *t.Topics) (string, any, *pr.Error) {
	// Get message length
	text_len := len(command.Key) + len(command.Body)

	// Check message length
	if text_len > pr.MaxMessageLength {
		return "", "", &pr.Error{
			Message: pr.ErrMessageTooLong,
			Code:    1,
		}
	}

	// Check topic is a valid name
	if !IsValidId(command.TopicName) {
		return "", "", &pr.Error{
			Message: fmt.Sprintf("'%s' %s", command.TopicName, pr.ErrInvalidName),
			Code:    1,
		}
	}

	// Check topic exist
	topic, exist := topics.Get(command.TopicName)
	if !exist {
		return "", "", &pr.Error{
			Message: fmt.Sprintf("'%s' %s", command.TopicName, pr.ErrTopicDoesNotExist),
			Code:    2,
		}
	}

	// Add new message
	select {
	case topic.Topic.NewMsgChan <- pr.Message{Key: command.Key, Body: command.Body}:
	case <-topic.Done:
		return "", "", &pr.Error{Message: "server shutting down", Code: 3}
	}

	return "", "", nil
}

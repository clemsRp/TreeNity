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

	// Get datas
	Key := command.Key
	Body := command.Body
	RawKey := command.RawKey
	RawBody := command.RawBody

	// Cast message
	if command.Raw {
		Key = string(RawKey)
		Body = string(RawBody)
	} else {
		RawKey = []byte(Key)
		RawBody = []byte(Body)
	}

	// Get message length
	text_len := len(Key) + len(Body)
	raw_len := len(RawKey) + len(RawBody)
	length := max(text_len, raw_len)

	// Check message length
	if length > pr.MaxMessageLength {
		return "", "", &pr.Error{
			Message: pr.ErrMessageTooLong,
			Code:    1,
		}
	}

	// Add new message
	topic.Topic.NewMsgChan <- pr.Message{
		Key:     Key,
		Body:    Body,
		RawKey:  RawKey,
		RawBody: RawBody,
	}

	return "", "", nil
}

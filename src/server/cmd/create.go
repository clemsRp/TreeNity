package cmd

import (
	"fmt"
	pr "treenity/src/protocol"
	t "treenity/src/server/topics"
)

func HandleCmdCreate(command pr.CmdCreateDatas, topics *t.Topics) (string, any, *pr.Error) {
	// Handle invalid topic name
	if _, ok := topics.Get(command.TopicName); ok {
		return "", "", &pr.Error{
			Message: fmt.Sprintf("'%s' %s", command.TopicName, pr.ErrTopicAlreadyExist),
			Code:    2,
		}

	} else if !IsValidId(command.TopicName) {
		return "", "", &pr.Error{
			Message: fmt.Sprintf("'%s' %s", command.TopicName, pr.ErrInvalidName),
			Code:    1,
		}
	}

	// Create Topic
	topics.Create(command.TopicName)

	return "topic created", "", nil
}

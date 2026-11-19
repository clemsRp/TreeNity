package cmd

import (
	"fmt"
	pr "treenity/src/protocol"
	t "treenity/src/server/topics"
)

func HandleCmdCreate(command pr.CmdCreateDatas, topics *t.Topics) (string, any, *pr.Error) {
	// Handle invalid topic name
	if !IsValidId(command.TopicName) {
		return "", "", &pr.Error{
			Message: fmt.Sprintf("'%s' %s", command.TopicName, pr.ErrInvalidName),
			Code:    1,
		}

	} else if _, ok := topics.Get(command.TopicName); ok {
		return "", "", &pr.Error{
			Message: fmt.Sprintf("'%s' %s", command.TopicName, pr.ErrTopicAlreadyExist),
			Code:    2,
		}
	}

	// Create Topic
	new_topic := topics.Create(command.TopicName)

	go new_topic.DistributeMessages()

	return "topic created", "", nil
}

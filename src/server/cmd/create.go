package cmd

import (
	"fmt"
	pr "treenity/src/protocol"
	t "treenity/src/server/topics"
)

func HandleCmdCreate(command pr.CmdCreateDatas, topics *t.Topics) (string, any, error) {
	// Handle invalid topic name
	if _, ok := topics.Get(command.TopicName); ok {
		return "", "", fmt.Errorf("'%s' %s", command.TopicName, pr.ErrTopicAlreadyExist)

	} else if !IsValidId(command.TopicName) {
		return "", "", fmt.Errorf("'%s' %s", command.TopicName, pr.ErrInvalidName)
	}

	// Create Topic
	topics.Create(command.TopicName)

	return "topic created", "", nil
}

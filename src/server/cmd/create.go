package cmd

import (
	"fmt"
	pr "treenity/src/protocol"
	"treenity/src/server/topics"
)

func HandleCmdCreate(command pr.CmdCreateDatas, topics *topics.Topics) (string, error) {
	// Handle invalid topic name
	if _, ok := topics.Get(command.TopicName); ok {
		return "", fmt.Errorf("'%s' topic already exist", command.TopicName)

	} else if !IsValidId(command.TopicName) {
		return "", fmt.Errorf("invalid topic name '%s', must respect this format `^[a-zA-Z0-9_.-]+$`", command.TopicName)
	}

	// Create Topic
	topics.Create(command.TopicName)

	return "topic created", nil
}

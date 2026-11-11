package cmd

import (
	"fmt"
	pr "treenity/src/protocol"
	h "treenity/src/server/hashmap"
	t "treenity/src/server/topics"
)

func HandleCmdCreate(command pr.CmdCreateDatas, topics *t.Topics, clients *h.Hashmap) (string, any, *pr.Error) {
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

	go new_topic.DistributeMessages(clients)

	return "topic created", "", nil
}

// Backward compatibility
func DistributeMessages(topic *t.TopicDatas, clients *h.Hashmap) {
	topic.DistributeMessages(clients)
}

// Backward compatibility
func AcceptNewMessages(topic *t.TopicDatas, clients *h.Hashmap) {
	topic.DistributeMessages(clients)
}

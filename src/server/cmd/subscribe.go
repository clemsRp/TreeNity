package cmd

import (
	"fmt"
	pr "treenity/src/protocol"
	h "treenity/src/server/hashmap"
	t "treenity/src/server/topics"
)

func HandleCmdSubscribe(command pr.CmdSubscribeDatas, topics *t.Topics, clients *h.Hashmap, ipcPath string) (string, any, *pr.Error) {
	var client *pr.Client
	var err error

	// Check topic exist
	topic, exist := topics.Get(command.TopicName)
	if !exist {
		return "", "", &pr.Error{
			Message: fmt.Sprintf("'%s' %s", command.TopicName, pr.ErrTopicDoesNotExist),
			Code:    2,
		}
	}

	client, err = clients.Get(command.SubscriberName)

	// Check subscriber is valid name
	if !IsValidId(command.SubscriberName) {
		return "", "", &pr.Error{
			Message: fmt.Sprintf("'%s' %s", command.SubscriberName, pr.ErrInvalidName),
			Code:    1,
		}
	}

	// Check subscriber doesn't exist or is reconnecting
	if err == nil {
		// Check client is connected
		if client.IsActive {
			return "", "", &pr.Error{
				Message: fmt.Sprintf("'%s' %s", client.Id, pr.ErrSubscriberAlreadyActive),
				Code:    2,
			}
		}

		// Update client datas
		client.IsActive = true

		// Remove client from previous topic tree if needed
		if client.Topic != command.TopicName {
			if old_topic, exist := topics.Get(client.Topic); exist {
				old_topic.Tree.Remove(client)
			}
			client.Topic = command.TopicName
		}

		// Update prefix if needed
		if command.Prefix != nil {
			client.PreviousPrefix = client.Prefix
			client.Prefix = *command.Prefix
		}

		// Update offset if needed
		if command.Offset != nil {
			client.Offset = uint32(*command.Offset)
		}

	} else {
		// Handle missing offset option
		offset := uint32(0)
		if command.Offset != nil {
			offset = uint32(*command.Offset)
		}

		// Handle missing prefix option
		prefix := ""
		if command.Prefix != nil {
			prefix = *command.Prefix
		}

		// Create subscriber
		client = &pr.Client{
			Id:       command.SubscriberName,
			Topic:    command.TopicName,
			Offset:   offset,
			Prefix:   prefix,
			IpcPath:  ipcPath,
			IsActive: true,
		}

		// Add client to Hashmap
		clients.Append(command.SubscriberName, client)
	}

	// Add client to topic prefix tree
	topic.Tree.Insert(client)

	return "subscribed to " + command.TopicName, "", nil
}

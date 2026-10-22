package cmd

import (
	"fmt"
	pr "treenity/src/protocol"
	h "treenity/src/server/hashmap"
	t "treenity/src/server/topics"
)

func HandleCmdSubscribe(command pr.CmdSubscribeDatas, topics *t.Topics, clients *h.Hashmap) (string, any, error) {
	var client *pr.Client
	var err error

	// Check topic exist
	topic, exist := topics.Get(command.TopicName)
	if !exist {
		return "", "", fmt.Errorf("'%s' topic does not exist", command.TopicName)
	}

	client, err = clients.Get(command.SubscriberName)

	// Check subscriber is valid name
	if !IsValidId(command.SubscriberName) {
		return "", "", fmt.Errorf("invalid subscriber name '%s', must respect this format `^[a-zA-Z0-9_.-]+$`", command.SubscriberName)
	}

	// Check subscriber doesn't exist or is reconnecting
	if err == nil {
		// Check client is connected
		if client.IsActive {
			return "", "", fmt.Errorf("'%s' subscriber is already active", client.Id)
		}

		// Update prefix
		client.PreviousPrefix = client.Prefix
		client.Prefix = command.Prefix
		client.IsActive = true
		client.Topic = command.TopicName

	} else {
		// Create subscriber
		client = &pr.Client{
			Id:       command.SubscriberName,
			Topic:    command.TopicName,
			Offset:   uint32(*command.Offset),
			Prefix:   command.Prefix,
			IpcPath:  "TODO",
			IsActive: true,
		}

		// Add client to Hashmap
		clients.Append(command.SubscriberName, client)
	}

	// Add client to topic prefix tree
	topic.Tree.Insert(client)

	return "subscribed to " + command.TopicName, "", nil
}

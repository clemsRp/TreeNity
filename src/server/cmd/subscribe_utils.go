package cmd

import (
	"fmt"
	pr "treenity/src/protocol"
	"treenity/src/protocol/ipc"
	h "treenity/src/server/hashmap"
	t "treenity/src/server/topics"
)

func CheckSubscribeCmd(command pr.CmdSubscribeDatas, topics *t.Topics, ipcPath string) *pr.Error {
	// Check response channel
	if ipcPath == "" {
		return &pr.Error{
			Message: "missing client IPC path",
			Code:    1,
		}
	}

	// Check subscriber is valid name
	if !IsValidId(command.SubscriberName) {
		return &pr.Error{
			Message: fmt.Sprintf("'%s' %s", command.SubscriberName, pr.ErrInvalidName),
			Code:    1,
		}
	}
	// Check topic is valid name
	if !IsValidId(command.TopicName) {
		return &pr.Error{
			Message: fmt.Sprintf("'%s' %s", command.TopicName, pr.ErrInvalidName),
			Code:    1,
		}
	}

	// Check topic exist
	_, exist := topics.Get(command.TopicName)
	if !exist {
		return &pr.Error{
			Message: fmt.Sprintf("'%s' %s", command.TopicName, pr.ErrTopicDoesNotExist),
			Code:    2,
		}
	}

	return nil
}

func AddNewSubscriber(command pr.CmdSubscribeDatas, clients *h.Hashmap, ipcPath string) *pr.Error {
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

	// Init client Writer
	writer, wri_err := ipc.NewWriter(ipcPath)
	if wri_err != nil {
		return &pr.Error{
			Message: wri_err.Error(),
			Code:    3,
		}
	}

	// Create subscriber
	client := &pr.Client{
		Id:             command.SubscriberName,
		Topic:          command.TopicName,
		Prefix:         prefix,
		PreviousPrefix: prefix,
		IpcPath:        ipcPath,
		Writer:         writer,
	}
	client.Offset.Store(offset)
	client.IsActive.Store(true)

	// Add client to Hashmap
	clients.Append(command.SubscriberName, client)

	return nil
}

func ReconnectSubscriber(command pr.CmdSubscribeDatas, topics *t.Topics, client *pr.Client, ipcPath string) {
	// Update client datas
	client.IsActive.Store(true)
	client.IpcPath = ipcPath

	// Remove client from previous topic tree if needed
	if client.Topic != command.TopicName {
		if old_topic, exist := topics.Get(client.Topic); exist {
			old_topic.Tree.Remove(client)
		}
		client.Topic = command.TopicName
		client.Offset.Store(0)
	}

	// Update prefix if needed
	client.PreviousPrefix = client.Prefix
	if command.Prefix != nil {
		client.Prefix = *command.Prefix
	}

	// Update offset if needed
	if command.Offset != nil {
		client.Offset.Store(uint32(*command.Offset))
	}
}

func StartSubscription(command pr.CmdSubscribeDatas, topics *t.Topics, clients *h.Hashmap) {
	topic, exist := topics.Get(command.TopicName)
	if !exist {
		return
	}

	client, err := clients.Get(command.SubscriberName)
	if err != nil || client == nil {
		return
	}

	topic.StartSubscriber(client, clients)
}

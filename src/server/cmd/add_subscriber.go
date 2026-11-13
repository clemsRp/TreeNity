package cmd

import (
	pr "treenity/src/protocol"
	"treenity/src/protocol/ipc"
	h "treenity/src/server/hashmap"
	t "treenity/src/server/topics"
)

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
		Offset:         offset,
		Prefix:         prefix,
		PreviousPrefix: prefix,
		IpcPath:        ipcPath,
		Writer:         writer,
		IsActive:       true,
	}

	// Add client to Hashmap
	clients.Append(command.SubscriberName, client)

	return nil
}

func ReconnectSubscriber(command pr.CmdSubscribeDatas, topics *t.Topics, client *pr.Client, ipcPath string) {
	// Update client datas
	client.IsActive = true
	client.IpcPath = ipcPath

	// Remove client from previous topic tree if needed
	if client.Topic != command.TopicName {
		if old_topic, exist := topics.Get(client.Topic); exist {
			old_topic.Tree.Remove(client)
		}
		client.Topic = command.TopicName
		client.Offset = 0
	}

	// Update prefix if needed
	client.PreviousPrefix = client.Prefix
	if command.Prefix != nil {
		client.Prefix = *command.Prefix
	}

	// Update offset if needed
	if command.Offset != nil {
		client.Offset = uint32(*command.Offset)
	}
}

package cmd

import (
	"fmt"
	pr "treenity/src/protocol"
	"treenity/src/protocol/ipc"
	h "treenity/src/server/hashmap"
	t "treenity/src/server/topics"
)

func HandleCmdCreate(command pr.CmdCreateDatas, topics *t.Topics, clients *h.Hashmap) (string, any, *pr.Error) {
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
	new_topic := topics.Create(command.TopicName)

	go AcceptNewMessages(new_topic, clients)

	return "topic created", "", nil
}

func AcceptNewMessages(topic *t.TopicDatas, clients *h.Hashmap) {
	offset := uint32(0)
	for {
		// Add message
		msg := <-topic.Topic.NewMsgChan
		msg.Offset = offset
		topic.Topic.Messages = append(topic.Topic.Messages, &msg)
		offset++

		// Broadcast to subscriber
		subscribers := topic.Tree.Filter(msg.Key)
		for _, client := range subscribers {
			// Skip inactive clients
			if !client.IsActive {
				continue
			}

			// Init writer if needed
			if client.Writer == nil {
				writer, err := ipc.NewWriter(client.IpcPath)
				if err != nil {
					disconnectSubscriber(client.Id, client, clients)
					continue
				}

				client.Writer = writer
			}

			// Write message
			if err := SendResponse(client.Writer, pr.ToServer{}, "", msg, nil); err != nil {
				disconnectSubscriber(client.Id, client, clients)
			}
		}
	}
}

func disconnectSubscriber(clientId string, client *pr.Client, clients *h.Hashmap) {
	target := client
	if clients != nil {
		if c, err := clients.Get(clientId); err == nil {
			target = c
		}
	}

	if target.Writer != nil {
		_ = target.Writer.Close()
		target.Writer = nil
	}
	target.IsActive = false
}

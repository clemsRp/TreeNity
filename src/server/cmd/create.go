package cmd

import (
	"fmt"
	pr "treenity/src/protocol"
	"treenity/src/protocol/ipc"
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
	new_topic := topics.Create(command.TopicName)

	go AcceptNewMessages(new_topic)

	return "topic created", "", nil
}

func AcceptNewMessages(topic *t.TopicDatas) {
	offset := uint32(0)
	for {
		// Add message
		msg := <-topic.Topic.NewMsgChan
		msg.Offset = offset
		topic.Topic.Messages = append(topic.Topic.Messages, &msg)
		offset++

		// Broadcast to subscriber
		clients := topic.Tree.Filter(msg.Key)
		for _, client := range clients {
			// Init writer if needed
			if client.Writer == nil {
				writer, err := ipc.NewWriter(client.IpcPath)
				if err != nil {
					continue
				}

				client.Writer = writer
			}

			// Write message
			if err := SendResponse(client.Writer, pr.ToServer{}, "", msg, nil); err != nil {
				client.Writer.Close()
				client.Writer = nil
				client.IsActive = false
			}
		}
	}
}

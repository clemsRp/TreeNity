package cmd

import (
	"encoding/json"
	"fmt"
	"strings"
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
		// 1. Detect and pop new messages added to the topic
		msg := <-topic.Topic.NewMsgChan
		msg.Offset = offset
		topic.Topic.Messages = append(topic.Topic.Messages, &msg)
		offset++

		// 2. Thread-safe subscriber retrieval (managed by topic.Tree with sync.RWMutex)
		subscribers := topic.Tree.Filter(msg.Key)

		// Format message as CmdMessageDatas for subscribers
		cmdData := pr.CmdMessageDatas{
			Key:       msg.Key,
			Body:      msg.Body,
			TopicName: topic.Topic.Name,
		}
		payload, err := json.Marshal(cmdData)
		if err != nil {
			continue
		}

		// 3. For each subscriber: verify IsActive, offset match, and key prefix
		for _, client := range subscribers {
			if !client.IsActive {
				continue
			}

			// Check subscriber offset
			if msg.Offset < client.Offset {
				continue
			}

			// Validate message key prefix
			if client.Prefix != "" && !strings.HasPrefix(msg.Key, client.Prefix) {
				continue
			}

			// 4. Open client FIFO (IpcPath) in write-only mode (os.O_WRONLY) if not already open
			if client.Writer == nil {
				writer, err := ipc.NewWriter(client.IpcPath)
				if err != nil {
					disconnectSubscriber(client.Id, client, clients)
					continue
				}
				client.Writer = writer
			}

			// Write formatted message
			if err := client.Writer.WriteLine(payload); err != nil {
				// 5. Handle write errors (syscall.EPIPE / broken pipe) by deactivating client and closing descriptor
				disconnectSubscriber(client.Id, client, clients)
				continue
			}

			// Update client offset after successful delivery
			client.Offset = msg.Offset + 1
		}
	}
}

func disconnectSubscriber(clientId string, client *pr.Client, clients *h.Hashmap) {
	if client != nil {
		if client.Writer != nil {
			_ = client.Writer.Close()
			client.Writer = nil
		}
		client.IsActive = false
	}

	if clients != nil {
		if c, err := clients.Get(clientId); err == nil && c != nil && c != client {
			if c.Writer != nil {
				_ = c.Writer.Close()
				c.Writer = nil
			}
			c.IsActive = false
		}
	}
}

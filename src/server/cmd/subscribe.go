package cmd

import (
	"fmt"
	"strings"
	pr "treenity/src/protocol"
	h "treenity/src/server/hashmap"
	t "treenity/src/server/topics"
)

func HandleCmdSubscribe(command pr.CmdSubscribeDatas, topics *t.Topics, clients *h.Hashmap, ipcPath string) (string, any, *pr.Error) {
	var client *pr.Client
	var err error

	// Check response channel
	if ipcPath == "" {
		return "", "", &pr.Error{
			Message: "missing client IPC path",
			Code:    1,
		}
	}

	// Subscriber IPC channel
	permanentIpcPath := ipcPath

	// Check topic exist
	_, exist := topics.Get(command.TopicName)
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
		client.IpcPath = permanentIpcPath
		if client.Writer != nil {
			client.Writer.Close()
			client.Writer = nil
		}

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
			IpcPath:  permanentIpcPath,
			IsActive: true,
		}

		// Add client to Hashmap
		clients.Append(command.SubscriberName, client)
	}

	return "subscribed to " + command.TopicName, "", nil
}

// Replay messages and insert client
func ReplayAndInsert(command pr.CmdSubscribeDatas, topics *t.Topics, clients *h.Hashmap) {
	topic, exist := topics.Get(command.TopicName)
	if !exist {
		return
	}

	client, err := clients.Get(command.SubscriberName)
	if err != nil || client == nil {
		return
	}

	topic.Mu.Lock()
	defer topic.Mu.Unlock()

	for _, msg := range topic.Topic.Messages {
		if !client.IsActive {
			break
		}
		if client.Offset <= msg.Offset {
			if strings.HasPrefix(msg.Key, client.Prefix) {
				if !topic.SendToSubscriber(msg, client, clients) {
					break
				}
			} else {
				client.Offset = msg.Offset + 1
			}
		}
	}

	// Subscriber worker goroutine
	if client.MsgChan == nil {
		client.MsgChan = make(chan *pr.Message, 100)
	}
	go topic.StartSubscriberWorker(client, clients)

	// Add client to topic prefix tree
	topic.Tree.Insert(client)
}

// Backward compatibility
func ReplayAndInsertSubscriber(command pr.CmdSubscribeDatas, topics *t.Topics, clients *h.Hashmap) {
	ReplayAndInsert(command, topics, clients)
}

package cmd

import (
	"fmt"
	"os"
	"strings"
	pr "treenity/src/protocol"
	"treenity/src/protocol/ipc"
	h "treenity/src/server/hashmap"
	t "treenity/src/server/topics"
)

func HandleCmdSubscribe(command pr.CmdSubscribeDatas, topics *t.Topics, clients *h.Hashmap, ipcPath string) (string, any, *pr.Error) {
	var client *pr.Client
	var err error

	sub_err := CheckSubscribeCmd(command, topics, ipcPath)
	if sub_err != nil {
		return "", "", sub_err
	}

	client, err = clients.Get(command.SubscriberName)

	// Check subscriber doesn't exist or is reconnecting
	if err == nil {
		// Check client is connected
		if client.IsActive {
			if clientAlive(client.IpcPath) {
				return "", "", &pr.Error{
					Message: fmt.Sprintf("'%s' %s", client.Id, pr.ErrSubscriberAlreadyActive),
					Code:    2,
				}
			}
			os.Remove(client.IpcPath)
		}

		client.Writer.Close()

		// Init client Writer
		writer, wri_err := ipc.NewWriter(ipcPath)
		if wri_err != nil {
			return "", "", &pr.Error{
				Message: wri_err.Error(),
				Code:    3,
			}
		}
		client.Writer = writer

		ReconnectSubscriber(command, topics, client, ipcPath)

	} else {
		err := AddNewSubscriber(command, clients, ipcPath)
		if err != nil {
			return "", "", err
		}
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
			if strings.HasPrefix(string(msg.Key), client.Prefix) {
				if !topic.SendToSubscriber(msg, client, clients) {
					break
				}
			}
		}
	}

	// Add client to topic prefix tree
	topic.Tree.Insert(client)

	// Subscriber worker goroutine
	if client.MsgChan == nil {
		client.MsgChan = make(chan *pr.Message, 100)
	}
	go topic.StartSubscriberWorker(client, clients)
}

// Backward compatibility
func ReplayAndInsertSubscriber(command pr.CmdSubscribeDatas, topics *t.Topics, clients *h.Hashmap) {
	ReplayAndInsert(command, topics, clients)
}

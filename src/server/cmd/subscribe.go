package cmd

import (
	"fmt"
	"os"
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
		if client.IsActive.Load() {
			if clientAlive(client.IpcPath) {
				return "", "", &pr.Error{
					Message: fmt.Sprintf("'%s' %s", client.Id, pr.ErrSubscriberAlreadyActive),
					Code:    2,
				}
			}
			os.Remove(client.IpcPath)
		}

		// Stop previous worker
		t.StopSubscriber(client)
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

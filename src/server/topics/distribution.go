package topics

import (
	"encoding/json"
	"os"
	"strings"
	"syscall"

	pr "treenity/src/protocol"
	h "treenity/src/server/hashmap"
)

func (td *TopicDatas) DistributeMessages(clients *h.Hashmap) {
	for msg := range td.Topic.NewMsgChan {

		td.Mu.Lock()
		msg.Offset = uint32(len(td.Topic.Messages))
		msgCopy := msg
		td.Topic.Messages = append(td.Topic.Messages, &msgCopy)

		subscribers := td.Tree.Filter(msg.Key)

		for _, client := range subscribers {
			if !client.IsActive {
				continue
			}

			if client.Offset <= msg.Offset {
				if strings.HasPrefix(msg.Key, client.Prefix) {
					if client.MsgChan == nil {
						client.MsgChan = make(chan *pr.Message, 100)
						go td.StartSubscriberWorker(client, clients)
					}
					select {
					case client.MsgChan <- &msgCopy:
					default:
						go func(c *pr.Client, m *pr.Message) {
							c.MsgChan <- m
						}(client, &msgCopy)
					}
				} else {
					client.Offset = msg.Offset + 1
				}
			}
		}
		td.Mu.Unlock()
	}
}

func (td *TopicDatas) StartSubscriberWorker(client *pr.Client, clients *h.Hashmap) {
	for msg := range client.MsgChan {
		if !client.IsActive {
			break
		}
		if !td.SendToSubscriber(msg, client, clients) {
			break
		}
	}
}

func (td *TopicDatas) SendToSubscriber(msg *pr.Message, client *pr.Client, clients *h.Hashmap) bool {
	var pipe *os.File
	var err error

	pipe, err = os.OpenFile(client.IpcPath, os.O_WRONLY|syscall.O_NONBLOCK, 0)
	if err != nil {
		disconnectSubscriber(client.Id, client, clients)
		return false
	}

	msgData := pr.CmdMessageDatas{
		Key:       msg.Key,
		Body:      msg.Body,
		RawKey:    msg.RawKey,
		RawBody:   msg.RawBody,
		Offset:    msg.Offset,
		TopicName: td.Topic.Name,
	}

	var payload []byte

	payload, err = json.Marshal(msgData)
	if err != nil {
		_ = pipe.Close()
		return false
	}

	SendClient := pr.ToClient{
		Datas: json.RawMessage(payload),
	}

	payload, err = json.Marshal(SendClient)
	if err != nil {
		_ = pipe.Close()
		return false
	}
	payload = append(payload, '\n')

	_, writeErr := pipe.Write(payload)
	_ = pipe.Close()

	if writeErr != nil {
		disconnectSubscriber(client.Id, client, clients)
		return false
	}

	client.Offset = msg.Offset + 1
	return true
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

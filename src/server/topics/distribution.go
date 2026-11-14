package topics

import (
	"encoding/json"

	pr "treenity/src/protocol"
	h "treenity/src/server/hashmap"
)

func (td *TopicDatas) dispatch(msg pr.Message, clients *h.Hashmap) {
	td.Mu.Lock()
	msg.Offset = uint32(len(td.Topic.Messages))
	msgCopy := msg
	td.Topic.Messages = append(td.Topic.Messages, &msgCopy)

	subscribers := td.Tree.Filter(string(msg.Key))

	for _, client := range subscribers {
		if !client.IsActive {
			continue
		}

		if client.Offset <= msg.Offset {
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
		}
	}
	td.Mu.Unlock()
}

func (td *TopicDatas) DistributeMessages(clients *h.Hashmap) {
	for {
		select {
		case <-td.Done:
			return
		case msg := <-td.Topic.NewMsgChan:
			td.dispatch(msg, clients)
		}
	}
}

func (td *TopicDatas) StartSubscriberWorker(client *pr.Client, clients *h.Hashmap) {
	for {
		select {
		case <-td.Done:
			return
		case msg := <-client.MsgChan:
			if !client.IsActive || !td.SendToSubscriber(msg, client, clients) {
				return
			}
		}
	}
}

func (td *TopicDatas) SendToSubscriber(msg *pr.Message, client *pr.Client, clients *h.Hashmap) bool {
	var err error

	msgData := pr.CmdMessageDatas{
		Key:       msg.Key,
		Body:      msg.Body,
		Offset:    msg.Offset,
		TopicName: td.Topic.Name,
	}

	var payload []byte

	payload, err = json.Marshal(msgData)
	if err != nil {
		return false
	}

	SendClient := pr.ToClient{
		Datas: json.RawMessage(payload),
	}

	payload, err = json.Marshal(SendClient)
	if err != nil {
		return false
	}

	writeErr := client.Writer.WriteLine(payload)

	if writeErr != nil {
		disconnectSubscriber(client.Id, client, clients)
		return false
	}

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

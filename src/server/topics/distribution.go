package topics

import (
	"encoding/json"

	pr "treenity/src/protocol"
	h "treenity/src/server/hashmap"
)

func (td *TopicDatas) dispatch(msg pr.Message) {
	type target struct {
		ch   chan *pr.Message
		done chan struct{}
	}

	td.Mu.Lock()
	msg.Offset = uint32(len(td.Topic.Messages))
	msgCopy := msg
	td.Topic.Messages = append(td.Topic.Messages, &msgCopy)

	var targets []target
	subscribers := td.Tree.Filter(string(msg.Key))

	// Get message targets
	for _, client := range subscribers {
		if client.IsActive.Load() && client.Offset.Load() <= msg.Offset && client.MsgChan != nil {
			targets = append(targets, target{client.MsgChan, client.Done})
		}
	}
	td.Mu.Unlock()

	// Conserv message order
	for _, t := range targets {
		select {
		case t.ch <- &msgCopy: // Send message
		case <-t.done: // Wroker dead, continue to next subscriber
		case <-td.Done: // Server shutdown
			return
		}
	}
}

func (td *TopicDatas) DistributeMessages(clients *h.Hashmap) {
	for {
		select {
		case <-td.Done:
			return
		case msg := <-td.Topic.NewMsgChan:
			td.dispatch(msg)
		}
	}
}

func (td *TopicDatas) StartSubscriberWorker(client *pr.Client, clients *h.Hashmap, ch <-chan *pr.Message, done chan struct{}) {
	// Deactive worker
	defer func() {
		close(done)
		td.Mu.Lock()
		client.IsWorking = false
		td.Mu.Unlock()
	}()

	for {
		select {
		case <-td.Done:
			return
		case msg := <-ch:
			if !client.IsActive.Load() || !td.SendToSubscriber(msg, client, clients) {
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
		client.IsActive.Store(false)
	}

	if clients != nil {
		if c, err := clients.Get(clientId); err == nil && c != nil && c != client {
			if c.Writer != nil {
				_ = c.Writer.Close()
				c.Writer = nil
			}
			c.IsActive.Store(false)
		}
	}
}

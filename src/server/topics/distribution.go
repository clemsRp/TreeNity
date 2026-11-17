package topics

import (
	"encoding/json"
	"strings"

	pr "treenity/src/protocol"
	h "treenity/src/server/hashmap"
)

func (td *TopicDatas) dispatch(msg pr.Message) {
	td.Mu.Lock()
	defer td.Mu.Unlock()

	msg.Offset = uint32(len(td.Topic.Messages))
	td.Topic.Messages = append(td.Topic.Messages, &msg)

	for _, client := range td.Tree.Filter(string(msg.Key)) {
		if client.Notify == nil {
			continue
		}
		select {
		case client.Notify <- struct{}{}:
		default:
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

func (td *TopicDatas) StartSubscriber(client *pr.Client, clients *h.Hashmap) {
	td.Mu.Lock()
	notify := make(chan struct{}, 1)
	client.Notify = notify
	client.Done = make(chan struct{})
	client.Stop = make(chan struct{})
	done, stop := client.Done, client.Stop
	_ = td.Tree.Insert(client)
	td.Mu.Unlock()

	go td.subscriberWorker(client, clients, notify, done, stop)
}

func StopSubscriber(client *pr.Client) {
	if client == nil || client.Stop == nil {
		return
	}

	// Stop worker
	close(client.Stop)
	client.Stop = nil

	// Close FIFO
	client.Writer.Close()

	// Wait for true end
	<-client.Done
}

func (td *TopicDatas) subscriberWorker(client *pr.Client, clients *h.Hashmap, notify chan struct{}, done, stop chan struct{}) {
	defer close(done)

	next := client.Offset.Load()

	for {
		var batch []*pr.Message
		td.Mu.Lock()
		if int(next) < len(td.Topic.Messages) {
			// Messages are append-only and never modified: snapshot is safe
			batch = td.Topic.Messages[next:]
		}
		td.Mu.Unlock()

		// Caught up: wait for a new message
		if len(batch) == 0 {
			select {
			case <-td.Done:
				return
			case <-stop:
				return
			case <-notify:
			}
			continue
		}

		for _, msg := range batch {
			select {
			case <-td.Done:
				return
			case <-stop:
				return
			default:
			}
			if !client.IsActive.Load() {
				return
			}

			next = msg.Offset + 1
			if !strings.HasPrefix(string(msg.Key), client.Prefix) {
				continue
			}
			if !td.deliver(msg, client, clients, stop) {
				return
			}
		}
	}
}

// func (td *TopicDatas) replay(client *pr.Client, clients *h.Hashmap, ch chan *pr.Message, stop chan struct{}) (uint32, bool) {
// 	next := client.Offset.Load()

// 	for {
// 		td.Mu.Lock()
// 		if int(next) >= len(td.Topic.Messages) {
// 			td.Tree.Insert(client)
// 			client.MsgChan = ch
// 			td.Mu.Unlock()
// 			return next, true
// 		}

// 		batch := td.Topic.Messages[next:]
// 		td.Mu.Unlock()

// 		for _, msg := range batch {
// 			select {
// 			case <-td.Done:
// 				return next, false
// 			case <-stop:
// 				return next, false
// 			default:
// 			}
// 			if !client.IsActive.Load() {
// 				return next, false
// 			}

// 			next = msg.Offset + 1
// 			if !strings.HasPrefix(string(msg.Key), client.Prefix) {
// 				continue
// 			}
// 			if !td.deliver(msg, client, clients, stop) {
// 				return next, false
// 			}
// 		}
// 	}
// }

func (td *TopicDatas) deliver(msg *pr.Message, client *pr.Client, clients *h.Hashmap, stop chan struct{}) bool {
	err := td.SendToSubscriber(msg, client)
	if err == nil {
		return true
	}

	select {
	case <-stop:
	default:
		disconnectSubscriber(client.Id, client, clients)
	}
	return false
}

func (td *TopicDatas) SendToSubscriber(msg *pr.Message, client *pr.Client) error {
	msgData := pr.CmdMessageDatas{
		Key:       msg.Key,
		Body:      msg.Body,
		Offset:    msg.Offset,
		TopicName: td.Topic.Name,
	}

	payload, err := json.Marshal(msgData)
	if err != nil {
		return err
	}

	SendClient := pr.ToClient{
		Datas: json.RawMessage(payload),
	}

	payload, err = json.Marshal(SendClient)
	if err != nil {
		return err
	}

	return client.Writer.WriteLineTimeout(payload, pr.SubscriberWriteTimeout)
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

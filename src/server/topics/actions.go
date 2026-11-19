package topics

import (
	"fmt"
	"slices"
	pr "treenity/src/protocol"
	p "treenity/src/server/prefix"
)

func (t *Topics) Create(topic_name string) *TopicDatas {
	// Handle Thread safety
	t.mu.Lock()
	defer t.mu.Unlock()

	// Check topic exist
	_, exist := t.AllTopics[topic_name]
	if exist {
		return t.AllTopics[topic_name]
	}

	// Create topic
	t.AllTopics[topic_name] = &TopicDatas{
		Topic: &pr.Topic{
			Name:       topic_name,
			Messages:   make([]*pr.Message, 0),
			NewMsgChan: make(chan pr.Message),
		},
		Tree: p.CreateTree(0),
		Done: make(chan struct{}),
	}

	// Add topic to topic list
	t.TopicsList = append(t.TopicsList, topic_name)

	return t.AllTopics[topic_name]
}

func (t *Topics) Delete(topic_name string) *pr.Error {
	// Handle Thread safety
	t.mu.Lock()
	defer t.mu.Unlock()

	// Check topic exist
	topic, exist := t.AllTopics[topic_name]
	if !exist {
		return &pr.Error{
			Message: fmt.Sprintf("'%s' %s", topic_name, pr.ErrTopicDoesNotExist),
			Code:    2,
		}
	}

	topic.Stop()

	// Delete topic
	delete(t.AllTopics, topic_name)
	t.TopicsList = slices.DeleteFunc(t.TopicsList, func(n string) bool {
		return n == topic_name
	})

	return nil
}

func (t *Topics) Get(topic_name string) (*TopicDatas, bool) {
	// Handle Thread safety
	t.mu.RLock()
	defer t.mu.RUnlock()

	// Get topic
	topic, exist := t.AllTopics[topic_name]

	return topic, exist
}

func (t *Topics) GetTopicNames() []string {
	// Handle Thread safety
	t.mu.RLock()
	defer t.mu.RUnlock()

	// Get topic
	topics := t.TopicsList

	return topics
}

func (t *Topics) StopAll() {
	t.mu.RLock()
	all := make([]*TopicDatas, 0, len(t.AllTopics))
	for _, td := range t.AllTopics {
		all = append(all, td)
	}
	t.mu.RUnlock()

	for _, td := range all {
		td.Stop()
	}
}

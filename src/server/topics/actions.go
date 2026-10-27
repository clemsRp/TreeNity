package topics

import (
	pr "treenity/src/protocol"
	p "treenity/src/server/prefix"
)

func (t *Topics) Create(topic_name string) *TopicDatas {
	// Handle Thread safety
	t.mu.Lock()
	defer t.mu.Unlock()

	// Create topic
	t.AllTopics[topic_name] = &TopicDatas{
		Topic: &pr.Topic{
			Name:       topic_name,
			Messages:   make([]*pr.Message, 0),
			NewMsgChan: make(chan pr.Message),
		},
		Tree: p.CreateTree(0),
	}

	// Add topic to topic list
	t.TopicsList = append(t.TopicsList, topic_name)

	return t.AllTopics[topic_name]
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

package topics

import (
	pr "treenity/src/protocol"
)

func (t *Topics) Create(topic_name string) {
	// Handle Thread safety
	t.mu.Lock()
	defer t.mu.Unlock()

	// Create topic
	t.AllTopics[topic_name] = &pr.Topic{
		Name:     topic_name,
		Messages: make([]*pr.Message, 0),
	}

	// Add topic to topic list
	t.TopicsList = append(t.TopicsList, topic_name)
}

func (t *Topics) Get(topic_name string) (*pr.Topic, bool) {
	// Handle Thread safety
	t.mu.Lock()
	defer t.mu.Unlock()

	// Get topic
	topic, exist := t.AllTopics[topic_name]

	return topic, exist
}

func (t *Topics) GetTopicNames() []string {
	// Handle Thread safety
	t.mu.Lock()
	defer t.mu.Unlock()

	// Get topic
	topics := t.TopicsList

	return topics
}

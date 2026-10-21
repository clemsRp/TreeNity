package topics

import (
	"sync"
	pr "treenity/src/protocol"
)

type Topics struct {
	AllTopics  map[string]*pr.Topic
	TopicsList []string
	mu         sync.RWMutex
}

func NewTopics() *Topics {
	return &Topics{
		AllTopics:  make(map[string]*pr.Topic),
		TopicsList: make([]string, 0),
	}
}

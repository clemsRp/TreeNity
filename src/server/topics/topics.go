package topics

import (
	"sync"
	pr "treenity/src/protocol"
	p "treenity/src/server/prefix"
)

type TopicDatas struct {
	Topic *pr.Topic
	Tree  *p.Tree
}

type Topics struct {
	AllTopics  map[string]*TopicDatas
	TopicsList []string
	mu         sync.RWMutex
}

func NewTopics() *Topics {
	return &Topics{
		AllTopics:  make(map[string]*TopicDatas),
		TopicsList: make([]string, 0),
	}
}

package server

import (
	pr "treenity/src/protocol"
	h "treenity/src/server/hashmap"
	t "treenity/src/server/topics"
)

type Server struct {
	Entrypoint string
	Topics     *t.Topics
	Clients    *h.Hashmap
}

func NewServer() *Server {
	return &Server{
		Entrypoint: pr.EntryPoint,
		Topics:     t.NewTopics(),
		Clients:    h.CreateHashmap(int(pr.ClientsHashmapSize)),
	}
}

func (s *Server) Start() (int, error) {
	pid := 1337

	return pid, nil
}

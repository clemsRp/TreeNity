package server

import (
	pr "treenity/src/protocol"
	"treenity/src/server/topics"
)

type Server struct {
	Entrypoint string
	Topics     *topics.Topics
}

func NewServer() *Server {
	return &Server{
		Entrypoint: pr.EntryPoint,
		Topics:     topics.NewTopics(),
	}
}

func (s *Server) Start() (int, error) {
	pid := 1337

	return pid, nil
}

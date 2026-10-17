package server

import (
	pr "treenity/src/protocol"
)

type Server struct {
	Entrypoint string
}

func NewServer() *Server {
	return &Server{
		Entrypoint: pr.EntryPoint,
	}
}

func (s *Server) Start() (int, error) {
	pid := 1337

	return pid, nil
}

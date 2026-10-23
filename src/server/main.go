package server

import (
	"fmt"
	"io"
	"os"
	pr "treenity/src/protocol"
	"treenity/src/protocol/ipc"
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

func (s *Server) Start() error {
	// Get FIFO path
	pid := os.Getpid()
	path := fmt.Sprintf("%s%d", pr.EntryPoint, pid)

	fmt.Println(path)

	// Create Server Reader
	reader, err := ipc.NewReader(path)
	if err != nil {
		return err
	}
	defer reader.Close()

	for {
		line, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				continue
			}
			return err
		}

		// TODO: Handle Command
		fmt.Printf("Commande reçue : %s\n", string(line))
	}
}

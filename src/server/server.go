package server

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	pr "treenity/src/protocol"
	"treenity/src/protocol/ipc"
	cmd "treenity/src/server/cmd"
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
		line, read_err := reader.ReadLine()
		if read_err != nil {
			if read_err == io.EOF {
				continue
			}
			return read_err
		}

		// Cast request
		var command pr.ToServer
		json_err := json.Unmarshal(line, &command)
		if json_err != nil {
			return json_err
		}

		// Get request response
		res, datas, err := s.HandleCommands(command)

		// Create Writer
		writer, write_err := ipc.NewWriter(command.Path)
		if write_err != nil {
			return write_err
		}

		// Send response to client
		cmd.SendResponse(writer, command, res, datas, err)
	}
}

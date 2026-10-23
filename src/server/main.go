package server

import (
	"encoding/json"
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

		// Cast request
		var command pr.ToServer
		json_err := json.Unmarshal(line, &command)
		if json_err != nil {
			return json_err
		}

		// Get request response
		res, datas, err := s.HandleCommands(command)

		// Send response to client
		s.SendResponse(command, res, datas, err)
	}
}

func (s *Server) SendResponse(command pr.ToServer, res string, datas any, server_err error) error {
	// Create Writer
	writer, err := ipc.NewWriter(command.Path)
	if err != nil {
		return err
	}

	// Convert cmd_struct to Raw Message
	rawDatas, err := json.Marshal(datas)
	if err != nil {
		return err
	}

	server_err_string := ""
	if server_err != nil {
		server_err_string = server_err.Error()
	}

	// Init response
	response := pr.ToClient{
		Err:      server_err_string,
		Response: res,
		Datas:    json.RawMessage(rawDatas),
	}

	// Cast response
	payload, err := json.Marshal(response)
	if err != nil {
		return err
	}

	// Send request to Server
	return writer.WriteLine(payload)
}

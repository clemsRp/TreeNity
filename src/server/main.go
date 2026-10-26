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

		// Send response to client
		s.SendResponse(command, res, datas, err)
	}
}

func (s *Server) SendResponse(command pr.ToServer, res string, datas any, server_err *pr.Error) *pr.Error {
	// Create Writer
	writer, write_err := ipc.NewWriter(command.Path)
	if write_err != nil {
		return &pr.Error{
			Message: write_err.Error(),
			Code:    3,
		}
	}

	// Convert cmd_struct to Raw Message
	rawDatas, marsh_err := json.Marshal(datas)
	if marsh_err != nil {
		return &pr.Error{
			Message: marsh_err.Error(),
			Code:    1,
		}
	}

	// Init response
	response := pr.ToClient{
		Err:      server_err,
		Response: res,
		Datas:    json.RawMessage(rawDatas),
	}

	// Cast response
	payload, err := json.Marshal(response)
	if err != nil {
		return &pr.Error{
			Message: err.Error(),
			Code:    1,
		}
	}

	// Send request to Server
	send_err := writer.WriteLine(payload)
	if send_err != nil {
		return &pr.Error{
			Message: send_err.Error(),
			Code:    3,
		}
	}

	return nil
}

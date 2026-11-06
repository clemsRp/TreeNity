package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
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

func (s *Server) closeClient() {
	for _, bucket := range s.Clients.Nodes {
		for node := bucket; node != nil; node = node.Next {
			client := node.Value
			if !client.IsActive {
				continue
			}

			if client.Writer == nil {
				writer, err := ipc.NewWriter(client.IpcPath)
				if err != nil {
					continue
				}
				client.Writer = writer
			}

			payload, _ := json.Marshal(pr.CmdMessageDatas{Key: "EOF", Body: "EOF"})

			client.Writer.WriteLine(payload)

			_ = client.Writer.Close()
			client.IsActive = false
		}
	}
}

func (s *Server) Start() (int, error) {
	// Get FIFO path
	pid := os.Getpid()
	path := fmt.Sprintf("%s%d", pr.EntryPoint, pid)

	fmt.Println(path)

	// Create Server Reader
	reader, err := ipc.NewReader(path)
	if err != nil {
		return 1, err
	}
	defer reader.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)

	lineChan := make(chan string)
	errorChan := make(chan error)

	go func()  {
		for {
			line, err := reader.ReadLine()
			if err != nil {
				errorChan <- err
				return
			}
			lineChan <- string(line)
		}
	}()

	go func() (int, error){
		defer wg.Done()

		for {
			select {
				case <-ctx.Done():
					fmt.Println("Server shutingdown")
					s.closeClient()
					if reader != nil {
						reader.Close()
					}
					return 0, err

				case errorGet := <-errorChan:
					if errorGet != nil {
						if errorGet == io.EOF {
							continue
						}
						return 1, errorGet
					}

				case lineGet := <-lineChan:
					// Cast request
					var command pr.ToServer
					json_err := json.Unmarshal([]byte(lineGet), &command)
					if json_err != nil {
						return 1, json_err
					}

					// Get request response
					res, datas, err := s.HandleCommands(command)

					// Create Writer
					writer, write_err := ipc.NewWriter(command.Path)
					if write_err != nil {
						return 1, write_err
					}

					// Send response to client
					cmd.SendResponse(writer, command, res, datas, err)
					writer.Close()
					if command.Command == pr.CmdSubscribe && err == nil && send_err == nil {
						var command_datas pr.CmdSubscribeDatas
						if json_err := json.Unmarshal(command.Datas, &command_datas); json_err == nil {
							cmd.ReplayAndInsert(command_datas, s.Topics, s.Clients)
						}
					}
			}
		}
	}()

	<-ctx.Done()
	wg.Wait()

	return 0, err
}

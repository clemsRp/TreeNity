package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/signal"
	"slices"
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

			SetQuit := pr.ToClient{Quit: true}
			payload, _ := json.Marshal(SetQuit)

			client.Writer.WriteLine(payload)

			_ = client.Writer.Close()
			client.IsActive = false
		}
	}
}

func ReadChannel(reader *ipc.Reader, lineChan chan string, errorChan chan error) {
	go func() {
		for {
			line, err := reader.ReadLine()
			if err != nil {
				errorChan <- err
				return
			}
			lineChan <- string(line)
		}
	}()
}

func (s *Server) RequestAction(reader *ipc.Reader, err error, lineChan chan string, errorChan chan error, ctx context.Context, wg *sync.WaitGroup) {
	go func() error {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				fmt.Fprintln(os.Stderr, "Server shutingdown")
				s.closeClient()
				if reader != nil {
					reader.Close()
				}
				return err

			case errorGet := <-errorChan:
				if errorGet != nil {
					if errorGet == io.EOF {
						continue
					}
					return errorGet
				}

			case lineGet := <-lineChan:
				// Cast request
				var command pr.ToServer
				json_err := json.Unmarshal([]byte(lineGet), &command)
				if json_err != nil {
					continue
				}

				// Get request response
				res, datas, err := s.HandleCommands(command)

				no_response_cmds := []string{pr.CmdAck, pr.CmdDisconnect}
				if slices.Contains(no_response_cmds, command.Command) {
					continue
				}

				// Create Writer
				writer, write_err := ipc.NewWriter(command.Path)
				if write_err != nil {
					fmt.Fprintln(os.Stderr, "client gone:", write_err)
					continue
				}

				// Send response to client
				send_err := cmd.SendResponse(writer, command, res, datas, err)
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

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)

	lineChan := make(chan string)
	errorChan := make(chan error)

	ReadChannel(reader, lineChan, errorChan)
	s.RequestAction(reader, err, lineChan, errorChan, ctx, &wg)

	<-ctx.Done()
	wg.Wait()

	return err
}

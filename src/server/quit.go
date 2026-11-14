package server

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
	pr "treenity/src/protocol"
	"treenity/src/protocol/ipc"
	"treenity/src/server/cmd"
)

func (s *Server) notifyQuit() {
	quit, _ := json.Marshal(pr.ToClient{Quit: true})

	for _, c := range s.Clients.Snapshot() {
		if !c.IsActive || c.Writer == nil {
			continue
		}
		go func(c *pr.Client) {
			if err := c.Writer.WriteLineTimeout(quit, pr.ShutdownTimeout); err != nil {
				c.IsActive = false
			}
		}(c)
	}
}

func (s *Server) activeCount() int {
	n := 0
	for _, c := range s.Clients.Snapshot() {
		if c.IsActive {
			n++
		}
	}
	return n
}

func (s *Server) gracefulShutdown(lineChan <-chan string, errorChan <-chan error) {
	fmt.Fprintln(os.Stderr, "Server shutting down")

	s.Topics.StopAll()
	s.notifyQuit()

	timeout := time.NewTimer(pr.ShutdownTimeout)
	defer timeout.Stop()
	tick := time.NewTicker(50 * time.Millisecond)
	defer tick.Stop()

	for s.activeCount() > 0 {
		select {
		case line := <-lineChan:
			var command pr.ToServer
			if json.Unmarshal([]byte(line), &command) != nil {
				continue
			}
			switch command.Command {
			case pr.CmdAck, pr.CmdDisconnect:
				s.HandleCommands(command)
			default:
				if w, err := ipc.NewWriter(command.Path); err == nil {
					cmd.SendResponse(w, command, "", "",
						&pr.Error{Message: "server shutting down", Code: 3})
					w.Close()
				}
			}
		case <-errorChan:
			return
		case <-tick.C:
		case <-timeout.C:
			fmt.Fprintln(os.Stderr, "shutdown timeout, forcing exit")
			return
		}
	}
}

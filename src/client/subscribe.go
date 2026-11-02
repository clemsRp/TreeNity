package client

import (
	"bufio"
	"fmt"
	"os"
	pr "treenity/src/protocol"
)

func (c *Client) ExecSubscribe(cmd_struct any) *pr.Error {
	command := cmd_struct.(pr.CmdSubscribeDatas)

	clientFifoPath := command.IpcIdentifier + "." + command.SubscriberName

	// Open IPC pipe in read-only mode
	file, err := os.OpenFile(clientFifoPath, os.O_RDONLY, os.ModeNamedPipe)
	if err != nil {
		return &pr.Error{
			Message: fmt.Sprintf("failed to open fifo %s: %s", clientFifoPath, err.Error()),
			Code:    3,
		}
	}
	defer file.Close()

	// Instantiate scanner on opened file
	scanner := bufio.NewScanner(file)

	// Blocking continuous listening loop
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	// Handle scanning errors on loop exit
	if err := scanner.Err(); err != nil {
		return &pr.Error{
			Message: fmt.Sprintf("error reading from fifo %s: %s", clientFifoPath, err.Error()),
			Code:    3,
		}
	}

	return nil
}


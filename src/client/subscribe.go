package client

import (
	"bufio"
	"fmt"
	"os"
	pr "treenity/src/protocol"
)

func (c *Client) ExecSubscribe(cmd_struct any) *pr.Error {
	fifoPath := c.Path

	// 1. Ouverture du FIFO dédié en lecture seule (os.O_RDONLY).
	// L'appel est bloquant sous Linux jusqu'à ce que le serveur ouvre l'extrémité d'écriture.
	file, err := os.OpenFile(fifoPath, os.O_RDONLY, 0)
	if err != nil {
		return &pr.Error{
			Message: fmt.Sprintf("failed to open fifo %s: %s", fifoPath, err.Error()),
			Code:    3,
		}
	}
	defer file.Close()

	// 2. Boucle continue de lecture ligne par ligne avec bufio.Scanner.
	// scanner.Scan() bloque jusqu'à réception d'une nouvelle ligne (\n).
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
	}

	// 3. Gestion d'erreur de lecture
	if err := scanner.Err(); err != nil {
		return &pr.Error{
			Message: fmt.Sprintf("error reading from fifo %s: %s", fifoPath, err.Error()),
			Code:    3,
		}
	}

	return nil
}

package client

import (
	"bufio"
	"encoding/binary"
	"io"
	"os"
	"strings"
	pr "treenity/src/protocol"
)

// Run exécute la production de messages sur un topic.
func (c *Client) Run(cmd pr.CmdProduceDatas) *pr.Error {
	if cmd.Raw {
		// 1. Mode binaire : ignorer la validation texte et la limite de 1024 octets
		reader := bufio.NewReader(os.Stdin)

		for {
			var keyLen int32
			err := binary.Read(reader, binary.LittleEndian, &keyLen)
			if err == io.EOF {
				break
			}
			if err != nil || keyLen < 0 {
				return &pr.Error{Message: pr.ErrPartialRecordEOF, Code: 1}
			}

			keyBuf := make([]byte, keyLen)
			if _, err := io.ReadFull(reader, keyBuf); err != nil {
				return &pr.Error{Message: pr.ErrPartialRecordEOF, Code: 1}
			}

			var valLen int32
			if err := binary.Read(reader, binary.LittleEndian, &valLen); err != nil || valLen < 0 {
				return &pr.Error{Message: pr.ErrPartialRecordEOF, Code: 1}
			}

			valBuf := make([]byte, valLen)
			if _, err := io.ReadFull(reader, valBuf); err != nil {
				return &pr.Error{Message: pr.ErrPartialRecordEOF, Code: 1}
			}

			msgStruct := pr.CmdMessageDatas{
				Key:       string(keyBuf),
				Body:      string(valBuf),
				RawKey:    keyBuf,
				RawBody:   valBuf,
				TopicName: cmd.TopicName,
			}
			if sendErr := c.SendRequest(msgStruct, pr.CmdMessage); sendErr != nil {
				return sendErr
			}
		}
		return nil
	}

	// 2. Mode texte par défaut
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return &pr.Error{Message: pr.ErrInvalidFormat, Code: 1}
		}
		key, body := parts[0], parts[1]

		if len(key)+len(body) > pr.MaxMessageLength {
			return &pr.Error{Message: pr.ErrMessageTooLong, Code: 1}
		}

		msgStruct := pr.CmdMessageDatas{
			Key:       key,
			Body:      body,
			TopicName: cmd.TopicName,
		}
		if sendErr := c.SendRequest(msgStruct, pr.CmdMessage); sendErr != nil {
			return sendErr
		}
	}
	if scanErr := scanner.Err(); scanErr != nil {
		return &pr.Error{Message: scanErr.Error(), Code: 1}
	}
	return nil
}

// ExecProduceText maintient la compatibilité avec l'ancien appel
func (c *Client) ExecProduceText(topic_name string) *pr.Error {
	return c.Run(pr.CmdProduceDatas{
		TopicName: topic_name,
		Raw:       false,
	})
}

// ExecProduceRaw maintient la compatibilité avec l'ancien appel
func (c *Client) ExecProduceRaw(topic_name string) *pr.Error {
	return c.Run(pr.CmdProduceDatas{
		TopicName: topic_name,
		Raw:       true,
	})
}

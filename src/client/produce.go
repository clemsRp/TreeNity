package client

import (
	"bufio"
	"context"
	"encoding/binary"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"
	pr "treenity/src/protocol"
)

func (c *Client) RawMode(cmd pr.CmdProduceDatas) *pr.Error {
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

		if int(keyLen)+int(valLen) > pr.MaxMessageLength {
			return &pr.Error{Message: pr.ErrMessageTooLong, Code: 1}
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

func (c *Client) TextMode(cmd pr.CmdProduceDatas) *pr.Error {
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

func (c *Client) SignalThread(ctx context.Context) {
	go func ()  {
		<-ctx.Done()
		c.Disconnect()
		os.Exit(0)
	}()
}

func (c *Client) Run(cmd pr.CmdProduceDatas) *pr.Error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	c.SignalThread(ctx)

	if cmd.Raw {
		return c.RawMode(cmd)
	} else {
		return c.TextMode(cmd)
	}
}

// Backward compatibility
func (c *Client) ExecProduceText(topic_name string) *pr.Error {
	return c.Run(pr.CmdProduceDatas{
		TopicName: topic_name,
		Raw:       false,
	})
}

// Backward compatibility
func (c *Client) ExecProduceRaw(topic_name string) *pr.Error {
	return c.Run(pr.CmdProduceDatas{
		TopicName: topic_name,
		Raw:       true,
	})
}

package client

import (
	"bufio"
	"encoding/binary"
	"io"
	"os"
	"strings"
	pr "treenity/src/protocol"
)

func (c *Client) ExecProduceText(topic_name string) *pr.Error {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return &pr.Error{
				Message: pr.ErrInvalidFormat,
				Code:    1,
			}
		}
		key, body := parts[0], parts[1]

		msg_err := c.SendMessage(key, body, topic_name)
		if msg_err != nil {
			return msg_err
		}
	}

	if scan_err := scanner.Err(); scan_err != nil {
		return &pr.Error{
			Message: scan_err.Error(),
			Code:    1,
		}
	}

	return nil
}

func (c *Client) ExecProduceRaw(topic_name string) *pr.Error {
	reader := bufio.NewReader(os.Stdin)

	for {
		key, key_err := c.GetRawData(reader, true)
		if key_err != nil {
			if key_err.Code == -1 {
				break
			}
			return key_err
		}

		body, body_err := c.GetRawData(reader, false)
		if body_err != nil {
			if body_err.Code == -1 {
				break
			}
			return body_err
		}

		msg_err := c.SendMessage(key, body, topic_name)
		if msg_err != nil {
			return msg_err
		}
	}

	return nil
}

func (c *Client) GetRawData(reader *bufio.Reader, isKey bool) (string, *pr.Error) {
	var dataLen int32
	err := binary.Read(reader, binary.LittleEndian, &dataLen)

	if err == io.EOF {
		if isKey {
			return "", &pr.Error{
				Message: pr.ErrEOF,
				Code:    -1,
			}
		}

		return "", &pr.Error{
			Message: pr.ErrPartialRecordEOF,
			Code:    1,
		}
	}

	if err != nil || dataLen < 0 {
		return "", &pr.Error{
			Message: pr.ErrPartialRecordEOF,
			Code:    1,
		}
	}

	dataBuf := make([]byte, dataLen)
	if _, err := io.ReadFull(reader, dataBuf); err != nil {
		return "", &pr.Error{
			Message: pr.ErrPartialRecordEOF,
			Code:    1,
		}
	}

	return string(dataBuf), nil
}

func (c *Client) SendMessage(key, body, topic_name string) *pr.Error {
	if len(key)+len(body) > pr.MaxMessageLength {
		return &pr.Error{
			Message: pr.ErrMessageTooLong,
			Code:    1,
		}
	}
	msg_struct := pr.CmdMessageDatas{
		Key:       key,
		Body:      body,
		TopicName: topic_name,
	}

	err := c.SendRequest(msg_struct, pr.CmdMessage)
	if err != nil {
		return err
	}

	return nil
}

package client

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	pr "treenity/src/protocol"
	"treenity/src/protocol/ipc"
)

func (c *Client) ExecSubscribe(cmd_struct any, permReader ...*ipc.Reader) *pr.Error {
	cmd := cmd_struct.(pr.CmdSubscribeDatas)

	var reader *ipc.Reader
	if len(permReader) > 0 && permReader[0] != nil {
		reader = permReader[0]
	} else {
		permanentPath := cmd.IpcIdentifier + "." + cmd.SubscriberName
		var err error
		reader, err = ipc.NewReader(permanentPath)
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("Server closed the connection")
				return nil
			}
			return &pr.Error{
				Message: err.Error(),
				Code:    1,
			}
		}
	}
	defer reader.Close()

	for {
		data_read, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return &pr.Error{
				Message: err.Error(),
				Code:    1,
			}
		}

		if len(data_read) == 0 {
			continue
		}

		var data_convert pr.CmdMessageDatas
		err = json.Unmarshal(data_read, &data_convert)
		if err != nil {
			return &pr.Error{
				Message: err.Error(),
				Code:    1,
			}
		}

		if data_convert.Key == "EOF" && data_convert.Body == "EOF" {
			return nil
		}


		if !cmd.Raw {
			fmt.Printf("%s:%s\n", data_convert.Key, data_convert.Body)
		} else {
			keyBytes := data_convert.RawKey
			if len(keyBytes) == 0 {
				keyBytes = []byte(data_convert.Key)
			}
			bodyBytes := data_convert.RawBody
			if len(bodyBytes) == 0 {
				bodyBytes = []byte(data_convert.Body)
			}

			buffer := new(bytes.Buffer)
			binary.Write(buffer, binary.LittleEndian, data_convert.Offset)
			binary.Write(buffer, binary.LittleEndian, int32(len(keyBytes)))
			buffer.Write(keyBytes)
			binary.Write(buffer, binary.LittleEndian, int32(len(bodyBytes)))
			buffer.Write(bodyBytes)

			if _, err := os.Stdout.Write(buffer.Bytes()); err != nil {
				return &pr.Error{
					Message: err.Error(),
					Code:    1,
				}
			}
		}
	}
}

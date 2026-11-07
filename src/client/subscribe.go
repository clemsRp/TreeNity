package client

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	pr "treenity/src/protocol"
)

func PrintData(dataMsg pr.CmdMessageDatas, cmd pr.CmdSubscribeDatas) {
		if !cmd.Raw {
			fmt.Printf("%s:%s\n", dataMsg.Key, dataMsg.Body)
		} else {
			keyByte := []byte(dataMsg.Key)
			bodyByte := []byte(dataMsg.Body)
			buffer := new(bytes.Buffer)
			binary.Write(buffer, binary.LittleEndian, int32(len(keyByte)))
			buffer.Write(keyByte)
			binary.Write(buffer, binary.LittleEndian, int32(len(bodyByte)))
			buffer.Write(bodyByte)

			fmt.Println(buffer.String())
		}

}

func (c *Client) ExecSubscribe(cmd_struct any) *pr.Error {
	cmd := cmd_struct.(pr.CmdSubscribeDatas)

	for {
		data_read, err := c.Reader.ReadLine()

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
		var data_convert pr.ToClient
		err = json.Unmarshal(data_read, &data_convert)

		if err != nil {
			return &pr.Error{
				Message: err.Error(),
				Code:    1,
			}
		}

		if data_convert.Quit {
			return nil
		}

		var dataMsg pr.CmdMessageDatas
		msgError := json.Unmarshal(data_convert.Datas, &dataMsg)
		if msgError != nil {
			return &pr.Error{
				Message: msgError.Error(),
				Code: 1,
			}
		}

		PrintData(dataMsg, cmd)
	}
}

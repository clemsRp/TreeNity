package client

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	pr "treenity/src/protocol"
)

func (c *Client) ExecSubscribe(cmd_struct any) *pr.Error {
	cmd := cmd_struct.(pr.CmdSubscribeDatas)

	for {
		data_read, err := c.Reader.ReadLine()
		if err != nil {
			return &pr.Error{
				Message: err.Error(),
				Code:    1,
			}
		}
		var data_convert pr.CmdMessageDatas
		err = json.Unmarshal(data_read, &data_convert)
		if err != nil {
			return &pr.Error{
				Message: err.Error(),
				Code:    1,
			}
		}

		fmt.Println(cmd.Raw)
		if !cmd.Raw {
			fmt.Printf("%s:%s\n", data_convert.Key, data_convert.Body)
		} else {
			keyByte := []byte(data_convert.Key)
			bodyByte := []byte(data_convert.Body)
			buffer := new(bytes.Buffer)
			binary.Write(buffer, binary.LittleEndian, int32(len(keyByte)))
			buffer.Write(keyByte)
			binary.Write(buffer, binary.LittleEndian, int32(len(bodyByte)))
			buffer.Write(bodyByte)

			fmt.Println(buffer.String())
		}
	}
}

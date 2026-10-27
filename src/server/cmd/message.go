package cmd

import (
	"fmt"
	pr "treenity/src/protocol"
	t "treenity/src/server/topics"
)

func HandleCmdMessage(command pr.CmdMessageDatas, topics *t.Topics) (string, any, *pr.Error) {
	// TODO: New message logic
	fmt.Println(command.Key, command.Body)

	return "", "", nil
}

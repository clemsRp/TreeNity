package cmd

import (
	"fmt"
	pr "treenity/src/protocol"
	h "treenity/src/server/hashmap"
)

func HandleCmdAck(command pr.CmdAckDatas, clients *h.Hashmap) (string, any, *pr.Error) {
	// Handle invalid subscriber name
	client, err := clients.Get(command.SubscriberName)
	if err != nil {
		return "", "", &pr.Error{
			Message: fmt.Sprintf("'%s' %s", command.SubscriberName, pr.ErrSubscriberDoesNotExist),
			Code:    2,
		}
	}

	// Update client Offset
	client.Offset = command.Offset

	return "", "", nil
}

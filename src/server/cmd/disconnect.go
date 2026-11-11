package cmd

import (
	"fmt"
	pr "treenity/src/protocol"
	h "treenity/src/server/hashmap"
)

func HandleCmdDisconnect(command pr.CmdDisconnectDatas, clients *h.Hashmap) (string, any, *pr.Error) {
	// Handle invalid subscriber name
	client, err := clients.Get(command.SubscriberName)
	if err != nil {
		return "", "", &pr.Error{
			Message: fmt.Sprintf("'%s' %s", command.SubscriberName, pr.ErrSubscriberDoesNotExist),
			Code:    2,
		}

	} else if !client.IsActive {
		return "", "", &pr.Error{
			Message: fmt.Sprintf("'%s' %s", command.SubscriberName, pr.ErrSubscriberAlreadyActive),
			Code:    2,
		}
	}

	// Update client status
	client.IsActive = false

	return "", "", nil
}

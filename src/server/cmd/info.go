package cmd

import (
	"fmt"
	pr "treenity/src/protocol"
	h "treenity/src/server/hashmap"
)

func HandleCmdInfo(command pr.CmdInfoDatas, clients *h.Hashmap) (string, any, *pr.Error) {
	// Handle invalid name
	if !IsValidId(command.SubscriberName) {
		return "", nil, &pr.Error{
			Message: fmt.Sprintf("'%s' %s", command.SubscriberName, pr.ErrInvalidName),
			Code:    1,
		}
	}

	// Get Client
	client, err := clients.Get(command.SubscriberName)

	// Handle invalid subscriber name
	if err != nil {
		return "", nil, &pr.Error{
			Message: err.Error(),
			Code:    2,
		}
	}

	return "", client, nil
}

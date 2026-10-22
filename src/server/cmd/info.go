package cmd

import (
	pr "treenity/src/protocol"
	h "treenity/src/server/hashmap"
)

func HandleCmdInfo(command pr.CmdInfoDatas, clients *h.Hashmap) (string, any, error) {
	// Get Client
	client, err := clients.Get(command.SubscriberName)

	// Handle invalid subscriber name
	if err != nil {
		return "", nil, err
	}

	return "", client, nil
}

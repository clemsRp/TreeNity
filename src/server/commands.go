package server

import (
	"encoding/json"
	"errors"
	pr "treenity/src/protocol"
	"treenity/src/server/cmd"
)

func (s *Server) HandleCommands(command pr.ToServer) (string, any, error) {
	// Define return variables
	var res string
	var datas any
	var err error

	// Handle Command
	switch command.Command {
	case pr.CmdCreate:
		var command_datas pr.CmdCreateDatas

		// Cast datas
		err := json.Unmarshal(command.Datas, &command_datas)
		if err != nil {
			return "", "", errors.New(pr.ErrCast)
		}

		// Run handler
		res, datas, err = cmd.HandleCmdCreate(command_datas, s.Topics)

	case pr.CmdList:
		var command_datas pr.CmdListDatas

		// Cast datas
		err := json.Unmarshal(command.Datas, &command_datas)
		if err != nil {
			return "", "", errors.New(pr.ErrCast)
		}

		// Run handler
		res, datas, err = cmd.HandleCmdList(command_datas, s.Topics)

	case pr.CmdInfo:
		var command_datas pr.CmdInfoDatas

		// Cast datas
		err := json.Unmarshal(command.Datas, &command_datas)
		if err != nil {
			return "", "", errors.New(pr.ErrCast)
		}

		// Run handler
		res, datas, err = cmd.HandleCmdInfo(command_datas, s.Clients)

	case pr.CmdProduce:
		var command_datas pr.CmdProduceDatas

		// Cast datas
		err := json.Unmarshal(command.Datas, &command_datas)
		if err != nil {
			return "", "", errors.New(pr.ErrCast)
		}

		// Run handler
		res, datas, err = cmd.HandleCmdProduce(command_datas, s.Topics)

	case pr.CmdSubscribe:
		var command_datas pr.CmdSubscribeDatas

		// Cast datas
		err := json.Unmarshal(command.Datas, &command_datas)
		if err != nil {
			return "", "", errors.New(pr.ErrCast)
		}

		// Run handler
		res, datas, err = cmd.HandleCmdSubscribe(command_datas, s.Topics, s.Clients)
	}

	return res, datas, err
}

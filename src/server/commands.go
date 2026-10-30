package server

import (
	"encoding/json"
	pr "treenity/src/protocol"
	"treenity/src/server/cmd"
)

func (s *Server) HandleCommands(command pr.ToServer) (string, any, *pr.Error) {
	// Define return variables
	var res string
	var datas any
	var err *pr.Error

	// Handle Command
	switch command.Command {
	case pr.CmdCreate:
		var command_datas pr.CmdCreateDatas

		// Cast datas
		json_err := json.Unmarshal(command.Datas, &command_datas)
		if json_err != nil {
			return "", "", &pr.Error{
				Message: pr.ErrCast,
				Code:    1,
			}
		}

		// Run handler
		res, datas, err = cmd.HandleCmdCreate(command_datas, s.Topics, s.Clients)

	case pr.CmdList:
		var command_datas pr.CmdListDatas

		// Cast datas
		json_err := json.Unmarshal(command.Datas, &command_datas)
		if json_err != nil {
			return "", "", &pr.Error{
				Message: pr.ErrCast,
				Code:    1,
			}
		}

		// Run handler
		res, datas, err = cmd.HandleCmdList(command_datas, s.Topics)

	case pr.CmdInfo:
		var command_datas pr.CmdInfoDatas

		// Cast datas
		json_err := json.Unmarshal(command.Datas, &command_datas)
		if json_err != nil {
			return "", "", &pr.Error{
				Message: pr.ErrCast,
				Code:    1,
			}
		}

		// Run handler
		res, datas, err = cmd.HandleCmdInfo(command_datas, s.Clients)

	case pr.CmdProduce:
		var command_datas pr.CmdProduceDatas

		// Cast datas
		json_err := json.Unmarshal(command.Datas, &command_datas)
		if json_err != nil {
			return "", "", &pr.Error{
				Message: pr.ErrCast,
				Code:    1,
			}
		}

		// Run handler
		res, datas, err = cmd.HandleCmdProduce(command_datas, s.Topics)

	case pr.CmdSubscribe:
		var command_datas pr.CmdSubscribeDatas

		// Cast datas
		json_err := json.Unmarshal(command.Datas, &command_datas)
		if json_err != nil {
			return "", "", &pr.Error{
				Message: pr.ErrCast,
				Code:    1,
			}
		}

		// Run handler
		res, datas, err = cmd.HandleCmdSubscribe(command_datas, s.Topics, s.Clients, command.Path)

	case pr.CmdMessage:
		var command_datas pr.CmdMessageDatas

		// Cast datas
		json_err := json.Unmarshal(command.Datas, &command_datas)
		if json_err != nil {
			return "", "", &pr.Error{
				Message: pr.ErrCast,
				Code:    1,
			}
		}

		// Run handler
		res, datas, err = cmd.HandleCmdMessage(command_datas, s.Topics)
	}

	return res, datas, err
}

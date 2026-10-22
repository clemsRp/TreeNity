package server

import (
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
		// Cast datas
		command_datas, ok := command.Datas.(pr.CmdCreateDatas)
		if !ok {
			return "", "", errors.New("cast error")
		}

		// Run handler
		res, datas, err = cmd.HandleCmdCreate(command_datas, s.Topics)

	case pr.CmdList:
		// Cast datas
		command_datas, ok := command.Datas.(pr.CmdListDatas)
		if !ok {
			return "", "", errors.New("cast error")
		}

		// Run handler
		res, datas, err = cmd.HandleCmdList(command_datas, s.Topics)

	case pr.CmdInfo:
		// Cast datas
		command_datas, ok := command.Datas.(pr.CmdInfoDatas)
		if !ok {
			return "", "", errors.New("cast error")
		}

		// Run handler
		res, datas, err = cmd.HandleCmdInfo(command_datas, s.Clients)

	case pr.CmdProduce:
		// Cast datas
		command_datas, ok := command.Datas.(pr.CmdProduceDatas)
		if !ok {
			return "", "", errors.New("cast error")
		}

		// Run handler
		res, datas, err = cmd.HandleCmdProduce(command_datas)

	case pr.CmdSubscribe:
		// Cast datas
		command_datas, ok := command.Datas.(pr.CmdSubscribeDatas)
		if !ok {
			return "", "", errors.New("cast error")
		}

		// Run handler
		res, datas, err = cmd.HandleCmdSubscribe(command_datas, s.Topics, s.Clients)
	}

	return res, datas, err
}

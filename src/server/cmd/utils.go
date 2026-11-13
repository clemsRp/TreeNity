package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"syscall"
	pr "treenity/src/protocol"
	"treenity/src/protocol/ipc"
	t "treenity/src/server/topics"
)

var CheckIdRegex = regexp.MustCompile(pr.ValidRegex)

func IsValidId(id string) bool {
	return CheckIdRegex.MatchString(id)
}

func IsBrokenPipe(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, syscall.EPIPE) ||
		errors.Is(err, syscall.ECONNRESET) ||
		errors.Is(err, io.ErrClosedPipe)
}

func SendResponse(writer *ipc.Writer, command pr.ToServer, res string, datas any, server_err *pr.Error) *pr.Error {
	// Convert cmd_struct to Raw Message
	rawDatas, marsh_err := json.Marshal(datas)
	if marsh_err != nil {
		return &pr.Error{
			Message: marsh_err.Error(),
			Code:    1,
		}
	}

	// Init response
	response := pr.ToClient{
		Err:      server_err,
		Response: res,
		Datas:    json.RawMessage(rawDatas),
	}

	// Cast response
	payload, err := json.Marshal(response)
	if err != nil {
		return &pr.Error{
			Message: err.Error(),
			Code:    1,
		}
	}

	// Send request to Server
	send_err := writer.WriteLine(payload)
	if send_err != nil {
		return &pr.Error{
			Message: send_err.Error(),
			Code:    3,
		}
	}

	return nil
}

func clientAlive(path string) bool {
	f, err := os.OpenFile(path, os.O_WRONLY|syscall.O_NONBLOCK, 0)
	if err != nil {
		return false
	}
	f.Close()
	return true
}

func CheckSubscribeCmd(command pr.CmdSubscribeDatas, topics *t.Topics, ipcPath string) *pr.Error {
	// Check response channel
	if ipcPath == "" {
		return &pr.Error{
			Message: "missing client IPC path",
			Code:    1,
		}
	}

	// Check subscriber is valid name
	if !IsValidId(command.SubscriberName) {
		return &pr.Error{
			Message: fmt.Sprintf("'%s' %s", command.SubscriberName, pr.ErrInvalidName),
			Code:    1,
		}
	}
	// Check topic is valid name
	if !IsValidId(command.TopicName) {
		return &pr.Error{
			Message: fmt.Sprintf("'%s' %s", command.TopicName, pr.ErrInvalidName),
			Code:    1,
		}
	}

	// Check topic exist
	_, exist := topics.Get(command.TopicName)
	if !exist {
		return &pr.Error{
			Message: fmt.Sprintf("'%s' %s", command.TopicName, pr.ErrTopicDoesNotExist),
			Code:    2,
		}
	}

	return nil
}

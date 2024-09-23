package protocol

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/nghiagtgsoft/gosf-socketio/color"
	"github.com/nghiagtgsoft/gosf-socketio/logger"
)

var (
	ErrorWrongMessageType = errors.New("Wrong message type")
	ErrorWrongPacket      = errors.New("Wrong packet")
)

func typeToText(msgType int) (string, error) {
	return "", nil
}

func Encode(msg *Message) (string, error) {
	output := ""
	output += strconv.Itoa(int(msg.EngineIoType))

	if msg.SocketType == SocketMessageTypeNone {
		return output, nil
	}
	output += strconv.Itoa(int(msg.SocketType))
	if !(msg.SocketEvent.EmitName == "" || msg.SocketEvent.EmitContent == "") {
		json, _ := json.Marshal([2]interface{}{msg.SocketEvent.EmitName, msg.SocketEvent.EmitContent})
		output += string(json)
	}

	return output, nil
}

func MustEncode(msg *Message) string {
	result, err := Encode(msg)
	if err != nil {
		panic(err)
	}

	return result
}

func GetEngineMessageType(data string) (EngineMessageType, error) {
	if len(data) == 0 {
		return 0, ErrorWrongMessageType
	}
	msgType, _ := strconv.Atoi(data[0:1])
	if msgType > 4 {
		return 0, ErrorWrongMessageType
	}
	return EngineMessageType(msgType), nil
}

func GetSocketMessageType(data string) (SocketMessageType, error) {
	if len(data) == 0 {
		return 0, ErrorWrongMessageType
	}
	msgType, _ := strconv.Atoi(data[1:2])
	if msgType > 6 {
		return 0, ErrorWrongMessageType
	}
	return SocketMessageType(msgType), nil
}

func GetSocketIoEmitName(data string) string {
	jsonevent := data[2:]
	var emit []interface{}
	err := json.Unmarshal([]byte(jsonevent), &emit)
	if err != nil {
		logger.LogErrorSocketIo(color.Red + "Error: " + err.Error() + color.Reset)
	}
	emitNameString, _ := emit[0].(string)
	return emitNameString
}

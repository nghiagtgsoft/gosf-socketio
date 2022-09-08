package protocol

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"

	"github.com/ambelovsky/gosf-socketio/color"
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

func getEngineMessageType(data string) (EngineMessageType, error) {
	if len(data) == 0 {
		return 0, ErrorWrongMessageType
	}
	msgType, _ := strconv.Atoi(data[0:1])
	if msgType > 4 {
		return 0, ErrorWrongMessageType
	}
	return EngineMessageType(msgType), nil
}

func getSocketMessageType(data string) (SocketMessageType, error) {
	if len(data) == 0 {
		return 0, ErrorWrongMessageType
	}
	msgType, _ := strconv.Atoi(data[1:2])
	if msgType > 6 {
		return 0, ErrorWrongMessageType
	}
	return SocketMessageType(msgType), nil
}

/**
Get ack id of current packet, if present
*/
// func getAck(text string) (ackId int, restText string, err error) {
// 	if len(text) < 4 {
// 		return 0, "", ErrorWrongPacket
// 	}
// 	text = text[2:]

// 	pos := strings.IndexByte(text, '[')
// 	if pos == -1 {
// 		return 0, "", ErrorWrongPacket
// 	}

// 	ack, err := strconv.Atoi(text[0:pos])
// 	if err != nil {
// 		return 0, "", err
// 	}

// 	return ack, text[pos:], nil
// }

/**
Get message method of current packet, if present
*/
// func getMethod(text string) (method, restText string, err error) {

// }

func Decode(data string) (*Message, error) {

	msg := &Message{}
	var err error

	msg.EngineIoType, err = getEngineMessageType(data)
	log.Println(color.Green + "Engine IO type: (" + data[0:1] + ") " + msg.EngineIoType.String() + color.Reset)

	if msg.EngineIoType == EngineMessageTypeMessage {
		msg.SocketType, err = getSocketMessageType(data)
		log.Println(color.Yellow + "Socket IO type: (" + data[1:2] + ") " + msg.SocketType.String() + color.Reset)
		if msg.SocketType == SocketMessageTypeEvent {
			jsonevent := data[2:]

			var emit []string
			json.Unmarshal([]byte(jsonevent), &emit)

			msg.SocketEvent = SocketEvent{EmitName: emit[0], EmitContent: emit[1]}
		}
	} else {
		msg.SocketType = SocketMessageTypeNone
		log.Println(color.Yellow + "Socket IO type: " + msg.SocketType.String() + color.Reset)
	}

	if err != nil {
		return nil, err
	}
	return msg, nil

}

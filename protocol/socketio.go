package protocol

import (
	"errors"
	"log"
	"strconv"
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
	output += strconv.Itoa(int(msg.SocketType))

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
	log.Println("Received msg of type", data[0:1])
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
	log.Println("Received msg of type", data[0:1])
	msgType, _ := strconv.Atoi(data[0:1])
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

	log.Println("Message type:", data[0:1])
	msg := &Message{}
	var err error

	msg.EngineIoType, err = getEngineMessageType(data)
	msg.SocketType, err = getSocketMessageType(data)

	if err != nil {
		return nil, err
	}
	log.Println("Got message of Engine type ", msg.EngineIoType.String(), " Socket type is, ", msg.SocketType.String())
	return msg, nil
	// var err error
	// msg := &Message{}
	// msg.Source = data

	// msg.Type, err = getMessageType(data)
	// if err != nil {
	// 	return nil, err
	// }

	// if msg.Type == MessageTypeOpen {
	// 	msg.Args = data[1:]
	// 	return msg, nil
	// }

	// if msg.Type == MessageTypeClose || msg.Type == MessageTypePing ||
	// 	msg.Type == MessageTypePong || msg.Type == MessageTypeEmpty {
	// 	return msg, nil
	// }

	// // ack, rest, err := getAck(data)
	// // msg.AckId = ack
	// // if msg.Type == MessageTypeAckResponse {
	// // 	if err != nil {
	// // 		return nil, err
	// // 	}
	// // 	msg.Args = rest[1 : len(rest)-1]
	// // 	return msg, nil
	// // }

	// if err != nil {
	// 	msg.Type = MessageTypeEmit
	// 	rest = data[2:]
	// }

	// msg.Method, msg.Args, err = getMethod(rest)
	// if err != nil {
	// 	return nil, err
	// }

	return nil, nil
}

package protocol

/*
https://github.com/socketio/socket.io-protocol#0---connect
*/
type EngineMessageType int64

const (
	EngineMessageTypeOpen    EngineMessageType = 0
	EngineMessageTypeClose   EngineMessageType = 1
	EngineMessageTypePing    EngineMessageType = 2
	EngineMessageTypePong    EngineMessageType = 3
	EngineMessageTypeMessage EngineMessageType = 4
)

/*https://github.com/socketio/socket.io-protocol#examples*/

type SocketMessageType int64

const (
	SocketMessageTypeConnect      SocketMessageType = 0
	SocketMessageTypeDisconnect   SocketMessageType = 1
	SocketMessageTypeEvent        SocketMessageType = 2
	SocketMessageTypeAck          SocketMessageType = 3
	SocketMessageTypeConnectError SocketMessageType = 4
	SocketMessageTypeBinaryEvent  SocketMessageType = 5
	SocketMessageTypeBinaryAck    SocketMessageType = 6
)

func (e EngineMessageType) String() string {
	switch e {
	case
		EngineMessageTypeOpen:
		return "EngineMessageTypeOpen"
	case
		EngineMessageTypeClose:
		return "EngineMessageTypeClose"
	case
		EngineMessageTypePing:
		return "EngineMessageTypeClose"
	case
		EngineMessageTypePong:
		return "EngineMessageTypeClose"
	case
		EngineMessageTypeMessage:
		return "EngineMessageTypeClose"
	}
	return "Unknown"
}

func (e SocketMessageType) String() string {
	switch e {

	case SocketMessageTypeConnect:
		return "SocketMessageTypeConnect"
	case SocketMessageTypeDisconnect:
		return "SocketMessageTypeDisconnect"
	case SocketMessageTypeEvent:
		return "SocketMessageTypeEvent"
	case SocketMessageTypeAck:
		return "SocketMessageTypeAck"
	case SocketMessageTypeConnectError:
		return "SocketMessageTypeConnectError"
	case SocketMessageTypeBinaryEvent:
		return "SocketMessageTypeBinaryEvent"
	case SocketMessageTypeBinaryAck:
		return "SocketMessageTypeBinaryAck"
	}
	return "Unknown"
}

type Message struct {
	EngineIoType EngineMessageType
	SocketType   SocketMessageType
}

package gosocketio

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/ambelovsky/gosf-socketio/color"
	"github.com/ambelovsky/gosf-socketio/logger"
	"github.com/ambelovsky/gosf-socketio/protocol"
)

const (
	OnConnection    = "connection"
	OnDisconnection = "disconnection"
	OnError         = "error"
)

/*
*
System handler function for internal event processing
*/
type systemHandler func(c *Channel)

/*
*
Contains maps of message processing functions
*/
type methods struct {
	messageHandlers     sync.Map
	messageHandlersLock sync.RWMutex

	onConnection    systemHandler
	onDisconnection systemHandler
}

/*
*
create messageHandlers map
*/
func (m *methods) initMethods() {
	//m.messageHandlers = make(sync.Map)
	m.messageHandlers = sync.Map{}

}

/*
*
Add message processing function, and bind it to given method
*/
func (m *methods) On(method string, f interface{}) error {
	c, err := newCaller(f)
	if err != nil {
		return err
	}

	m.messageHandlers.Store(method, c)
	return nil
}

/*
*
Find message processing function associated with given method
*/
func (m *methods) findMethod(method string) (*caller, bool) {
	if f, ok := m.messageHandlers.Load(method); ok {
		return f.(*caller), true
	}

	return nil, false
}

/*
*
Check incoming message
On ack_resp - look for waiter
On ack_req - look for processing function and send ack_resp
On emit - look for processing function
*/
func (m *methods) processIncomingMessage(c *Channel, engineIoType protocol.EngineMessageType, pkg string) {
	switch engineIoType {
	case protocol.EngineMessageTypeOpen:
		m.processOpenMessage(c, pkg)
	case protocol.EngineMessageTypePing:
		m.processPingMessage(c)
	case protocol.EngineMessageTypeMessage:
		m.processSocketMessage(c, pkg)
	case protocol.EngineMessageTypeClose:
		m.processDisconnectMessage(c)

	}
}

func (m *methods) processSocketMessage(c *Channel, pkg string) {
	socketType, err := protocol.GetSocketMessageType(pkg)
	if err != nil {
		return
	}
	logger.LogDebugSocketIo(color.Yellow + "Socket IO type: (" + socketType.String() + color.Reset)
	if socketType == protocol.SocketMessageTypeEvent { //Decode socket.io message type
		emitName := protocol.GetSocketIoEmitName(pkg)
		if emitName == "" {
			return
		}
		f, ok := m.findMethod(emitName)
		if !ok {
			return
		}
		if !f.ArgsPresent {
			f.callFunc(c, &struct{}{})
			return
		}
		data := f.getArgs()
		structReceived := []interface{}{"", data}
		jsonevent := pkg[2:]
		err := json.Unmarshal([]byte(jsonevent), &structReceived)
		if err != nil {
			logger.LogDebug(color.Red + "Error: " + err.Error() + color.Reset)
		}
		if err != nil {
			// if reflect.TypeOf(data) == reflect.TypeOf(&data) { //check if it is ok without JSON encoding, mostly for strings
			// 	data = &data
			// 	f.callFunc(c, data)
			// 	return
			// }
			logger.LogErrorSocketIo(color.Red + "Decoding error: " + err.Error() + color.Reset)
		}
		f.callFunc(c, data)
	}

}

func (m *methods) processPingMessage(c *Channel) {

	reply := protocol.Message{}
	reply.EngineIoType = protocol.EngineMessageTypePong
	reply.SocketType = protocol.SocketMessageTypeNone
	command, _ := protocol.Encode(&reply)
	send(command, c)

}
func (m *methods) processOpenMessage(c *Channel, pkg string) {

	c.SetHeader(pkg[1:])

	reply := protocol.Message{}
	reply.EngineIoType = protocol.EngineMessageTypeMessage
	reply.SocketType = protocol.SocketMessageTypeConnect

	command, _ := protocol.Encode(&reply)
	send(command, c)
	time.Sleep(100 * time.Millisecond) //just to be sure the open message is processed
	f, _ := m.findMethod(OnConnection)
	if f == nil {
		return
	}
	f.callFunc(c, &struct{}{})

}
func (m *methods) processDisconnectMessage(c *Channel) {
	m.onDisconnection(c)

}

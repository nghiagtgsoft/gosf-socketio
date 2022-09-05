package gosocketio

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/ambelovsky/gosf-socketio/color"
	"github.com/ambelovsky/gosf-socketio/protocol"
)

const (
	OnConnection    = "connection"
	OnDisconnection = "disconnection"
	OnError         = "error"
)

/**
System handler function for internal event processing
*/
type systemHandler func(c *Channel)

/**
Contains maps of message processing functions
*/
type methods struct {
	messageHandlers     sync.Map
	messageHandlersLock sync.RWMutex

	onConnection    systemHandler
	onDisconnection systemHandler
}

/**
create messageHandlers map
*/
func (m *methods) initMethods() {
	//m.messageHandlers = make(sync.Map)
}

/**
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

/**
Find message processing function associated with given method
*/
func (m *methods) findMethod(method string) (*caller, bool) {
	if f, ok := m.messageHandlers.Load(method); ok {
		return f.(*caller), true
	}

	return nil, false
}

/**
Check incoming message
On ack_resp - look for waiter
On ack_req - look for processing function and send ack_resp
On emit - look for processing function
*/
func (m *methods) processIncomingMessage(c *Channel, msg *protocol.Message) {
	switch msg.EngineIoType {
	case protocol.EngineMessageTypeOpen:
		m.processOpenMessage(c, msg)
	case protocol.EngineMessageTypePing:
		m.processPingMessage(c, msg)
	case protocol.EngineMessageTypeMessage:
		m.processSocketMessage(c, msg)
	case protocol.EngineMessageTypeClose:
		m.processDisconnectMessage(c, msg)

	}
}
func (m *methods) processSocketMessage(c *Channel, msg *protocol.Message) {
	if msg.SocketEvent.EmitName == "" {
		return
	}
	log.Println(color.Cyan + "MSG: " + msg.SocketEvent.EmitName + " VALUE: " + msg.SocketEvent.EmitContent + color.Reset)
	f, ok := m.findMethod(msg.SocketEvent.EmitName)
	if !ok {
		return
	}

	if !f.ArgsPresent {
		f.callFunc(c, &struct{}{})
		return
	}

	data := f.getArgs()
	err := json.Unmarshal([]byte(msg.SocketEvent.EmitContent), &data)
	if err != nil {
		log.Println("ERROR DECODING!!")

		if _, ok := data.(string); ok {
			f.callFunc(c, ok)
		} else {
			return
		}
	}

	f.callFunc(c, data)

}

func (m *methods) processPingMessage(c *Channel, msg *protocol.Message) {

	reply := protocol.Message{}
	reply.EngineIoType = protocol.EngineMessageTypePong
	command, _ := protocol.Encode(&reply)
	send(command, c)

}
func (m *methods) processOpenMessage(c *Channel, msg *protocol.Message) {

	reply := protocol.Message{}
	reply.EngineIoType = protocol.EngineMessageTypeMessage
	reply.SocketType = protocol.SocketMessageTypeConnect

	command, _ := protocol.Encode(&reply)
	send(command, c)

	f, _ := m.findMethod(OnConnection)

	f.callFunc(c, &struct{}{})

}
func (m *methods) processDisconnectMessage(c *Channel, msg *protocol.Message) {
	m.onDisconnection(c)

}

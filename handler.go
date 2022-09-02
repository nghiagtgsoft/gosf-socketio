package gosocketio

import (
	"log"
	"sync"

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

func (m *methods) callLoopEvent(c *Channel, event string) {
	if m.onConnection != nil && event == OnConnection {
		log.Println("new connection!", event)
		m.onConnection(c)
	}
	if m.onDisconnection != nil && event == OnDisconnection {
		m.onDisconnection(c)
	}

	f, ok := m.findMethod(event)
	if !ok {
		return
	}

	f.callFunc(c, &struct{}{})
}

/**
Check incoming message
On ack_resp - look for waiter
On ack_req - look for processing function and send ack_resp
On emit - look for processing function
*/
func (m *methods) processIncomingMessage(c *Channel, msg *protocol.Message) {
	log.Println("PROCESS INCOMING MESSAGE")
	switch msg.EngineIoType {
	case protocol.EngineMessageTypeOpen:
		m.processOpenMessage(c, msg)
	}
}

func (m *methods) processOpenMessage(c *Channel, msg *protocol.Message) {
	print("PROCESS OPEN")
	reply := protocol.Message{}
	reply.EngineIoType = protocol.EngineMessageTypeMessage
	reply.SocketType = protocol.SocketMessageTypeConnect

	send(&reply, c, nil)
}

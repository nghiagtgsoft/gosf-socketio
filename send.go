package gosocketio

import (
	"errors"
	"log"

	"github.com/nghiagtgsoft/gosf-socketio/protocol"
)

var (
	ErrorSendTimeout     = errors.New("Timeout")
	ErrorSocketOverflood = errors.New("Socket overflood")
)

/*
*
Send message packet to socket
*/
func send(msg string, c *Channel) error {
	//preventing json/encoding "index out of range" panic
	defer func() {
		if r := recover(); r != nil {
			log.Println("socket.io send panic: ", r)
		}
	}()

	if len(c.out) == queueBufferSize {
		return ErrorSocketOverflood
	}

	c.out <- msg

	return nil
}

/*
*
Create packet based on given data and send it
*/
func (c *Channel) Emit(method string, args interface{}) error {
	// msg := &protocol.Message{
	// 	Type:   protocol.MessageTypeEmit,
	// 	Method: method,
	// }
	msg := protocol.Message{
		EngineIoType: protocol.EngineMessageTypeMessage,
		SocketType:   protocol.SocketMessageTypeEvent,
	}
	//content, _ := json.Marshal(args)
	msg.SocketEvent.EmitName = method
	msg.SocketEvent.EmitContent = args

	command, err := protocol.Encode(&msg)
	if err != nil {
		return err
	}

	send(command, c)
	return nil
}

/**
Create ack packet based on given data and send it and receive response
*/
// func (c *Channel) Ack(method string, args interface{}, timeout time.Duration) (string, error) {
// 	msg := &protocol.Message{
// 		Type:   protocol.MessageTypeAckRequest,
// 		AckId:  c.ack.getNextId(),
// 		Method: method,
// 	}

// 	waiter := make(chan string)
// 	//	c.ack.addWaiter(msg.AckId, waiter)

// 	err := send(msg, c, args)
// 	if err != nil {
// 		//	c.ack.removeWaiter(msg.AckId)
// 	}

// 	select {
// 	case result := <-waiter:
// 		return result, nil
// 		// case <-time.After(timeout):
// 		// 	c.ack.removeWaiter(msg.AckId)
// 		// 	return "", ErrorSendTimeout
// 	}
// }

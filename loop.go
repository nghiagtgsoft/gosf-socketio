package gosocketio

import (
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/ambelovsky/gosf-socketio/color"
	"github.com/ambelovsky/gosf-socketio/protocol"
	"github.com/ambelovsky/gosf-socketio/transport"
)

const (
	queueBufferSize = 10000
)

var (
	ErrorWrongHeader = errors.New("Wrong header")
)

/**
engine.io header to send or receive
*/
type Header struct {
	Sid          string   `json:"sid"`
	Upgrades     []string `json:"upgrades"`
	PingInterval int      `json:"pingInterval"`
	PingTimeout  int      `json:"pingTimeout"`
}

/**
socket.io connection handler

use IsAlive to check that handler is still working
use Dial to connect to websocket
use In and Out channels for message exchange
Close message means channel is closed
ping is automatic
*/
type Channel struct {
	conn transport.Connection

	out    chan string
	header Header

	alive     bool
	aliveLock sync.Mutex

	//ack ackProcessor

	//server  *Server
	ip      string
	request *http.Request
}

/**
create channel, map, and set active
*/
func (c *Channel) initChannel() {
	//TODO: queueBufferSize from constant to server or client variable
	c.out = make(chan string, queueBufferSize)
	//c.ack.resultWaiters = make(map[int](chan string))
	c.setAliveValue(true)
}

/**
Get id of current socket connection
*/
func (c *Channel) Id() string {
	return c.header.Sid
}

/**
Checks that Channel is still alive
*/
func (c *Channel) IsAlive() bool {
	c.aliveLock.Lock()
	isAlive := c.alive
	c.aliveLock.Unlock()

	return isAlive
}

func (c *Channel) setAliveValue(value bool) {
	c.aliveLock.Lock()
	c.alive = value
	c.aliveLock.Unlock()
}

/**
Close channel
*/
func closeChannel(c *Channel, m *methods, args ...interface{}) error {
	if !c.IsAlive() {
		//already closed
		return nil
	}

	c.conn.Close()

	c.setAliveValue(false)

	//clean outloop
	for len(c.out) > 0 {
		<-c.out
	}

	//c.out <- protocol.CloseMessage
	//m.callLoopEvent(c, OnDisconnection)

	deleteOverflooded(c)

	return nil
}

//incoming messages loop, puts incoming messages to In channel
func inLoop(c *Channel, m *methods) error {
	for {
		pkg, err := c.conn.GetMessage()
		if err != nil {
			return closeChannel(c, m, err)
		}
		msg, err := protocol.Decode(pkg)
		if err != nil {
			closeChannel(c, m, protocol.ErrorWrongPacket)
			return err
		}
		go m.processIncomingMessage(c, msg)

	}
}

var overflooded sync.Map

func deleteOverflooded(c *Channel) {
	overflooded.Delete(c)
}

func storeOverflow(c *Channel) {
	overflooded.Store(c, struct{}{})
}

/**
outgoing messages loop, sends messages from channel to socket
*/
func outLoop(c *Channel, m *methods) error {
	for {
		outBufferLen := len(c.out)
		if outBufferLen >= queueBufferSize-1 {
			return closeChannel(c, m, ErrorSocketOverflood)
		} else if outBufferLen > int(queueBufferSize/2) {
			storeOverflow(c)
		} else {
			deleteOverflooded(c)
		}

		msg := <-c.out

		log.Println(color.Purple + "Sending message: " + msg + color.Reset)
		err := c.conn.WriteMessage(msg)
		if err != nil {
			return closeChannel(c, m, err)
		}
	}
}

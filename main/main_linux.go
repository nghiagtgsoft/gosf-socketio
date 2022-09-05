package main

import (
	"log"
	"time"

	gosocketio "github.com/ambelovsky/gosf-socketio"
	"github.com/ambelovsky/gosf-socketio/transport"
)

func main() {
	sockClient, err := gosocketio.Dial(
		gosocketio.GetUrl("localhost", 8001, false),
		transport.GetDefaultWebsocketTransport(),
	)

	if err != nil {
		log.Println(err)
	}

	sockClient.On("connection", func(c *gosocketio.Channel) {
		log.Println("On connection")
		sockClient.Emit("message", "hello")

	})
	sockClient.On("message", func(c *gosocketio.Channel, args string) {
		log.Println("this is the messgage", args)
		time.Sleep(5 * time.Second)
		//	sockClient.Emit("message", "hello you!")
	})
	sockClient.On(gosocketio.OnError, func(c *gosocketio.Channel, e error) {
		log.Println("Error occurs", e)
	})
	sockClient.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel) {
		log.Println("Socket to signalserver disconnected")
	})

	time.Sleep(200 * time.Second)
	log.Println("session ended ", err)

}

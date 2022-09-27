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

	})
	sockClient.On("connectionOpen", func(c *gosocketio.Channel, args string) {
		sockClient.Emit("connectionOpenResponse", `{"userId": "jonas.delrue@jimber.org", "publicKey": "YX1an214Gbk6QNwOpTYeqsdY+3+JopHA6gwtlT2lrl8="}`)
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

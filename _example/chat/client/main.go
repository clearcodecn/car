package main

import (
	"fmt"
	"github.com/clearcodecn/cargo"
	"github.com/gorilla/websocket"
	"log"
)

func main() {
	d := websocket.Dialer{}
	conn, _, err := d.Dial("wss://localhost:9527", nil)
	if err != nil {
		log.Fatal(err)
	}

	client := cargo.NewWebsocketClient(conn)

	client.Write(cargo.NewPacket([]byte("hello world")))

	packet, err := client.NextPacket()

	fmt.Println(packet.String())
}

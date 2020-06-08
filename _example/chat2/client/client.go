package main

import (
	"fmt"
	"github.com/clearcodecn/cargo/gogogo"
	"github.com/clearcodecn/cargo/proto"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

func main() {

	d := websocket.Dialer{}
	conn, _, err := d.Dial("ws://localhost:6666", nil)
	if err != nil {
		log.Fatal(err)
	}

	gogogo.RegisterMessageAndRouter("helloReply", new(proto.HelloRequest), func(msg interface{}, agent *gogogo.Agent) {
		fmt.Println("111")
		fmt.Printf(msg.(*proto.HelloReply).Hahaha)
	})
	ag := gogogo.NewClient(conn)
	ag.WriteMessage("helloRequest", &proto.HelloRequest{
		Wocaonima: "cabi",
	})

	time.Sleep(10 * time.Second)
}

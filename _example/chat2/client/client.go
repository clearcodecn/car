package main

import (
	"fmt"
	"github.com/clearcodecn/cargo"
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

	cargo.RegisterHandler("helloReply", new(proto.HelloReply), func(msg interface{}, agent *cargo.Agent) {
		fmt.Println("111")
		fmt.Printf(msg.(*proto.HelloReply).Hahaha)
	})
	ag := cargo.NewClient(conn)
	ag.WriteMessage("helloRequest", &proto.HelloRequest{
		Wocaonima: "cabi",
	})
	time.Sleep(10 * time.Second)
}

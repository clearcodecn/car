package main

import (
	"fmt"
	"github.com/clearcodecn/car"
	"github.com/clearcodecn/car/proto"
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

	car.RegisterHandler("helloReply", new(proto.HelloReply), func(msg interface{}, agent *car.Agent) {
		fmt.Println("111")
		fmt.Printf(msg.(*proto.HelloReply).Hahaha)
	})
	ag := car.NewClient(conn)
	ag.WriteMessage("helloRequest", &proto.HelloRequest{
		Wocaonima: "cabi",
	})
	time.Sleep(10 * time.Second)
}

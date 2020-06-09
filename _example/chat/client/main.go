package main

import (
	"fmt"
	"github.com/clearcodecn/car/cluster"
	"github.com/clearcodecn/car/packet"
	"github.com/clearcodecn/car/proto"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

var (
	writeChan chan *packet.Packet
)

func main() {
	d := websocket.Dialer{}
	conn, _, err := d.Dial("ws://0.0.0.0:6300/ws", nil)
	if err != nil {
		log.Fatal(err)
	}
	cluster.RegisterHandle(1, func(ctx *cluster.Context, v *proto.Message) error {
		fmt.Println("server1 -> ", string(v.Body))
		return nil
	})
	cluster.RegisterHandle(2, func(ctx *cluster.Context, v *proto.Message) error {
		fmt.Println("server2 -> ", string(v.Body))
		ctx.WriteMessage(&proto.Message{
			Type: 3,
			Body: []byte("要"),
		})
		return nil
	})

	ctx := cluster.NewClientContext(conn)
	ctx.WriteMessage(&proto.Message{
		Type: 1,
		Body: []byte("我是客户端"),
	})
	time.Sleep(1*time.Second)

	ctx.Close()

	time.Sleep(10 * time.Hour)
}

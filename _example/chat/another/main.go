package main

import (
	"fmt"
	"github.com/clearcodecn/cargo"
	"github.com/clearcodecn/cargo/cluster"
	"github.com/clearcodecn/cargo/proto"
	"log"
)

func main() {
	n := cargo.Default()
	cluster.RegisterHandle(proto.MsgType(1), func(ctx *cluster.Context, v *proto.Message) error {
		ctx.WriteMessage(&proto.Message{
			Type: 1,
			Body: []byte("我说鸡蛋你说要"),
		})
		ctx.Stop()
		return nil
	}, func(ctx *cluster.Context, v *proto.Message) error {
		ctx.WriteMessage(&proto.Message{
			Type: 2,
			Body: []byte("鸡蛋"),
		})
		return nil
	})
	cluster.RegisterHandle(proto.MsgType(3), func(ctx *cluster.Context, v *proto.Message) error {
		fmt.Println(string(v.Body))
		v.Type = 2
		v.Body = []byte("鸡蛋")
		ctx.WriteMessage(v)
		return nil
	})
	if err := n.Start(); err != nil {
		log.Fatal(err)
	}
}

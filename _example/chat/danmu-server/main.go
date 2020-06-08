package main

import (
	"github.com/clearcodecn/cargo"
	"github.com/clearcodecn/cargo/cluster"
	"github.com/clearcodecn/cargo/proto"
	"log"
	"time"
)

func main() {
	n := cargo.New(cargo.WithTimeout(time.Hour * 1))
	cluster.RegisterHandle(proto.Auth, func(ctx *cluster.Context, v *proto.Message) error {
		return nil
	})
	cluster.RegisterHandle(proto.Text, func(ctx *cluster.Context, message *proto.Message) error {
		n.SendToAll(message)
		return nil
	})
	if err := n.Start(); err != nil {
		log.Fatal(err)
	}
}

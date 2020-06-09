package main

import (
	"github.com/clearcodecn/car"
	"github.com/clearcodecn/car/cluster"
	"github.com/clearcodecn/car/proto"
	"log"
	"time"
)

func main() {
	n := car.New(car.WithTimeout(time.Hour * 1))
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

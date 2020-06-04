package main

import (
	"github.com/clearcodecn/cargo"
	"github.com/clearcodecn/cargo/cluster"
	"github.com/clearcodecn/cargo/proto"
	"github.com/sirupsen/logrus"
	"log"
	"time"
)

func main() {
	n := cargo.New(cargo.WithTimeout(time.Hour * 1))
	cluster.RegisterHandle(proto.MsgType(1), func(ctx *cluster.Context, v *proto.Message) error {
		logrus.Infof("new msg: %s", v.Body)
		n.SendToAll(v)
		return nil
	})
	if err := n.Start(); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"fmt"
	"github.com/clearcodecn/cargo/gogogo"
	"github.com/clearcodecn/cargo/proto"
)

func main() {

	gogogo.RegisterMessageAndRouter("helloRequest", new(proto.HelloRequest), func(msg interface{}, agent *gogogo.Agent) {
		fmt.Println("new request",msg.(*proto.HelloRequest).Wocaonima)
		agent.WriteMessage("helloReply", &proto.HelloReply{
			Hahaha: "nimabi",
		})
	})

	s := gogogo.Server{}
	s.Start()
}

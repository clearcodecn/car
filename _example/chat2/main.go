package main

import (
	"github.com/clearcodecn/car"
)

type Server struct {
	imServer *car.Server

	handlers []car.Handler
}

func main() {
	var s = new(Server)
}

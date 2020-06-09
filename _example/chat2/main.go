package main

import (
	"github.com/clearcodecn/cargo"
)

type Server struct {
	imServer *cargo.Server

	handlers []cargo.Handler
}

func main() {
	var s = new(Server)
}

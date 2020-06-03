package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

func main() {
	ln, err := net.Listen("tcp", ":1111")
	go func() {
		if err != nil {
			return
		}
		for {
			_, err := ln.Accept()
			if err != nil {
				log.Print(1, err)
				return
			}

		}
	}()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)

	<-ch
	err = ln.Close()
	if err != nil {
		log.Println(err)
	}
	time.Sleep(2 * time.Second)
}

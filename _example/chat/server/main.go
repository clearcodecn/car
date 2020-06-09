package main

import (
	"github.com/clearcodecn/car/old"
	"log"
	"net"
	"net/http"
	"time"
)

func main() {
	server := old.NewServer(old.ServerOption{
		ConnectionIdleTimeout: 30 * time.Second,
		Timeout:               30 * time.Second,
		ReaderBufferSize:      1024,
		WriteBufferSize:       1024,
		EnableCompress:        false,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	})

	ln, err := net.Listen("tcp", ":9527")
	if err != nil {
		log.Fatal(err)
	}

	defer server.Close()

	err = server.Serve(ln)
	if err != nil {
		log.Fatal(err)
	}
}

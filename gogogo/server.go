package gogogo

import (
	"bufio"
	"github.com/clearcodecn/cargo/config"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
	"sync"
)

type Server struct {
	wg sync.WaitGroup

	PushServer PushServer

	opts Option

	upgrader websocket.Upgrader
}

func (s *Server) Start() {
	hp := net.JoinHostPort(config.Ip, config.Port)
	http.ListenAndServe(hp, s)
}

type responseWriterWrapper struct {
	http.ResponseWriter
	isSent bool
}

func (rw *responseWriterWrapper) Write(b []byte) (int, error) {
	rw.isSent = true
	return rw.ResponseWriter.Write(b)
}

func (rw *responseWriterWrapper) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return rw.ResponseWriter.(http.Hijacker).Hijack()
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rw := new(responseWriterWrapper)
	rw.ResponseWriter = w
	w = rw
	if s.opts.BeforeConnect != nil {
		s.opts.BeforeConnect(w, r)
	}
	if rw.isSent {
		return
	}

	conn, err := s.upgrader.Upgrade(rw, r, nil)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	agent := newAgent(conn)
	agent.Run()

	agent.Close()
}
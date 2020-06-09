package car

import (
	"bufio"
	"github.com/clearcodecn/car/config"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	wg       sync.WaitGroup
	opts     Option
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

	agent := newAgent(conn, func(agent *Agent) {
		s.opts.PushServer.DelMember(agent)
		s.opts.PushServer.DelGroupAgent(agent)
	})
	agent.Run()

	agent.Close()
}

func NewServer(opts ...Options) *Server {
	var option Option
	for _, o := range opts {
		o(&option)
	}
	s := new(Server)
	s.opts = option
	s.upgrader = websocket.Upgrader{
		HandshakeTimeout:  30 * time.Second,
		ReadBufferSize:    1024,
		WriteBufferSize:   1024,
		CheckOrigin:       option.CheckOrigin,
		EnableCompression: false,
	}

	return s
}

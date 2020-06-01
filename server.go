package cargo

import (
	"errors"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
)

var (
	increaseID uint32
)

type Server struct {
	mu    sync.Mutex
	conns map[uint32]Transport
	serve bool
	ln    net.Listener

	wg   sync.WaitGroup
	done chan struct{}

	isWebsocket       bool
	upgrader          websocket.Upgrader
	onWebsoketConnect http.HandlerFunc

	isTcp bool

	opt ServerOption
}

func (s *Server) Serve(ln net.Listener) error {
	s.mu.Lock()
	if s.serve {
		s.mu.Unlock()
		return errors.New("call twice serve")
	}
	if s.conns == nil {
		s.conns = make(map[uint32]Transport)
	}
	s.serve = true
	s.ln = ln
	s.wg.Add(1)
	defer func() {
		s.wg.Done()
		<-s.done
	}()
	var err error
	if s.isWebsocket {
		s.upgrader = websocket.Upgrader{
			HandshakeTimeout:  s.opt.Timeout,
			ReadBufferSize:    s.opt.ReaderBufferSize,
			WriteBufferSize:   s.opt.WriteBufferSize,
			CheckOrigin:       s.opt.CheckOrigin,
			EnableCompression: s.opt.EnableCompress,
		}
		err = s.serveWebsocket()
	} else {
		err = s.serveTCP()
	}
	return err
}

func nextID() uint32 {
	return atomic.AddUint32(&increaseID, 1)
}

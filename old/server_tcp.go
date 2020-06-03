package old

import (
	"time"
)

func (s *Server) serveTCP() error {
	var tempDelay time.Duration
	for {
		rawConn, err := s.ln.Accept()
		if err != nil {
			if ne, ok := err.(interface {
				Temporary() bool
			}); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				timer := time.NewTimer(tempDelay)
				select {
				case <-timer.C:
				case <-s.done:
					timer.Stop()
					return nil
				}
				continue
			}
			select {
			case <-s.done:
				return nil
			default:
			}
			return err
		}
		_ = rawConn
		tempDelay = 0
		// Start a new goroutine to deal with rawConn so we don't stall this Accept
		// loop goroutine.
		//
		// Make sure we account for the goroutine so GracefulStop doesn't nil out
		// s.conns before this conn can be added.
		s.wg.Add(1)
		go func() {
			s.wg.Done()
		}()
	}
}

//
//func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	rw := &responseWriter{isSent: false, ResponseWriter: w}
//	if s.onWebsoketConnect != nil {
//		s.onWebsoketConnect(rw, r)
//	}
//	if rw.isSent {
//		return
//	}
//	conn, err := s.upgrader.Upgrade(w, r, nil)
//	if err != nil {
//		return
//	}
//	tr := newWebsocketTransport(conn, s.opt.Timeout)
//	s.mu.Lock()
//	s.conns[nextID()] = tr
//	s.mu.Unlock()
//	config := &ContextConfig{
//		readChannelSize:  1024,
//		writeChannelSize: 1024,
//	}
//	ctx := newContext(tr, config)
//
//	go ctx.writeLoop()
//	ctx.readLoop()
//}

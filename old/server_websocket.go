package old

import (
	"net/http"
)

type responseWriter struct {
	isSent bool
	http.ResponseWriter
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.isSent = true
	return rw.ResponseWriter.Write(b)
}

func (s *Server) serveWebsocket() error {
	erch := make(chan error)
	go func() {
		erch <- http.Serve(s.ln, s)
	}()
	return <-erch
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rw := &responseWriter{isSent: false, ResponseWriter: w}
	if s.onWebsoketConnect != nil {
		s.onWebsoketConnect(rw, r)
	}
	if rw.isSent {
		return
	}
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	tr := newWebsocketTransport(conn, s.opt.Timeout)
	s.mu.Lock()
	s.conns[nextID()] = tr
	s.mu.Unlock()
	config := &ContextConfig{
		readChannelSize:  1024,
		writeChannelSize: 1024,
	}
	ctx := newContext(tr, config)

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		ctx.writeLoop()
	}()
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		ctx.listenLoop()
	}()
	ctx.readLoop()
}

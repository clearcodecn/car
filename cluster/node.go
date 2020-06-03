package cluster

import (
	"crypto/tls"
	"errors"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
	"strings"
	"sync"
)

type Node struct {
	Id string

	opt ServerOptions

	BeforeWebsocketConnect func(w http.ResponseWriter, r *http.Request)

	wg sync.WaitGroup

	Handlers []HandlerFunc
}

type responseWriter struct {
	isSent bool
	http.ResponseWriter
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.isSent = true
	return rw.ResponseWriter.Write(b)
}

func (n *Node) startWebsocketServer() error {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:    n.opt.ReaderBufferSize,
		WriteBufferSize:   n.opt.WriteBufferSize,
		CheckOrigin:       n.opt.CheckOrigin,
		EnableCompression: n.opt.EnableCompress,
		HandshakeTimeout:  n.opt.Timeout,
	}

	path := "/" + strings.TrimPrefix(n.opt.WebsocketPath, "/")

	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		w = &responseWriter{ResponseWriter: w}
		if n.BeforeWebsocketConnect != nil {
			n.BeforeWebsocketConnect(w, r)
		}
		if w.(*responseWriter).isSent {
			return
		}
		wc, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		netConn := newWsConn(wc, n.opt.Timeout)
		n.wg.Add(1)
		defer n.wg.Done()

		n.handleConnection(netConn)
	})
	var err error
	var hostPort = net.JoinHostPort(n.opt.Host, n.opt.Port)

	if n.opt.Cert != "" && n.opt.Key != "" {
		err = http.ListenAndServeTLS(hostPort, n.opt.Cert, n.opt.Key, nil)
	} else {
		err = http.ListenAndServe(hostPort, nil)
	}
	if err != nil {
		if errors.Is(http.ErrServerClosed, err) {
			return nil
		}
	}
	return err
}

func (n *Node) startTCPServer() error {
	var (
		hostPort = net.JoinHostPort(n.opt.Host, n.opt.Port)
		ln       net.Listener
		err      error
	)
	if n.opt.Cert != "" && n.opt.Key != "" {
		cert, er := tls.LoadX509KeyPair(n.opt.Cert, n.opt.Key)
		if er != nil {
			return er
		}
		cfg := &tls.Config{Certificates: []tls.Certificate{cert}}
		ln, err = tls.Listen("tcp", hostPort, cfg)
	} else {
		ln, err = net.Listen("tcp", hostPort)
	}
	if err != nil {
		return err
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// server closed.
			if strings.Contains(err.Error(), "use of closed network connection") {
				return nil
			}
			return err
		}
		n.wg.Add(1)
		go func() {
			defer n.wg.Done()
			n.handleConnection(conn)
		}()
	}
}

func (n *Node) handleConnection(conn net.Conn) {
	config := &ContextConfig{
		readBufferSize:    n.opt.ReaderBufferSize,
		writeBufferSize:   n.opt.WriteBufferSize,
		readChannelSize:   n.opt.ReaderBufferSize,
		writeChannelSize:  n.opt.WriteBufferSize,
		AfterReadHandler:  []OnIOHandler{AfterRead},
		AfterWriteHandler: []OnIOHandler{AfterWrite},
		Handlers:          n.Handlers,
	}
	ctx := newContext(conn, config)
	go ctx.loop()
	ctx.readLoop()
}

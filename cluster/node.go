package cluster

import (
	"bufio"
	"context"
	"crypto/tls"
	"errors"
	"github.com/clearcodecn/cargo/proto"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Node struct {
	Id string

	opt ServerOptions

	wg sync.WaitGroup

	state uint32
	done  chan struct{}

	httpServer *http.Server

	mu          sync.Mutex
	connections map[uint32]*Context

	groupManager *contextSet
	ipManager    *contextSet
}

type responseWriter struct {
	isSent bool
	http.ResponseWriter
}

func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return rw.ResponseWriter.(http.Hijacker).Hijack()
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
	var err error
	var hostPort = net.JoinHostPort(n.opt.Host, n.opt.Port)
	var httpServer = &http.Server{Addr: hostPort}
	var path = "/" + strings.TrimPrefix(n.opt.WebsocketPath, "/")
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		w = &responseWriter{ResponseWriter: w}
		if n.opt.BeforeWebsocketConnect != nil {
			n.opt.BeforeWebsocketConnect(w, r)
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
	n.httpServer = httpServer

	if n.opt.Cert != "" && n.opt.Key != "" {
		nodeLogger.Debugf("websocket start at: wss://%s%s", hostPort, path)
		err = httpServer.ListenAndServeTLS(n.opt.Cert, n.opt.Key)
	} else {
		nodeLogger.Debugf("websocket start at: ws://%s%s", hostPort, path)
		err = httpServer.ListenAndServe()
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
		writeBufferSize:   n.opt.WriteChannelSize,
		readChannelSize:   n.opt.ReadChannelSize,
		writeChannelSize:  n.opt.WriteBufferSize,
		AfterReadHandler:  []OnIOHandler{AfterRead},
		AfterWriteHandler: []OnIOHandler{AfterWrite},
		codec:             n.opt.Codec,
	}
	ctx := newContext(conn, n, config)
	defer ctx.Close()

	go ctx.loop()
	ctx.readLoop()
}

func (n *Node) HttpServer() *http.Server {
	return n.httpServer
}

func NewNode(opt ServerOptions) *Node {
	node := new(Node)
	node.opt = opt
	logrus.SetLevel(opt.LogLevel)
	node.connections = make(map[uint32]*Context)
	node.groupManager = newContextSet()
	node.ipManager = newContextSet()
	node.done = make(chan struct{})

	return node
}

func (n *Node) Start() error {
	n.wg.Add(1)
	go func() {
		defer n.wg.Done()
		n.cleanLoop()
	}()

	var err error
	if n.opt.IsWebsocket {
		err = n.startWebsocketServer()
	} else {
		err = n.startTCPServer()
	}

	return err
}

func (n *Node) ShutDown() error {
	if !atomic.CompareAndSwapUint32(&n.state, 0, 1) {
		return errors.New("node already closed")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	n.httpServer.Shutdown(ctx)
	close(n.done)

	n.mu.Lock()
	for _, conn := range n.connections {
		conn.Close()
	}
	n.mu.Unlock()
	n.wg.Wait()

	return nil
}

func (n *Node) SendToAll(message *proto.Message) {
	n.mu.Lock()
	ctxs := n.connections
	n.mu.Unlock()
	for _, c := range ctxs {
		if err := c.WriteMessage(message); err != nil {
			nodeLogger.WithField("method", "SendToAll").WithError(err).Warnf("write message failed")
		}
	}
}

func (n *Node) SendByID(id string, message *proto.Message) error {
	var ctx *Context
	n.mu.Lock()
	for _, c := range n.connections {
		if c.id == id {
			ctx = c
		}
	}
	n.mu.Unlock()
	if ctx == nil {
		return errors.New("id not found")
	}
	return ctx.WriteMessage(message)
}

func (n *Node) AddToGroup(group string, ctx *Context) {
	n.groupManager.add(group, ctx)
}

func (n *Node) cleanLoop() {
	tk := time.NewTicker(n.opt.CleanDuration)
	defer tk.Stop()

	for {
		select {
		case <-n.done:
		case <-tk.C:
			n.ipManager.cleanClosed()
			n.groupManager.cleanClosed()
		}
	}
}

package cluster

import (
	"github.com/gorilla/websocket"
	"net"
	"time"
)

type wsConn struct {
	net.Conn
	wsConn *websocket.Conn

	timeout time.Duration
}

func (c *wsConn) Read(b []byte) (n int, err error) {
	if c.timeout != 0 {
		c.wsConn.SetReadDeadline(time.Now().Add(c.timeout))
	}
	_, r, err := c.wsConn.NextReader()
	if err != nil {
		return 0, err
	}
	n, err = r.Read(b)
	c.wsConn.SetReadDeadline(time.Time{})
	return n, err
}

func (c *wsConn) Write(b []byte) (n int, err error) {
	if c.timeout != 0 {
		c.wsConn.SetWriteDeadline(time.Now().Add(c.timeout))
	}
	w, err := c.wsConn.NextWriter(websocket.BinaryMessage)
	if err != nil {
		return 0, err
	}
	defer w.Close()
	n, err = w.Write(b)
	c.wsConn.SetWriteDeadline(time.Time{})
	return n, err
}

func newWsConn(c *websocket.Conn, timeout time.Duration) net.Conn {
	return &wsConn{wsConn: c, timeout: timeout, Conn: c.UnderlyingConn()}
}

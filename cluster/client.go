package cluster

import (
	"github.com/clearcodecn/cargo/codec"
	"github.com/gorilla/websocket"
	"time"
)

func NewClientContext(conn *websocket.Conn) *Context {
	cc := newWsConn(conn, 30*time.Second)
	ctx := newContext(cc, nil, &ContextConfig{
		readBufferSize:    1024,
		writeBufferSize:   1024,
		readChannelSize:   1024,
		writeChannelSize:  1024,
		AfterReadHandler:  nil,
		AfterWriteHandler: nil,
		codec:             &codec.JsonCodec{},
	})

	go ctx.loop()
	go ctx.readLoop()

	return ctx
}

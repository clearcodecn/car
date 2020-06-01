package cargo

import "net"

type Transport interface {
	net.Conn
}

type MsgHandler func(*Context, *Packet) error

type Packet struct{}

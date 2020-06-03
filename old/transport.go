package old

import (
	"bytes"
	"net"
)

type SendMode int

const (
	SendModeBlock SendMode = iota + 1
	SendModeAsync
)

func (sm SendMode) isBlock() bool {
	return sm == SendModeBlock
}

type Transport interface {
	net.Conn
}

type MsgHandler func(*Context, *Packet)

type Packet struct {
	body            *bytes.Buffer
	Event           Event
	isClusterPacket bool
}

func NewPacket(data []byte) *Packet {
	p := new(Packet)
	p.body = bytes.NewBuffer(data)
	return p
}

func (p *Packet) String() string {
	return p.body.String()
}

type Event uint32
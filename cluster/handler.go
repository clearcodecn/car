package cluster

import (
	"errors"
	"fmt"
	"github.com/clearcodecn/cargo/packet"
	"github.com/clearcodecn/cargo/proto"
)

var (
	ErrInvalidPayload = errors.New("invalid payload")
	handlers          = make(map[proto.MsgType][]HandlerFunc)
)

type MiddlewareFunc func(ctx *Context, p *packet.Packet) error

type HandlerFunc func(ctx *Context, message *proto.Message) error

func RecoveryHandler(ctx *Context, p *packet.Packet) MiddlewareFunc {
	return func(ctx *Context, p *packet.Packet) error {
		var err error
		defer func() {
			if er := recover(); er != nil {
				err = fmt.Errorf("%v", er)
			}
		}()
		return err
	}
}

func (c *Context) handleMessage(payload []byte) error {
	var msg = new(proto.Message)
	err := c.codec.UnMarshal(payload, msg)
	if err != nil {
		return ErrInvalidPayload
	}
	var hs []HandlerFunc
	var ok bool
	nhs, ok := handlers[msg.Type]
	if ok {
		hs = append(hs, nhs...)
	}
	if len(hs) == 0 {
		return errors.New("not found handler")
	}

	for _, h := range hs {
		if err := h(c, msg); err != nil {
			return err
		}
		if c.stop {
			break
		}
	}

	c.stop = false
	return nil
}

func RegisterHandle(t proto.MsgType, handlerFunc ...HandlerFunc) {
	if _, ok := handlers[t]; !ok {
		handlers[t] = make([]HandlerFunc, 0)
	}
	handlers[t] = append(handlers[t], handlerFunc...)
}

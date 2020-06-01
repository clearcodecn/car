package cargo

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"sync"
	"sync/atomic"
	"time"
)

type ContextConfig struct {
	readChannelSize  int
	writeChannelSize int
	beforeWrite      func(ctx *Context, packet *Packet) bool
	afterWrite       func(ctx *Context, packet *Packet, n int64)
	afterRead        func(ctx *Context, packet *Packet)
}

type Context struct {
	transport Transport

	createTime int64

	beforeWrite func(ctx *Context, packet *Packet) bool

	afterWrite func(ctx *Context, packet *Packet, n int)

	afterRead func(ctx *Context, packet *Packet, n int)

	readChannel chan *Packet

	writeChannel chan *Packet

	codec Codec

	done chan struct{}

	state uint32

	readBufferSize int

	writeBufferSize int

	decodeError int
}

func newContext(transport Transport, config *ContextConfig) *Context {
	ctx := new(Context)
	ctx.transport = transport
	ctx.createTime = time.Now().Unix()
	ctx.readChannel = make(chan *Packet, config.readChannelSize)
	ctx.writeChannel = make(chan *Packet, config.writeChannelSize)

	return ctx
}

func (c *Context) sendPacket(packet *Packet) error {
	select {
	case <-c.done:
		return errors.New("channel already been closed")
	case c.writeChannel <- packet:
	}
	return nil
}

func (c *Context) writeLoop() {
	var w = bufio.NewWriter(c.transport)
	if c.writeBufferSize != 0 {
		w = bufio.NewWriterSize(c.transport, c.writeBufferSize)
	}
	defer func() {
		if !c.isClosed() {
			c.closeChannel()
		}
	}()

	for {
		select {
		case <-c.done:
			return
		case packet, ok := <-c.writeChannel:
			if !ok {
				return
			}
			if c.beforeWrite != nil && !c.beforeWrite(c, packet) {
				continue
			}
			data, err := c.codec.Marshal(packet)
			if err != nil {
				// TODO:: log.
				return
			}
			n := len(data)
			var header []byte
			binary.BigEndian.PutUint16(header, uint16(n))
			n, err = w.Write(header)
			if err != nil {
				return
			}
			n, err = w.Write(data)
			if err != nil {
				return
			}
			err = w.Flush()
			if err != nil {
				return
			}
			if c.afterWrite != nil {
				c.afterWrite(c, packet, n+2)
			}
		}
	}
}

// readLoop
// first 2 + body's length
func (c *Context) readLoop() {
	var r io.Reader = c.transport
	if c.readBufferSize != 0 {
		r = bufio.NewReaderSize(c.transport, c.readBufferSize*2)
	}
	defer func() {
		if !c.isClosed() {
			c.closeChannel()
		}
	}()
	for {
		select {
		case <-c.done:
			return
		default:
		}
		var b = make([]byte, 2)
		n, err := io.ReadFull(r, b)
		if err != nil {
			if !c.isClosed() {
				c.closeChannel()
			}
			return
		}
		length := binary.BigEndian.Uint16(b)
		b = make([]byte, length)
		n, err = io.ReadFull(r, b)
		if err != nil {
			if !c.isClosed() {
				c.closeChannel()
			}
			return
		}
		var packet = new(Packet)
		if err = c.codec.UnMarshal(b[:n], packet, c); err != nil {
			c.decodeError++
			if c.decodeError >= 10 {
				c.closeChannel()
				return
			}
			continue
		}
		select {
		case <-c.done:
			return
		case c.readChannel <- packet:
		}
		if c.afterRead != nil {
			c.afterRead(c, packet, n+2)
		}
	}
}

func (c *Context) closeChannel() error {
	if atomic.CompareAndSwapUint32(&c.state, 1, 1) {
		return errors.New("channel already been closed")
	}
	atomic.StoreUint32(&c.state, 1)

	c.transport.Close()

	close(c.done)
	close(c.writeChannel)
	close(c.readChannel)
	return nil
}

func (c *Context) isClosed() bool {
	return atomic.CompareAndSwapUint32(&c.state, 1, 1)
}

type contextSet struct {
	collection map[string][]*Context
	mu         sync.Mutex
}

func (c *contextSet) add(key string, ctx *Context) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.collection[key]; !ok {
		c.collection[key] = make([]*Context, 0)
	}
	c.collection[key] = append(c.collection[key], ctx)
}

func (c *contextSet) RangeByKey(key string, f func(ctx *Context) bool) bool {
	var coll []*Context
	var ok bool
	c.mu.Lock()
	coll, ok = c.collection[key]
	c.mu.Unlock()

	if !ok {
		return false
	}

	for _, ctx := range coll {
		if !f(ctx) {
			break
		}
	}
	return true
}

func (c *contextSet) remove(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.collection, key)
}

func (c *contextSet) get(key string) ([]*Context, bool) {
	var coll []*Context
	var ok bool
	c.mu.Lock()
	coll, ok = c.collection[key]
	c.mu.Unlock()
	return coll, ok
}

func (c *contextSet) removeContext(key string, ctx *Context) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if coll, ok := c.collection[key]; ok {
		var index int = -1
		for i, cc := range coll {
			if cc == ctx {
				index = i
				break
			}
		}
		if index < 0 {
			return false
		}

		coll = append(coll[:index], coll[index+1:]...)
		c.collection[key] = coll
		return true
	}

	return false
}

func (c *contextSet) clean() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.collection = make(map[string][]*Context)
}

func (c *contextSet) all() []*Context {
	var ctxs []*Context
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, c := range c.collection {
		ctxs = append(ctxs, c...)
	}
	return ctxs
}

func newContextSet() *contextSet {
	cs := new(contextSet)
	cs.collection = make(map[string][]*Context)
	return cs
}

type ContextManager interface {
	Bind(token string, ctx *Context) error
	UnBind(token string, ctx *Context) error
	SendByToken(ctx *Context, token string, packet *Packet) error
	BSendByToken(ctx *Context, token string, packet *Packet) error

	SendToAll(packet *Packet) error
	BSendToAll(packet *Packet) error
}

package cluster

import (
	"bufio"
	"encoding/binary"
	"errors"
	"github.com/clearcodecn/cargo/codec"
	"github.com/clearcodecn/cargo/packet"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type OnIOHandler func(ctx *Context, packet *packet.Packet, n int)

type ContextConfig struct {
	readBufferSize    int
	writeBufferSize   int
	readChannelSize   int
	writeChannelSize  int
	AfterReadHandler  []OnIOHandler
	AfterWriteHandler []OnIOHandler
	codec             codec.Codec
	Handlers          []HandlerFunc
}

type Context struct {
	id   string
	ip   string
	port string

	mu            sync.Mutex
	createTime    int64
	lastReadTime  int64
	lastWriteTime int64
	readCount     uint32
	writeCount    uint32

	value map[interface{}]interface{}

	conn net.Conn
	br   *bufio.Reader
	bw   *bufio.Writer

	writeChan chan *packet.Packet
	readChan  chan *packet.Packet

	done chan struct{}

	AfterReadHandler  []OnIOHandler
	AfterWriteHandler []OnIOHandler
	Handlers          []HandlerFunc

	state uint32

	codec codec.Codec
}

type HandlerFunc func(ctx *Context, packet *packet.Packet)

func newContext(conn net.Conn, config *ContextConfig) *Context {
	ctx := new(Context)
	ctx.readChan = make(chan *packet.Packet, config.readChannelSize)
	ctx.writeChan = make(chan *packet.Packet, config.writeBufferSize)
	ctx.done = make(chan struct{})
	ctx.br = bufio.NewReaderSize(conn, config.readBufferSize)
	ctx.bw = bufio.NewWriterSize(conn, config.writeBufferSize)
	ctx.conn = conn
	ctx.value = make(map[interface{}]interface{})
	ctx.codec = config.codec

	host, port, _ := net.SplitHostPort(conn.RemoteAddr().String())
	ctx.ip = host
	ctx.port = port

	ctx.lastReadTime = time.Now().Unix()
	ctx.lastReadTime = time.Now().Unix()
	ctx.lastWriteTime = time.Now().Unix()
	ctx.AfterReadHandler = config.AfterReadHandler
	ctx.AfterWriteHandler = config.AfterWriteHandler

	return ctx
}

func (ctx *Context) readLoop() {
	for {
		select {
		case <-ctx.done:
			return
		default:
		}
		var header = make([]byte, 2)
		var count int
		n, err := io.ReadFull(ctx.br, header)
		if err != nil {
			return
		}
		count += n
		l := binary.BigEndian.Uint16(header)
		var command = make([]byte, 2)
		n, err = io.ReadFull(ctx.br, command)
		if err != nil {
			return
		}
		cmdType := binary.BigEndian.Uint16(command)
		count += n
		var body = make([]byte, l-cmdType)
		n, err = io.ReadFull(ctx.br, body)
		if err != nil {
			return
		}
		count += n

		p, err := packet.NewPacket(packet.Command(cmdType), body)
		if err != nil {
			// TODO:: invalid packet, log it.
			continue
		}

		select {
		case <-ctx.done:
			return
		case ctx.writeChan <- p:
		}
		if len(ctx.AfterReadHandler) != 0 {
			for _, h := range ctx.AfterReadHandler {
				h(ctx, p, count)
			}
		}
	}
}

// protocol
// | header (2) | command(2) | body(unknown) |
// header's value = command length + body length.
func (ctx *Context) loop() {
	for {
		select {
		case <-ctx.done:
			return
		case p, ok := <-ctx.readChan:
			if !ok {
				return
			}
			switch p.Command {
			case packet.CommandHeartBeat:
				ctx.writeChan <- packet.PacketPong
			case packet.CommandKick:
				ctx.Close()
				return
			case packet.CommandClose:
				ctx.Close()
				return
			case packet.CommandData:
				for _, h := range ctx.Handlers {
					h(ctx, p)
				}
			}
		case p, ok := <-ctx.writeChan:
			if !ok {
				return
			}
			cmd := p.Command
			var cmdHeader = make([]byte, 2)
			binary.BigEndian.PutUint16(cmdHeader, uint16(cmd))
			var header = make([]byte, 2)
			binary.BigEndian.PutUint16(header, uint16(2+len(p.Payload)))
			var count int
			n, err := ctx.bw.Write(header)
			if err != nil {
				return
			}
			count += n
			n, err = ctx.bw.Write(cmdHeader)
			if err != nil {
				return
			}
			count += n
			n, err = ctx.bw.Write(p.Payload)
			if err != nil {
				return
			}
			count += n
			if err := ctx.bw.Flush(); err != nil {
				return
			}
			if len(ctx.AfterWriteHandler) != 0 {
				for _, h := range ctx.AfterWriteHandler {
					h(ctx, p, count)
				}
			}
		}
	}
}

// BindID goroutine unsafe
func (ctx *Context) BindID(id string) {
	ctx.id = id
}

func (ctx *Context) Ip() string {
	return ctx.ip
}

func (ctx *Context) ClientAddr() string {
	return ctx.ip + ":" + ctx.port
}

// Reset goroutine unsafe
func (ctx *Context) Reset() {
	ctx.id = ""
	ctx.br = nil
	ctx.bw = nil
	ctx.readChan = nil
	ctx.writeChan = nil
	ctx.done = nil
	ctx.value = nil
	return
}

func (ctx *Context) ReadCount() uint32 {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	return ctx.readCount
}

func (ctx *Context) WriteCount() uint32 {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	return ctx.writeCount
}

func (ctx *Context) IsClose() bool {
	return atomic.CompareAndSwapUint32(&ctx.state, 1, 1)
}

func (ctx *Context) Close() error {
	if !atomic.CompareAndSwapUint32(&ctx.state, 0, 1) {
		return errors.New("connection already closed")
	}
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	close(ctx.done)
	close(ctx.readChan)
	close(ctx.writeChan)
	ctx.conn.Close()
	ctx.value = nil
	return nil
}

func AfterRead(ctx *Context, packet2 *packet.Packet, n int) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	ctx.readCount += uint32(n)
	ctx.lastReadTime = time.Now().Unix()
}

func AfterWrite(ctx *Context, packet *packet.Packet, n int) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	ctx.writeCount += uint32(n)
	ctx.lastWriteTime = time.Now().Unix()
}

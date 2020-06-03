package old

import (
	"bufio"
	"encoding/binary"
	"errors"
	"github.com/gorilla/websocket"
	"io"
	"sync"
	"sync/atomic"
	"time"
)

type WebSocketClient struct {
	Codec        Codec
	Transport    Transport
	ReadChannel  chan *Packet
	WriteChannel chan *Packet
	done         chan struct{}
	state        uint32
	wg           sync.WaitGroup
}

func (wsc *WebSocketClient) writeLoop() {
	defer func() {
		wsc.Close()
	}()

	w := bufio.NewWriterSize(wsc.Transport, 1024)
	for {
		select {
		case <-wsc.done:
			return
		case packet, ok := <-wsc.WriteChannel:
			if !ok {
				return
			}
			data, err := wsc.Codec.Marshal(packet)
			if err != nil {
				continue
			}
			l := len(data)
			var header []byte
			binary.BigEndian.PutUint16(header, uint16(l))
			_, err = w.Write(header)
			if err != nil {
				return
			}
			_, err = w.Write(data)
			if err != nil {
				return
			}

			err = w.Flush()
			if err != nil {
				return
			}

		default:
		}
	}
}

func (wsc *WebSocketClient) Close() error {
	if atomic.CompareAndSwapUint32(&wsc.state, 0, 1) {
		atomic.StoreUint32(&wsc.state, 1)
		close(wsc.done)
		close(wsc.WriteChannel)

		wsc.wg.Wait()
		return nil
	}

	return errors.New("client already closed")
}

func (wsc *WebSocketClient) Write(packet *Packet) error {
	select {
	case <-wsc.done:
		return errors.New("client closed")
	default:
	}
	wsc.WriteChannel <- packet
	return nil
}

func NewWebsocketClient(conn *websocket.Conn) *WebSocketClient {
	tr := newWebsocketTransport(conn, 30*time.Second)
	wsc := &WebSocketClient{
		Codec:        &defaultCodec{},
		Transport:    tr,
		ReadChannel:  make(chan *Packet, 1024),
		WriteChannel: make(chan *Packet, 1024),
		done:         make(chan struct{}),
		state:        0,
	}

	wsc.wg.Add(1)
	go func() {
		defer wsc.wg.Done()
		wsc.writeLoop()
	}()

	wsc.wg.Add(1)
	go func() {
		defer wsc.wg.Done()
		wsc.readLoop()
	}()

	return wsc
}

func (wsc *WebSocketClient) readLoop() {
	for {
		select {
		case <-wsc.done:
			return
		default:
		}
		var data = make([]byte, 2)
		_, err := io.ReadFull(wsc.Transport, data)
		if err != nil {
			return
		}
		l := binary.BigEndian.Uint16(data)
		data = make([]byte, l)
		_, err = io.ReadFull(wsc.Transport, data)
		if err != nil {
			return
		}
		packet := NewPacket(data)
		select {
		case <-wsc.done:
			return
		case wsc.ReadChannel <- packet:
		}
	}
}

func (wsc *WebSocketClient) NextPacket() (*Packet, error) {
	select {
	case packet, ok := <-wsc.ReadChannel:
		if !ok {
			return nil, io.EOF
		}
		return packet, nil
	case <-wsc.done:
		return nil, io.EOF
	}
}

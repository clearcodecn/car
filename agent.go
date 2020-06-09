package car

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

var (
	ErrConnectionAlreadyClosed = errors.New("connection already closed")
)

type Agent struct {
	writeChan chan []byte
	readChan  chan []byte
	conn      *websocket.Conn

	mu         sync.Mutex
	id         string
	isClose    bool
	writeCount int
	readCount  int

	MessageProcess MessageProcess
	Codec          MessageCodec

	onClose []func(agent *Agent)
}

func newAgent(conn *websocket.Conn, onClose ...func(agent *Agent)) *Agent {
	ag := new(Agent)
	ag.writeChan = make(chan []byte, 16)
	ag.readChan = make(chan []byte, 1)
	ag.isClose = false
	ag.MessageProcess = defaultMessageProcess
	ag.Codec = defaultCodec
	ag.conn = conn

	go func() {
		for b := range ag.writeChan {
			w, err := conn.NextWriter(websocket.BinaryMessage)
			if err != nil {
				break
			}
			n, err := w.Write(b)
			if err != nil {
				break
			}
			err = w.Close()
			if err != nil {
				break
			}
			ag.mu.Lock()
			ag.writeCount += n
			ag.mu.Unlock()
		}
		ag.Close()
	}()

	return ag
}

func (ag *Agent) Run() {
	defer ag.Close()

	for {
		_, data, err := ag.conn.ReadMessage()
		if err != nil {
			return
		}
		// |     4       | 1      | 1        |  method |  body |
		// | body length | version | command |  dynamic|  dyncmic |
		var (
			methodLength = data[:2]
			bodyLength   = data[2:4]
		)
		ml := binary.BigEndian.Uint16(methodLength)
		bl := binary.BigEndian.Uint16(bodyLength)
		method := data[6 : 6+ml]
		body := data[6+ml : 6+ml+bl]

		ag.mu.Lock()
		ag.readCount += len(data)
		ag.mu.Unlock()
		msg := ag.MessageProcess.Clone(string(method))
		if msg == nil {
			continue
		}
		go func() {
			if err := recover(); err != nil {
				log.Println(err)
			}
			if err = ag.Codec.UnMarshal(body, msg); err != nil {
				return
			}
			if !ag.MessageProcess.Route(string(method), msg, ag) {
				return
			}
		}()
	}
}

func (ag *Agent) WriteMessage(event string, msg interface{}) error {
	data, err := ag.Codec.Marshal(msg)
	if err != nil {
		return err
	}
	var header = make([]byte, 6)
	var w = bytes.NewBuffer(nil)
	binary.BigEndian.PutUint16(header[:2], uint16(len(event)))
	binary.BigEndian.PutUint16(header[2:4], uint16(len(data)))
	header[4] = 0
	header[5] = 0

	w.Write(header)
	w.WriteString(event)
	w.Write(data)

	ag.mu.Lock()
	defer ag.mu.Unlock()

	if ag.isClose {
		return ErrConnectionAlreadyClosed
	}
	ag.writeChan <- w.Bytes()
	return nil
}

func (ag *Agent) Close() error {
	ag.mu.Lock()
	defer ag.mu.Unlock()

	if ag.isClose {
		return ErrConnectionAlreadyClosed
	}
	ag.conn.Close()
	ag.isClose = true

	close(ag.writeChan)
	close(ag.readChan)

	for _, f := range ag.onClose {
		f(ag)
	}

	return nil
}

func (ag *Agent) BindID(id string) error {
	ag.mu.Lock()
	defer ag.mu.Unlock()
	if ag.isClose {
		return ErrConnectionAlreadyClosed
	}
	ag.id = id
	return nil
}

func (ag *Agent) ID() string {
	ag.mu.Lock()
	defer ag.mu.Unlock()
	return ag.id
}

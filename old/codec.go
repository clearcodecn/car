package old

import "bytes"

type Codec interface {
	Marshal(packet *Packet) ([]byte, error)
	UnMarshal(b []byte, packet *Packet, ctx *Context) error
}

type defaultCodec struct{}

func (j defaultCodec) Marshal(packet *Packet) ([]byte, error) {
	return append([]byte{uint8(packet.Event)}, packet.body.Bytes()...), nil
}

func (j defaultCodec) UnMarshal(b []byte, packet *Packet, ctx *Context) error {
	if len(b) > 0 {
		packet.Event = Event(b[0])
	}
	packet.body = bytes.NewBuffer(b[1:])
	return nil
}

package cargo

import "bytes"

type Codec interface {
	Marshal(packet *Packet) ([]byte, error)
	UnMarshal(b []byte, packet *Packet, ctx *Context) error
}

type defaultCodec struct{}

func (j defaultCodec) Marshal(packet *Packet) ([]byte, error) {
	return packet.body.Bytes(), nil
}

func (j defaultCodec) UnMarshal(b []byte, packet *Packet, ctx *Context) error {
	packet.body = bytes.NewBuffer(b)
	return nil
}
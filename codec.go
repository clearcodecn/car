package cargo

type Codec interface {
	Marshal(packet *Packet) ([]byte, error)
	UnMarshal(b []byte, packet *Packet, ctx *Context) error
}

package codec

import (
	"github.com/clearcodecn/cargo/packet"
)

type Codec interface {
	Marshal(packet *packet.Packet) ([]byte, error)
	UnMarshal([]byte, *packet.Packet) error
}

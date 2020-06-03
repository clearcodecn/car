package packet

import "errors"

type Command int

const (
	CommandHeartBeat = iota + 1
	CommandData
	CommandClose
	CommandKick
)

type Packet struct {
	Command Command
	Payload []byte
}

var (
	PacketPing = &Packet{Command: CommandHeartBeat, Payload: []byte("ping")}
	PacketPong = &Packet{Command: CommandHeartBeat, Payload: []byte("pong")}
)

func NewPacket(cmd Command, payload []byte) (*Packet, error) {
	if !validCommand(cmd) {
		return nil, errors.New("invalid command")
	}
	p := &Packet{}
	p.Command = cmd
	p.Payload = payload
	return p, nil
}

func validCommand(cmd Command) bool {
	return cmd >= CommandHeartBeat && cmd <= CommandKick
}

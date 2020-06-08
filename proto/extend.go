package proto

const (
	Auth MsgType = iota + 1
	Text
)

type HelloRequest struct {
	Wocaonima string
}

type HelloReply struct {
	Hahaha string
}

func (*HelloReply) Clone() interface{} {
	return new(HelloReply)
}

func (HelloRequest) Clone() interface{} {
	return new(HelloRequest)
}

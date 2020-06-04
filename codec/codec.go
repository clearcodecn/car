package codec

import (
	"encoding/json"
	"errors"
	"github.com/golang/protobuf/proto"
)

type Codec interface {
	Marshal(interface{}) ([]byte, error)
	UnMarshal([]byte, interface{}) error
}

type JsonCodec struct{}

func (j JsonCodec) Marshal(i interface{}) ([]byte, error) {
	return json.Marshal(i)
}

func (j JsonCodec) UnMarshal(bytes []byte, i interface{}) error {
	return json.Unmarshal(bytes, i)
}

type ptotoCodec struct{}

func (j ptotoCodec) Marshal(i interface{}) ([]byte, error) {
	m, ok := i.(proto.Message)
	if !ok {
		return nil, errors.New("invalid message")
	}
	return proto.Marshal(m)
}

func (j ptotoCodec) UnMarshal(bytes []byte, i interface{}) error {
	m, ok := i.(proto.Message)
	if !ok {
		return errors.New("invalid message")
	}

	return proto.Unmarshal(bytes, m)
}

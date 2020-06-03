package codec

import "encoding/json"

type Codec interface {
	Marshal(interface{}) ([]byte, error)
	UnMarshal([]byte, interface{}) error
}

type jsonCodec struct{}

func (j jsonCodec) Marshal(i interface{}) ([]byte, error) {
	return json.Marshal(i)
}

func (j jsonCodec) UnMarshal(bytes []byte, i interface{}) error {
	return json.Unmarshal(bytes, i)
}
package car

import "encoding/json"

type MessageHandler func(msg interface{}, agent *Agent)

type MessageCloneable interface {
	Clone() interface{}
}

type MessageProcess interface {
	RegisterMessageAndRouter(method string, v MessageCloneable, handler MessageHandler)
	Clone(method string) interface{}
	Route(method string, msg interface{}, agent *Agent) bool
}

type MessageCodec interface {
	Marshal(v interface{}) ([]byte, error)
	UnMarshal([]byte, interface{}) error
}

var (
	defaultMessageProcess MessageProcess
	defaultCodec          MessageCodec
)

func init() {
	defaultMessageProcess = newMessageProcess()
	defaultCodec = &jsonCodec{}
}

type messageProcess struct {
	messages map[string]MessageCloneable
	handlers map[string]MessageHandler
}

func (m *messageProcess) RegisterMessageAndRouter(method string, v MessageCloneable, handler MessageHandler) {
	m.messages[method] = v
	m.handlers[method] = handler
}

func (m *messageProcess) Clone(method string) interface{} {
	return m.messages[method]
}

func newMessageProcess() MessageProcess {
	mp := new(messageProcess)
	mp.messages = make(map[string]MessageCloneable)
	mp.handlers = make(map[string]MessageHandler)
	return mp
}

func (m *messageProcess) Route(method string, msg interface{}, agent *Agent) bool {
	handler, ok := m.handlers[method]
	if ok {
		handler(msg, agent)
		return true
	}
	return false
}

func RegisterHandler(h Handler) {
	defaultMessageProcess.RegisterMessageAndRouter(h.Route, h.Message, h.Handler)
}

type jsonCodec struct{}

func (j jsonCodec) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (j jsonCodec) UnMarshal(bytes []byte, i interface{}) error {
	return json.Unmarshal(bytes, i)
}

type Handler struct {
	Route   string
	Message MessageCloneable
	Handler MessageHandler
}

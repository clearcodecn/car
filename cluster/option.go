package cluster

import (
	"github.com/clearcodecn/cargo/codec"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type ServerOptions struct {
	ConnectionIdleTimeout time.Duration
	Timeout               time.Duration
	CleanDuration         time.Duration
	ReaderBufferSize      int
	ReadChannelSize       int
	WriteBufferSize       int
	WriteChannelSize      int
	EnableCompress        bool

	IsWebsocket            bool
	CheckOrigin            func(r *http.Request) bool
	BeforeWebsocketConnect func(w http.ResponseWriter, r *http.Request)
	WebsocketPath          string

	Cert string
	Key  string
	Host string
	Port string

	Codec codec.Codec

	LogLevel logrus.Level
}

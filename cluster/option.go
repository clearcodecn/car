package cluster

import (
	"net/http"
	"time"
)

type ServerOptions struct {
	ConnectionIdleTimeout time.Duration
	Timeout               time.Duration
	ReaderBufferSize      int
	WriteBufferSize       int
	EnableCompress        bool

	CheckOrigin   func(r *http.Request) bool
	WebsocketPath string

	Cert string
	Key  string
	Host string
	Port string
}

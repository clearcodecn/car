package cargo

import (
	"net/http"
	"time"
)

type ServerOption struct {
	ConnectionIdleTimeout time.Duration
	Timeout               time.Duration
	ReaderBufferSize      int
	WriteBufferSize       int
	EnableCompress        bool
	CheckOrigin           func(r *http.Request) bool
}

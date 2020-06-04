package cargo

import (
	"github.com/clearcodecn/cargo/cluster"
	"github.com/clearcodecn/cargo/codec"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"time"
)

type Option func(options *cluster.ServerOptions)

func WithTimeout(t time.Duration) Option {
	return func(options *cluster.ServerOptions) {
		options.Timeout = t
	}
}

func WithCodec(codec codec.Codec) Option {
	return func(options *cluster.ServerOptions) {
		options.Codec = codec
	}
}

func WithRWBufferSize(rsize int, wsize int) Option {
	return func(options *cluster.ServerOptions) {
		options.ReaderBufferSize = rsize
		options.WriteBufferSize = wsize
	}
}

func EnableCompress() Option {
	return func(options *cluster.ServerOptions) {
		options.EnableCompress = true
	}
}

func WithWebsocketPath(path string) Option {
	return func(options *cluster.ServerOptions) {
		options.WebsocketPath = path
	}
}

func WithTls(cert string, key string, ) Option {
	return func(options *cluster.ServerOptions) {
		options.Cert = cert
		options.Key = key
	}
}

func WithAddr(addr string) Option {
	return func(options *cluster.ServerOptions) {
		host, port, _ := net.SplitHostPort(addr)
		options.Host = host
		options.Port = port
	}
}

func WithCheckOrigin(f func(r *http.Request) bool) Option {
	return func(options *cluster.ServerOptions) {
		options.CheckOrigin = f
	}
}

func WithLogLevel(level logrus.Level) Option {
	return func(options *cluster.ServerOptions) {
		options.LogLevel = level
	}
}

func WithWebsocket() Option {
	return func(options *cluster.ServerOptions) {
		options.IsWebsocket = true
	}
}

func WithTCP() Option {
	return func(options *cluster.ServerOptions) {
		options.IsWebsocket = false
	}
}

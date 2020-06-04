package cargo

import (
	"github.com/clearcodecn/cargo/cluster"
	"github.com/clearcodecn/cargo/codec"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

var (
	node          *cluster.Node
	defaultOption cluster.ServerOptions

	defaultTimeout         = 30 * time.Second
	defaultReadBufferSize  = 2048
	defaultWriteBufferSize = 2048

	defaultReadChannelSize  = 1024
	defaultWriteChannelSize = 1024

	defaultWebsocketPath = "/ws"

	defaultHost = "0.0.0.0"
	defaultPort = "6300"

	defaultCodec = &codec.JsonCodec{}

	defaultCleanDuration = 5 * time.Minute
)

func init() {
	defaultOption = cluster.ServerOptions{
		ConnectionIdleTimeout:  0,
		Timeout:                defaultTimeout,
		ReaderBufferSize:       defaultReadBufferSize,
		ReadChannelSize:        defaultReadChannelSize,
		WriteBufferSize:        defaultWriteBufferSize,
		WriteChannelSize:       defaultWriteChannelSize,
		EnableCompress:         false,
		CheckOrigin:            defaultCheckOrigin,
		BeforeWebsocketConnect: nil,
		WebsocketPath:          defaultWebsocketPath,
		Cert:                   "",
		Key:                    "",
		Host:                   defaultHost,
		Port:                   defaultPort,
		Codec:                  defaultCodec,
		LogLevel:               logrus.DebugLevel,
		IsWebsocket:            true,
		CleanDuration:          defaultCleanDuration,
	}

	node = cluster.NewNode(defaultOption)
}

func defaultCheckOrigin(r *http.Request) bool {
	return true
}

func Default() *cluster.Node {
	return node
}

func New(opt ...Option) *cluster.Node {
	for _, o := range opt {
		o(&defaultOption)
	}
	node = cluster.NewNode(defaultOption)
	return node
}

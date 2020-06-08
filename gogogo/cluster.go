package gogogo

import (
	"github.com/clearcodecn/cargo/gogogo/pkg/registry"
	"sync"
)

type Cluster struct {
	mu    sync.Mutex
	nodes map[string]*Server

	registry registry.ServiceRegistry
}

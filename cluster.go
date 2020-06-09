package cargo

import (
	"github.com/clearcodecn/cargo/pkg/registry"
	"sync"
)

type Cluster struct {
	mu    sync.Mutex
	nodes map[string]*Server

	registry registry.ServiceRegistry
}

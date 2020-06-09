package car

import (
	"github.com/clearcodecn/car/pkg/registry"
	"sync"
)

type Cluster struct {
	mu    sync.Mutex
	nodes map[string]*Server

	registry registry.ServiceRegistry
}

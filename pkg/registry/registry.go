package registry

import "context"

type Service struct {
	Schema string
	Ip     string
	Port   string
	Label  map[string]string
}

type ServiceFilter func(s *Service) bool

type ServiceRegistry interface {
	Register(ctx context.Context, service *Service) error
	DeRegistry(ctx context.Context, service *Service) error
}

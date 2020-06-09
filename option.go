package cargo

import "net/http"

type Option struct {
	BeforeConnect func(w http.ResponseWriter, r *http.Request)

	PushServer PushServer

	CheckOrigin func(r *http.Request) bool
}

type Options func(option *Option)

func WithPushServer(ps PushServer) Options {
	return func(o *Option) {
		o.PushServer = ps
	}
}

func WithDefaultPushServer(ps PushServer) Options {
	return func(o *Option) {
		o.PushServer = defaultPushServer
	}
}

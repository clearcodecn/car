package gogogo

import "net/http"

type Option struct {
	BeforeConnect func(w http.ResponseWriter, r *http.Request)
}

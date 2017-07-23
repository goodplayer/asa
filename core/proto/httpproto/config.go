package httpproto

import "net/http"

type HttpHandlerConfig struct {
	Header map[string]func(r *http.Request) map[string]string
}

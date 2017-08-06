package constant

import "net/http"

const (
	DefaultServerName = "_default_"
	ServerHeader      = "asa"
)

type ProxyType int

const (
	HTTP_PROXY ProxyType = iota
	FCGI_PROXY
	FILE
)

type ProxyFunc func(w http.ResponseWriter, r *http.Request)

type ProxyParam struct {
	Header map[string]string
}

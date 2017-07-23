package varfunc

import (
	"net/http"
	"strings"
)

var httpVarMap map[string]func(*http.Request) string

func init() {
	httpVarMap = make(map[string]func(*http.Request) string)
	httpVarMap["remote_addr"] = remoteAddr
	httpVarMap["http_proxy_add_x_forwarded_for"] = xForwardedFor
}

func GetAllHttpInRaw() map[string]func(*http.Request) string {
	return httpVarMap
}

func remoteAddr(r *http.Request) string {
	return stripPort(r.RemoteAddr)
}

func xForwardedFor(r *http.Request) string {
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded == "" {
		return stripPort(r.RemoteAddr)
	} else {
		return forwarded + ", " + forwarded
	}
}

func stripPort(hostport string) string {
	colon := strings.IndexByte(hostport, ':')
	if colon == -1 {
		return hostport
	}
	if i := strings.IndexByte(hostport, ']'); i != -1 {
		return strings.TrimPrefix(hostport[:i], "[")
	}
	return hostport[:colon]
}

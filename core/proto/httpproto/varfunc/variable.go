package varfunc

import (
	"net/http"

	"github.com/goodplayer/asa/util"
)

var httpVarMap map[string]func(*http.Request) string

func init() {
	httpVarMap = make(map[string]func(*http.Request) string)
	httpVarMap["remote_addr"] = remoteAddr
	httpVarMap["http_proxy_add_x_forwarded_for"] = xForwardedFor
	httpVarMap["http_host"] = httpHost
}

func GetAllHttpInRaw() map[string]func(*http.Request) string {
	return httpVarMap
}

func remoteAddr(r *http.Request) string {
	return util.SplitHost(r.RemoteAddr)
}

func xForwardedFor(r *http.Request) string {
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded == "" {
		return util.SplitHost(r.RemoteAddr)
	} else {
		return forwarded + ", " + forwarded
	}
}

func httpHost(r *http.Request) string {
	return r.Host
}

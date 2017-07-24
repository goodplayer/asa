package main

import (
	"net/http"

	"github.com/goodplayer/asa/core/proto/httpproto"
	"github.com/goodplayer/asa/core/proto/httpproto/varfunc"
	"github.com/goodplayer/asa/core/trans"
	"github.com/goodplayer/asa/core/upstream"
)

func main() {
	ups := new(upstream.Upstream)
	ups.List = []string{"tcp://127.0.0.1:8888"}
	defaultUpstream := new(upstream.Upstream)
	defaultUpstream.List = []string{"tcp://127.0.0.1:80"}

	h := new(httpproto.HttpProxyHandler)
	h.UpstreamMap = map[string]*upstream.Upstream{
		"__default__": defaultUpstream,
		"localhost":   ups,
	}
	h.Config = httpproto.HttpHandlerConfig{
		Header: map[string]func(r *http.Request) map[string]string{
			"localhost": func(r *http.Request) map[string]string {
				m := varfunc.GetAllHttpInRaw()
				newMap := make(map[string]string)
				newMap["X-Real-IP"] = m["remote_addr"](r)
				newMap["X-Forwarded-For"] = m["http_proxy_add_x_forwarded_for"](r)
				newMap["Host"] = m["http_host"](r)
				return newMap
			},
			"__default__": func(r *http.Request) map[string]string {
				m := varfunc.GetAllHttpInRaw()
				newMap := make(map[string]string)
				newMap["X-Real-IP"] = m["remote_addr"](r)
				newMap["X-Forwarded-For"] = m["http_proxy_add_x_forwarded_for"](r)
				newMap["Host"] = m["http_host"](r)
				return newMap
			},
		},
	}

	tcpConfig := trans.TcpTransConfig{
		UseReuseport: true,
	}

	listener, err := trans.NewTcpTrans().NewListener("tcp", "0.0.0.0:8080", tcpConfig)
	if err != nil {
		panic(err)
	}

	err = httpproto.StartHttpServer(listener, h)
	if err != nil {
		panic(err)
	}
}

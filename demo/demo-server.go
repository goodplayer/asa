package main

import (
	"fmt"
	"net/http"

	"github.com/goodplayer/asa/constant"
	"github.com/goodplayer/asa/core/proto/httpproto"
	"github.com/goodplayer/asa/core/proto/httpproto/varfunc"
	"github.com/goodplayer/asa/core/proto/httpproto/vhost"
	"github.com/goodplayer/asa/core/trans"
	"github.com/goodplayer/asa/core/upstream"
)

func defaultRequestVariable(r *http.Request) map[string]string {
	m := varfunc.GetAllHttpInRaw()
	newMap := make(map[string]string)
	newMap["X-Real-IP"] = m["remote_addr"](r)
	newMap["X-Forwarded-For"] = m["http_proxy_add_x_forwarded_for"](r)
	newMap["Host"] = m["http_host"](r)
	return newMap
}

func defaultParamFunc(w http.ResponseWriter, r *http.Request) constant.ProxyParam {
	return constant.ProxyParam{
		Header: defaultRequestVariable(r),
	}
}

func main() {
	ups := new(upstream.Upstream)
	ups.List = []string{"tcp://127.0.0.1:8888"}
	defaultUpstream := new(upstream.Upstream)
	defaultUpstream.List = []string{"tcp://127.0.0.1:80"}

	keylessUpstream := new(upstream.Upstream)
	keylessUpstream.List = []string{"tcp://127.0.0.1:8999"}

	h := new(httpproto.HttpProxyHandler)
	h.VhostMap = map[string]*vhost.Vhost{
		"localhost": {
			VhostName:                "localhost",
			Passthrough:              true,
			PassthroughBackend:       ups,
			PassthroughType:          constant.HTTP_PROXY,
			PassthroughParamFunction: defaultParamFunc,
			SslEnable:                true,
			//ecdsa
			//SslCertPath:              "./demo/key/ecdsa_server.crt",
			//SslPrivateKeyPath:        "./demo/key/private_ecdsa_for_crt.key",
			//SslKeylessServer:         keylessUpstream,
			//rsa
			SslCertPath: "./demo/key/rsa_server.crt",
			//SslPrivateKeyPath: "./demo/key/private_rsa.key",
			SslKeylessServer: keylessUpstream,
		},
		constant.DefaultServerName: {
			VhostName:                constant.DefaultServerName,
			Passthrough:              true,
			PassthroughBackend:       defaultUpstream,
			PassthroughType:          constant.HTTP_PROXY,
			PassthroughParamFunction: defaultParamFunc,
			SslEnable:                true,
			SslCertPath:              "./demo/key/ecdsa_server.crt",
			SslPrivateKeyPath:        "./demo/key/private_ecdsa_for_crt.key",
		},
	}

	tcpConfig := trans.TcpTransConfig{
		UseReuseport: true,
	}

	listener, err := trans.NewTcpTrans().NewListener("tcp", "0.0.0.0:8080", tcpConfig)
	if err != nil {
		panic(err)
	}

	fmt.Println("starting...")

	go func() {
		tlsTcpConfig := trans.TcpTransConfig{
			UseReuseport: true,
		}

		tlsListener, err := trans.NewTcpTrans().NewListener("tcp", "0.0.0.0:8443", tlsTcpConfig)
		if err != nil {
			panic(err)
		}

		err = httpproto.StartHttpTlsServer(tlsListener, h, h.VhostMap)
		if err != nil {
			panic(err)
		}
	}()

	err = httpproto.StartHttpServer(listener, h)
	if err != nil {
		panic(err)
	}
}

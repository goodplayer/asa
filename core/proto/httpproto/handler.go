package httpproto

import (
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/goodplayer/asa/core/upstream"
)

type HttpProxyHandler struct {
	UpstreamMap map[string]*upstream.Upstream
	Config      HttpHandlerConfig
}

func (this *HttpProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqHost := r.Host
	if urla, err := url.Parse("http://" + reqHost); err == nil {
		reqHost = urla.Hostname()
	}

	ups, ok := this.UpstreamMap[reqHost]
	if !ok {
		ups, ok = this.UpstreamMap["__default__"]
		if !ok {
			panic(errors.New("no upstream mapped - httpproto"))
		}
	}

	//TODO select tcp
	conn, err := ups.SelectTcp(0, false)
	if err != nil {
		panic(err)
	}

	//TODO support http/https/etc.
	reqProxy, err := http.NewRequest(r.Method, "http://"+conn.Addr, r.Body)
	if err != nil {
		panic(err)
	}
	proxyHeader := reqProxy.Header
	// request header
	reqHeader := r.Header
	for k, v := range reqHeader {
		for idx, d := range v {
			if idx > 0 {
				proxyHeader.Add(k, d)
			} else {
				proxyHeader.Set(k, d)
			}
		}
	}
	// add new header k/v to reqHeader
	var newReqHeaderMap map[string]string
	newReqHeaderMapFunc, ok := this.Config.Header[reqHost]
	if !ok {
		newReqHeaderMapFunc, ok = this.Config.Header["__default__"]
	}
	if ok {
		newReqHeaderMap = newReqHeaderMapFunc(r)
		for k, v := range newReqHeaderMap {
			proxyHeader.Set(k, v)
		}
	}

	resp, err := http.DefaultClient.Do(reqProxy)
	if err != nil {
		panic(err)
	}
	defer func() {
		body := resp.Body
		if body != nil {
			body.Close()
		}
	}()

	respHeader := resp.Header
	header := w.Header()

	// add new header k/v to reqHeader
	header.Set("Server", "asa")

	// response header
	for k, v := range respHeader {
		for idx, d := range v {
			if idx > 0 {
				header.Add(k, d)
			} else {
				header.Set(k, d)
			}
		}
	}

	//TODO custom buffer size
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		panic(err)
	}
}
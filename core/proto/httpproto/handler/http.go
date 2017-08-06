package handler

import (
	"io"
	"net/http"
	"time"

	"github.com/goodplayer/asa/constant"
	"github.com/goodplayer/asa/core/upstream"
	"github.com/goodplayer/asa/util"
)

var httpClient *http.Client

func init() {
	httpClient = new(http.Client)
	httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
}

func HttpHandle(w http.ResponseWriter, r *http.Request, ups *upstream.Upstream, paramFunc func(w http.ResponseWriter, r *http.Request) constant.ProxyParam) {

	reqTime := time.Now()
	//TODO select stream (tcp & unix socket & etc.)
	conn, err := ups.SelectTcp(util.HashIntInt(reqTime.Nanosecond()), false)
	if err != nil {
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
	}

	//TODO support http/https/etc. & support unix socket
	reqProxy, err := http.NewRequest(r.Method, "http://"+conn.Addr+r.URL.RequestURI(), r.Body)
	if err != nil {
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
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
	if paramFunc != nil {
		newReqHeaderMap = paramFunc(w, r).Header
		for k, v := range newReqHeaderMap {
			proxyHeader.Set(k, v)
		}
		//TODO other header options
		reqProxy.Host = newReqHeaderMap["Host"]
	}

	resp, err := httpClient.Do(reqProxy)
	if err != nil {
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
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
	header.Set("Server", constant.ServerHeader)

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

	w.WriteHeader(resp.StatusCode)

	//TODO custom buffer size
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
	}
}

package httpproto

import (
	"errors"
	"net/http"

	"github.com/goodplayer/asa/constant"
	"github.com/goodplayer/asa/core/proto/httpproto/vhost"
	"github.com/goodplayer/asa/util"
)

type HttpProxyHandler struct {
	HttpProxyClient *http.Client
	VhostMap        map[string]*vhost.Vhost
}

func (this *HttpProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqHost := util.SplitHost(r.Host)

	vh, ok := this.VhostMap[reqHost]
	if !ok {
		vh, ok = this.VhostMap[constant.DefaultServerName]
		if !ok {
			panic(errors.New("no upstream mapped - httpproto"))
		}
	}

	url := r.URL
	path := url.Path // decode & dereference, not combine slash
	// URI processing
	internalRewriteUri := vh.InternalRewriteUri(path)
	redirectUrl, ok, statusCode := vh.Rewrite(internalRewriteUri)
	if ok {
		http.Redirect(w, r, redirectUrl, statusCode)
		return
	}
	finalUri := normalizeUri(internalRewriteUri) // %xx / reference / slash

	proxyFunc := vh.Location(finalUri)

	if proxyFunc != nil {
		proxyFunc(w, r)
	} else {
		// 404
		w.WriteHeader(http.StatusNotFound)
	}
}

func normalizeUri(uri string) string {
	//TODO normalize uri
	return uri
}

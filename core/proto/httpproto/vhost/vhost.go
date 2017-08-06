package vhost

import (
	"net/http"

	"github.com/goodplayer/asa/constant"
	"github.com/goodplayer/asa/core/proto/httpproto/handler"
	"github.com/goodplayer/asa/core/upstream"
)

type Vhost struct {
	VhostName string

	SslEnable   bool
	SslCertPath string // only support RSA/ECDSA cert
	// set only one configuration below
	SslPrivateKeyPath string
	SslKeylessServer  *upstream.Upstream

	Passthrough              bool
	PassthroughBackend       *upstream.Upstream
	PassthroughType          constant.ProxyType
	PassthroughParamFunction func(w http.ResponseWriter, r *http.Request) constant.ProxyParam

	//TODO location
}

func (this *Vhost) InternalRewriteUri(uri string) string {
	//TODO
	return uri
}

func (this *Vhost) Rewrite(uri string) (string, bool, int) {
	//TODO client rewrite
	return uri, false, 302
}

func (this *Vhost) Location(normalizedUri string) constant.ProxyFunc {
	if this.Passthrough {
		switch this.PassthroughType {
		case constant.HTTP_PROXY:
			return func(w http.ResponseWriter, r *http.Request) {
				handler.HttpHandle(w, r, this.PassthroughBackend, this.PassthroughParamFunction)
			}
		default:
			return func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Bad Gateway", http.StatusBadGateway)
			}
		}
	}

	//TODO matching
	//TODO =  exact matching  (terminate maching)
	//TODO    prefix matching  (longest first) (^~ do not check regex)
	//TODO    suffix matching
	//TODO ~* regex case-insensitive
	//TODO ~  regex case-sensitive
	//TODO    determine regex or else not found then using prefix matching

	return nil
}

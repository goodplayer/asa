package httpproto_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/goodplayer/asa/core/proto/httpproto"
	"github.com/goodplayer/asa/core/trans"
)

func TestStartHttpServer(t *testing.T) {
	tcpConfig := trans.TcpTransConfig{}
	tcpConfig.UseReuseport = true

	l, err := trans.NewTcpTrans().NewListener("tcp", ":8012", tcpConfig)
	if err != nil {
		panic(err)
	}

	err = httpproto.StartHttpServer(l, handler{})
	if err != nil {
		panic(err)
	}
}

type handler struct {
}

func (this handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Header)
}

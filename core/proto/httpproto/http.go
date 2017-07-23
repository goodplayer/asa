package httpproto

import (
	"net"
	"net/http"
	"time"
)

func StartHttpServer(l net.Listener, handler http.Handler) error {
	srv := &http.Server{Handler: handler}
	return srv.Serve(keepAliveListener{l.(*net.TCPListener)})
}

func StartHttpsServer() {
	//TODO
	//http.ListenAndServeTLS()
}

type keepAliveListener struct {
	*net.TCPListener
}

func (ln keepAliveListener) Accept() (net.Conn, error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return nil, err
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}
